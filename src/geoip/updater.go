package geoip

import (
	"fmt"
	"log"
	"time"
)

// UpdaterConfig holds configuration for the database updater
type UpdaterConfig struct {
	DataDir        string
	CheckInterval  time.Duration // How often to check for updates
	AutoUpdate     bool          // Whether to automatically update
	OnUpdateFunc   func()        // Callback after successful update
	OnErrorFunc    func(error)   // Callback on error
}

// Updater manages automatic GeoIP database updates
type Updater struct {
	config  *UpdaterConfig
	stopCh  chan struct{}
	running bool
}

// NewUpdater creates a new database updater
func NewUpdater(config *UpdaterConfig) *Updater {
	if config.CheckInterval == 0 {
		config.CheckInterval = 24 * time.Hour // Default: check daily
	}

	return &Updater{
		config: config,
		stopCh: make(chan struct{}),
	}
}

// Start begins the automatic update checker
func (u *Updater) Start() {
	if u.running {
		return
	}

	u.running = true
	go u.run()
}

// Stop stops the automatic update checker
func (u *Updater) Stop() {
	if !u.running {
		return
	}

	u.running = false
	close(u.stopCh)
}

// run is the main update loop
func (u *Updater) run() {
	ticker := time.NewTicker(u.config.CheckInterval)
	defer ticker.Stop()

	// Check immediately on start
	u.checkAndUpdate()

	for {
		select {
		case <-ticker.C:
			u.checkAndUpdate()
		case <-u.stopCh:
			return
		}
	}
}

// checkAndUpdate checks for updates and downloads if available
func (u *Updater) checkAndUpdate() {
	log.Println("Checking for GeoIP database updates...")

	// Get current version (from file metadata or release tag)
	currentVersion := u.getCurrentVersion()

	// Check for updates
	hasUpdate, newVersion, err := CheckForUpdates(currentVersion)
	if err != nil {
		log.Printf("Error checking for updates: %v", err)
		if u.config.OnErrorFunc != nil {
			u.config.OnErrorFunc(err)
		}
		return
	}

	if !hasUpdate {
		log.Printf("GeoIP databases are up to date (version: %s)", currentVersion)
		return
	}

	log.Printf("New GeoIP database version available: %s (current: %s)", newVersion, currentVersion)

	// Only auto-update if configured
	if !u.config.AutoUpdate {
		log.Println("Auto-update disabled. Skipping download.")
		return
	}

	// Download new databases
	log.Println("Downloading updated databases...")
	dbFiles, err := DownloadDatabases(u.config.DataDir)
	if err != nil {
		log.Printf("Error downloading databases: %v", err)
		if u.config.OnErrorFunc != nil {
			u.config.OnErrorFunc(err)
		}
		return
	}

	// Reload the GeoIP instance with new databases
	if instance := GetInstance(); instance != nil {
		if err := instance.Reload(dbFiles.CityIPv4DB, dbFiles.CityIPv6DB, dbFiles.CountryDB, dbFiles.ASNDB); err != nil {
			log.Printf("Error reloading databases: %v", err)
			if u.config.OnErrorFunc != nil {
				u.config.OnErrorFunc(err)
			}
			return
		}
	}

	// Save new version
	u.saveCurrentVersion(newVersion)

	log.Printf("Successfully updated GeoIP databases to version %s", newVersion)

	// Call update callback
	if u.config.OnUpdateFunc != nil {
		u.config.OnUpdateFunc()
	}
}

// getCurrentVersion reads the current database version
func (u *Updater) getCurrentVersion() string {
	// TODO: Store version in a file or database
	// For now, return empty to always check
	return ""
}

// saveCurrentVersion saves the current database version
func (u *Updater) saveCurrentVersion(version string) {
	// TODO: Store version in a file or database
	// For now, do nothing
}

// ManualUpdate triggers a manual database update
func (u *Updater) ManualUpdate() error {
	log.Println("Manual GeoIP database update triggered...")

	// Download new databases
	dbFiles, err := DownloadDatabases(u.config.DataDir)
	if err != nil {
		return fmt.Errorf("failed to download databases: %w", err)
	}

	// Reload the GeoIP instance
	if instance := GetInstance(); instance != nil {
		if err := instance.Reload(dbFiles.CityIPv4DB, dbFiles.CityIPv6DB, dbFiles.CountryDB, dbFiles.ASNDB); err != nil {
			return fmt.Errorf("failed to reload databases: %w", err)
		}
	}

	log.Println("Manual update completed successfully")
	return nil
}

// GetScheduledTask returns a function suitable for use with a cron scheduler
func GetScheduledTask(dataDir string) func() {
	return func() {
		log.Println("Scheduled GeoIP database update starting...")

		// Download databases
		dbFiles, err := DownloadDatabases(dataDir)
		if err != nil {
			log.Printf("Scheduled update failed: %v", err)
			return
		}

		// Reload databases
		if instance := GetInstance(); instance != nil {
			if err := instance.Reload(dbFiles.CityIPv4DB, dbFiles.CityIPv6DB, dbFiles.CountryDB, dbFiles.ASNDB); err != nil {
				log.Printf("Failed to reload databases: %v", err)
				return
			}
		}

		log.Println("Scheduled GeoIP database update completed successfully")
	}
}
