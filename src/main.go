package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/apimgr/zipcodes/src/config"
	"github.com/apimgr/zipcodes/src/database"
	"github.com/apimgr/zipcodes/src/geoip"
	"github.com/apimgr/zipcodes/src/mode"
	"github.com/apimgr/zipcodes/src/paths"
	"github.com/apimgr/zipcodes/src/server"
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
	var (
		port            = flag.String("port", "", "Port to listen on")
		address         = flag.String("address", "", "Address to bind to")
		configDir       = flag.String("config", "", "Configuration directory")
		dataDir         = flag.String("data", "", "Data directory")
		logsDir         = flag.String("logs", "", "Logs directory")
		showVersion     = flag.Bool("version", false, "Show version information")
		showStatus      = flag.Bool("status", false, "Show server status (for health checks)")
		showHelp        = flag.Bool("help", false, "Show help message")
		serviceCmd      = flag.String("service", "", "Service command (install, uninstall, start, stop, restart, status)")
		maintenanceMode = flag.String("maintenance", "", "Maintenance mode (on/off)")
		modeFlag        = flag.String("mode", "", "Application mode (dev/development, prod/production)")
		updateCmd       = flag.String("update", "", "Update command (stable, beta, nightly)")
	)
	flag.Parse()

	// Show help
	if *showHelp {
		printHelp()
		return
	}

	// Show version
	if *showVersion {
		fmt.Println(Version)
		return
	}

	// Handle update command
	if *updateCmd != "" {
		handleUpdateCommand(*updateCmd)
		return
	}

	// Initialize mode
	if err := mode.Initialize(*modeFlag); err != nil {
		log.Printf("Warning: invalid mode: %v", err)
	}

	// Status check (for health checks)
	if *showStatus {
		os.Exit(checkServerStatus())
	}

	// Handle service commands
	if *serviceCmd != "" {
		handleServiceCommand(*serviceCmd)
		return
	}

	// Handle maintenance mode
	if *maintenanceMode != "" {
		handleMaintenanceMode(*maintenanceMode)
		return
	}

	log.Printf("Starting Zipcodes API v%s", Version)

	// Determine directories
	appName := "zipcodes"
	if *configDir == "" {
		*configDir = os.Getenv("CONFIG_DIR")
		if *configDir == "" {
			*configDir = paths.GetConfigDir(appName)
		}
	}
	if *dataDir == "" {
		*dataDir = os.Getenv("DATA_DIR")
		if *dataDir == "" {
			*dataDir = paths.GetDataDir(appName)
		}
	}
	if *logsDir == "" {
		*logsDir = os.Getenv("LOGS_DIR")
		if *logsDir == "" {
			*logsDir = paths.GetLogsDir(appName)
		}
	}

	// Ensure directories exist
	for _, dir := range []string{*configDir, *dataDir, *logsDir} {
		if err := paths.EnsureDir(dir); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	log.Printf("Config directory: %s", *configDir)
	log.Printf("Data directory: %s", *dataDir)
	log.Printf("Logs directory: %s", *logsDir)

	// Load configuration
	configPath := filepath.Join(*configDir, "server.yml")
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Printf("Warning: Failed to load config: %v (using defaults)", err)
		cfg = config.DefaultConfig()
	}

	// Initialize database for zipcode data
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = filepath.Join(*dataDir, "zipcodes.db")
	}

	log.Printf("Database path: %s", dbPath)
	db, err := database.Initialize(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("Database initialized successfully")

	// Load zipcode data from embedded JSON
	log.Println("Loading zipcode data from embedded JSON...")
	if err := db.LoadFromJSON(zipcodesData); err != nil {
		log.Fatalf("Failed to load zipcode data: %v", err)
	}

	// Get stats
	stats, err := db.GetStats()
	if err != nil {
		log.Printf("Warning: Failed to get stats: %v", err)
	} else {
		log.Printf("Loaded %d zipcodes from %d states", stats["total_zipcodes"], stats["total_states"])
	}

	// Initialize GeoIP databases (zipcodes needs location features)
	if err := initializeGeoIP(*dataDir); err != nil {
		log.Printf("Warning: GeoIP initialization failed: %v", err)
		log.Println("GeoIP features will be unavailable")
	} else {
		log.Println("GeoIP databases initialized successfully")
	}

	// Determine server configuration
	if *port == "" {
		*port = os.Getenv("PORT")
		if *port == "" {
			*port = cfg.Server.Port
		}
		if *port == "" {
			// Generate random port in range 64000-64999
			*port = strconv.Itoa(64000 + rand.Intn(1000))
		}
	}
	if *address == "" {
		*address = os.Getenv("ADDRESS")
		if *address == "" {
			*address = cfg.Server.Address
		}
		if *address == "" {
			*address = "0.0.0.0"
		}
	}

	// Create HTTP server
	srv := server.New(db, cfg, *address, *port, Version, BuildDate, Commit, zipcodesData)

	// Setup graceful shutdown
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", *address, *port),
		Handler:      srv.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on %s:%s", *address, *port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

func initializeGeoIP(dataDir string) error {
	// Check if databases already exist
	if !geoip.DatabasesExist(dataDir) {
		log.Println("GeoIP databases not found. Downloading from GitHub...")

		// Download databases
		dbFiles, err := geoip.DownloadDatabases(dataDir)
		if err != nil {
			return fmt.Errorf("failed to download databases: %w", err)
		}

		log.Printf("Downloaded databases:")
		if dbFiles.CityIPv4DB != "" {
			log.Printf("  - City IPv4: %s", dbFiles.CityIPv4DB)
		}
		if dbFiles.CityIPv6DB != "" {
			log.Printf("  - City IPv6: %s", dbFiles.CityIPv6DB)
		}
		if dbFiles.CountryDB != "" {
			log.Printf("  - Country: %s", dbFiles.CountryDB)
		}
		if dbFiles.ASNDB != "" {
			log.Printf("  - ASN: %s", dbFiles.ASNDB)
		}
	} else {
		log.Println("Found existing GeoIP databases")
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

// handleServiceCommand handles service management commands
func handleServiceCommand(cmd string) {
	switch cmd {
	case "install":
		fmt.Println("Service installation not yet implemented")
		fmt.Println("Use systemd/launchd/rc.d to manage the service")
	case "uninstall":
		fmt.Println("Service uninstallation not yet implemented")
	case "start":
		fmt.Println("Use 'systemctl start zipcodes' or run the binary directly")
	case "stop":
		fmt.Println("Use 'systemctl stop zipcodes' or send SIGTERM to the process")
	case "restart":
		fmt.Println("Use 'systemctl restart zipcodes'")
	case "status":
		fmt.Println("Use 'systemctl status zipcodes' or --status flag")
	default:
		fmt.Printf("Unknown service command: %s\n", cmd)
		fmt.Println("Available commands: install, uninstall, start, stop, restart, status")
	}
}

// handleMaintenanceMode handles maintenance mode toggle
func handleMaintenanceMode(m string) {
	switch m {
	case "on":
		fmt.Println("Maintenance mode: ON")
		fmt.Println("Note: Maintenance mode is handled at runtime, not persisted")
	case "off":
		fmt.Println("Maintenance mode: OFF")
	default:
		fmt.Printf("Invalid maintenance mode: %s (use 'on' or 'off')\n", m)
	}
}

func printHelp() {
	fmt.Printf(`Zipcodes API v%s

Usage: zipcodes [options]

Options:
  --help             Show this help message
  --version          Show version information
  --status           Check service status
  --config DIR       Set configuration directory
  --data DIR         Set data directory
  --logs DIR         Set logs directory
  --address ADDR     Set listen address (default: 0.0.0.0)
  --port PORT        Set listen port (default: random)
  --mode MODE        Application mode (dev, prod)
  --update BRANCH    Update from branch (stable, beta, nightly)

Service Management:
  --service install    Install as system service
  --service uninstall  Uninstall system service
  --service start      Start the service
  --service stop       Stop the service
  --service restart    Restart the service
  --service status     Show service status

Maintenance:
  --maintenance on     Enable maintenance mode
  --maintenance off    Disable maintenance mode

Examples:
  zipcodes                          Start with default settings
  zipcodes --port 8080              Start on port 8080
  zipcodes --mode dev --port 8080   Start in development mode
  zipcodes --update stable          Update to stable branch

Documentation: https://zipcodes.apimgr.us
`, Version)
}

func handleUpdateCommand(branch string) {
	validBranches := map[string]bool{
		"stable":  true,
		"beta":    true,
		"nightly": true,
	}

	if !validBranches[branch] {
		fmt.Printf("Error: invalid update branch %q (valid: stable, beta, nightly)\n", branch)
		os.Exit(1)
	}

	fmt.Printf("Updating Zipcodes API from %s branch...\n", branch)

	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("Error: git is not installed")
		os.Exit(1)
	}

	cmd := exec.Command("git", "pull", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Update failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Update complete. Please rebuild the application.")
}
