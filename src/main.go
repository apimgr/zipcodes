package main

import (
	_ "embed"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/apimgr/zipcodes/src/database"
	"github.com/apimgr/zipcodes/src/geoip"
	"github.com/apimgr/zipcodes/src/paths"
	"github.com/apimgr/zipcodes/src/server"
	"github.com/apimgr/zipcodes/src/utils"
)

//go:embed data/zipcodes.json
var zipcodesData []byte

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
	dbPath := flag.String("db-path", "", "Set SQLite database path")
	devMode := flag.Bool("dev", false, "Run in development mode")

	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Println(Version)
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
		fmt.Println("  --db-path PATH    Set SQLite database path")
		fmt.Println("  --dev             Run in development mode")
		fmt.Println("\nEnvironment Variables:")
		fmt.Println("  CONFIG_DIR        Configuration directory")
		fmt.Println("  DATA_DIR          Data directory")
		fmt.Println("  LOGS_DIR          Logs directory")
		fmt.Println("  DB_PATH           SQLite database path")
		fmt.Println("  PORT              Server port")
		fmt.Println("  ADDRESS           Listen address")
		fmt.Println("  ADMIN_USER        Admin username (first run only)")
		fmt.Println("  ADMIN_PASSWORD    Admin password (first run only)")
		fmt.Println("  ADMIN_TOKEN       Admin API token (first run only)")
		os.Exit(0)
	}

	// Handle status flag
	if *showStatus {
		os.Exit(checkServerStatus())
	}

	// Store configuration
	config := &Config{
		Port:      *port,
		Address:   *address,
		DataDir:   *dataDir,
		ConfigDir: *configDir,
		LogsDir:   *logsDir,
		DBPath:    *dbPath,
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
	DBPath    string
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

	// Determine database path with priority order:
	// 1. Command-line flag
	// 2. Environment variable DB_PATH
	// 3. Default: {DATA_DIR}/zipcodes.db
	dbPath := config.DBPath
	if dbPath == "" {
		dbPath = os.Getenv("DB_PATH")
	}
	if dbPath == "" {
		dbPath = filepath.Join(dataDir, "zipcodes.db")
	}

	fmt.Printf("📂 Database path: %s\n", dbPath)
	db, err := database.NewAppDB(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	fmt.Println("✅ Database initialized successfully")

	// Load zipcode data from embedded JSON
	fmt.Println("📥 Loading zipcode data from embedded JSON...")

	if err := db.LoadFromJSON(zipcodesData); err != nil {
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
		// Note: rand is auto-seeded in Go 1.20+, no need for rand.Seed()
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

	// Display admin credentials if they were just created (with port)
	if err := database.DisplayAdminCredentials(db.GetConn(), port, address); err != nil {
		fmt.Printf("Warning: Failed to display credentials: %v\n", err)
	}

	// Create and start server
	srv := server.New(db, port, zipcodesData)

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
		if dbFiles.CityIPv4DB != "" {
			fmt.Printf("  - City IPv4: %s\n", dbFiles.CityIPv4DB)
		}
		if dbFiles.CityIPv6DB != "" {
			fmt.Printf("  - City IPv6: %s\n", dbFiles.CityIPv6DB)
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
	if err := geoip.Initialize(dbPaths.CityIPv4DB, dbPaths.CityIPv6DB, dbPaths.CountryDB, dbPaths.ASNDB); err != nil {
		return fmt.Errorf("failed to initialize GeoIP: %w", err)
	}

	return nil
}

// checkServerStatus checks if the server is running and healthy
// Returns exit code: 0 = healthy, 1 = unhealthy
func checkServerStatus() int {
	// Get port from environment or default
	port := os.Getenv("PORT")
	if port == "" {
		// Try to find a running instance by checking common ports
		// or reading from PID file
		fmt.Println("Status: Unknown (no PORT specified)")
		fmt.Println("Hint: Set PORT environment variable or use --port flag")
		return 1
	}

	// Try to connect to health endpoint
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	healthURL := fmt.Sprintf("http://127.0.0.1:%s/healthz", port)
	resp, err := client.Get(healthURL)
	if err != nil {
		fmt.Printf("Status: Unhealthy (cannot connect to %s)\n", healthURL)
		fmt.Printf("Error: %v\n", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Status: Healthy")
		fmt.Printf("Server: Running on port %s\n", port)
		return 0
	}

	fmt.Printf("Status: Unhealthy (HTTP %d)\n", resp.StatusCode)
	return 1
}
