package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/apimgr/zipcodes/src/database"
	"github.com/apimgr/zipcodes/src/geoip"
	"github.com/apimgr/zipcodes/src/paths"
	"github.com/apimgr/zipcodes/src/server"
	"github.com/apimgr/zipcodes/src/utils"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Command-line flags
	showVersion := flag.Bool("version", false, "Show version information")
	showStatus := flag.Bool("status", false, "Show server status and exit")
	showHelp := flag.Bool("help", false, "Show help message")
	port := flag.String("port", "", "Set port (default: random 64000-64999)")
	address := flag.String("address", "0.0.0.0", "Set listen address")
	dataDir := flag.String("data", "", "Set data directory")
	configDir := flag.String("config", "", "Set config directory")
	logsDir := flag.String("logs", "", "Set logs directory")
	devMode := flag.Bool("dev", false, "Run in development mode")

	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("zipcodes version %s\n", Version)
		fmt.Printf("Built: %s\n", BuildDate)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Go: %s\n", runtime.Version())
		fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// Handle help flag
	if *showHelp {
		fmt.Println("Usage: zipcodes [OPTIONS]")
		fmt.Println("\nOptions:")
		fmt.Println("  --help            Show this help message")
		fmt.Println("  --version         Show version information")
		fmt.Println("  --status          Show server status and exit with code")
		fmt.Println("  --port PORT       Set port (default: random 64000-64999)")
		fmt.Println("  --address ADDR    Set listen address (default: 0.0.0.0)")
		fmt.Println("  --config DIR      Set config directory")
		fmt.Println("  --data DIR        Set data directory")
		fmt.Println("  --logs DIR        Set logs directory")
		fmt.Println("  --dev             Run in development mode")
		fmt.Println("\nEnvironment Variables:")
		fmt.Println("  CONFIG_DIR        Configuration directory")
		fmt.Println("  DATA_DIR          Data directory")
		fmt.Println("  LOGS_DIR          Logs directory")
		fmt.Println("  PORT              Server port")
		fmt.Println("  ADDRESS           Listen address")
		fmt.Println("  ADMIN_USER        Admin username (first run only)")
		fmt.Println("  ADMIN_PASSWORD    Admin password (first run only)")
		fmt.Println("  ADMIN_TOKEN       Admin API token (first run only)")
		os.Exit(0)
	}

	// Handle status flag
	if *showStatus {
		// TODO: Implement status check
		fmt.Println("❌ Server: Not running")
		os.Exit(1)
	}

	// Store configuration
	config := &Config{
		Port:      *port,
		Address:   *address,
		DataDir:   *dataDir,
		ConfigDir: *configDir,
		LogsDir:   *logsDir,
		DevMode:   *devMode,
	}

	// Start server
	fmt.Printf("Starting zipcodes v%s...\n", Version)
	if err := StartServer(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type Config struct {
	Port      string
	Address   string
	DataDir   string
	ConfigDir string
	LogsDir   string
	DevMode   bool
}

func StartServer(config *Config) error {
	// Get OS-specific directories with priority order:
	// 1. Command-line flags (highest)
	// 2. Environment variables
	// 3. OS-specific defaults (lowest)
	configDir, dataDir, logsDir := paths.GetDirs("zipcodes", config.ConfigDir, config.DataDir, config.LogsDir)

	// Set CONFIG_DIR environment variable for admin credentials
	os.Setenv("CONFIG_DIR", configDir)

	// Create directories
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	fmt.Printf("📂 Config directory: %s\n", configDir)
	fmt.Printf("📂 Data directory: %s\n", dataDir)
	fmt.Printf("📂 Logs directory: %s\n", logsDir)

	// Initialize database with auth tables
	dbPath := filepath.Join(dataDir, "zipcodes.db")
	db, err := database.NewAppDB(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	fmt.Println("✅ Database initialized successfully")

	// Load zipcode data from JSON
	jsonPath := filepath.Join("src", "data", "zipcodes.json")
	fmt.Printf("📥 Loading zipcode data from %s...\n", jsonPath)

	if err := db.LoadFromJSON(jsonPath); err != nil {
		return fmt.Errorf("failed to load zipcode data: %w", err)
	}

	// Initialize GeoIP databases
	if err := initializeGeoIP(dataDir); err != nil {
		fmt.Printf("⚠️  Warning: GeoIP initialization failed: %v\n", err)
		fmt.Println("   GeoIP features will be unavailable")
	} else {
		fmt.Println("✅ GeoIP databases initialized successfully")
	}

	// Determine port with priority order:
	// 1. Command-line flag
	// 2. Environment variable PORT
	// 3. Random port 64000-64999 (spec default)
	port := config.Port
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		// Generate random port in range 64000-64999
		rand.Seed(time.Now().UnixNano())
		port = strconv.Itoa(64000 + rand.Intn(1000))
	}

	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("invalid port: %s", port)
	}

	// Get listen address (flag or env or default)
	address := config.Address
	if address == "" {
		address = os.Getenv("ADDRESS")
	}
	if address == "" {
		address = "0.0.0.0"
	}

	// Create and start server
	srv := server.New(db, port)

	// Get display address (external IP, hostname, or fallback)
	displayAddr := utils.GetDisplayAddress(address)

	fmt.Println("\n🚀 Server starting...")
	fmt.Printf("   URL: http://%s:%s\n\n", displayAddr, port)

	return srv.Start(displayAddr, address)
}

func initializeGeoIP(dataDir string) error {
	// Check if databases already exist
	if !geoip.DatabasesExist(dataDir) {
		fmt.Println("GeoIP databases not found. Downloading from GitHub...")

		// Download databases
		dbFiles, err := geoip.DownloadDatabases(dataDir)
		if err != nil {
			return fmt.Errorf("failed to download databases: %w", err)
		}

		fmt.Printf("Downloaded databases:\n")
		if dbFiles.CityDB != "" {
			fmt.Printf("  - City: %s\n", dbFiles.CityDB)
		}
		if dbFiles.CountryDB != "" {
			fmt.Printf("  - Country: %s\n", dbFiles.CountryDB)
		}
		if dbFiles.ASNDB != "" {
			fmt.Printf("  - ASN: %s\n", dbFiles.ASNDB)
		}
	} else {
		fmt.Println("Found existing GeoIP databases")
	}

	// Get database paths
	dbPaths := geoip.GetDatabasePaths(dataDir)

	// Initialize GeoIP with the databases
	if err := geoip.Initialize(dbPaths.CityDB, dbPaths.CountryDB, dbPaths.ASNDB); err != nil {
		return fmt.Errorf("failed to initialize GeoIP: %w", err)
	}

	return nil
}
