package database

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/apimgr/zipcodes/src/utils"
)

// InitializeAdminSchema creates admin-only authentication tables
func InitializeAdminSchema(db *sql.DB) error {
	schema := `
	-- Admin credentials table (admin-only auth)
	CREATE TABLE IF NOT EXISTS admin_credentials (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		token_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Settings table
	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('string', 'number', 'boolean', 'json')),
		category TEXT NOT NULL,
		description TEXT,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	-- Audit log table
	CREATE TABLE IF NOT EXISTS audit_log (
		id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
		username TEXT,
		action TEXT NOT NULL,
		resource TEXT NOT NULL,
		old_value TEXT,
		new_value TEXT,
		ip_address TEXT NOT NULL,
		user_agent TEXT,
		success INTEGER NOT NULL,
		error_message TEXT,
		timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	-- Scheduled tasks table
	CREATE TABLE IF NOT EXISTS scheduled_tasks (
		id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
		name TEXT UNIQUE NOT NULL,
		cron_expression TEXT NOT NULL,
		command TEXT NOT NULL,
		enabled INTEGER DEFAULT 1,
		last_run DATETIME,
		next_run DATETIME NOT NULL,
		last_status TEXT,
		last_error TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_audit_log_timestamp ON audit_log(timestamp);
	CREATE INDEX IF NOT EXISTS idx_settings_category ON settings(category);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Insert default settings
	if err := insertAdminDefaultSettings(db); err != nil {
		return fmt.Errorf("failed to insert default settings: %w", err)
	}

	// Initialize admin credentials silently (don't display yet)
	if err := initializeAdminCredentials(db); err != nil {
		return fmt.Errorf("failed to initialize admin credentials: %w", err)
	}

	return nil
}

// DisplayAdminCredentials displays admin credentials with server URL
// Should be called AFTER port is determined
func DisplayAdminCredentials(db *sql.DB, port, address string) error {
	// Check if credentials were just created
	var username, password, token string
	var createdAt time.Time

	err := db.QueryRow(`
		SELECT username, created_at FROM admin_credentials WHERE id = 1
	`).Scan(&username, &createdAt)
	if err != nil {
		return err
	}

	// Only display if recently created (within last 10 seconds)
	if time.Since(createdAt) > 10*time.Second {
		return nil
	}

	// Get plaintext password and token from environment
	password = os.Getenv("ADMIN_PASSWORD")
	token = os.Getenv("ADMIN_TOKEN")

	if password == "" || token == "" {
		// Can't display if not in environment
		return nil
	}

	// Get config directory for file path
	configDir := os.Getenv("CONFIG_DIR")

	// Write credentials file with port
	if configDir != "" {
		writeCredentialsFileWithPort(configDir, username, password, token, port, address)
	}

	// Get display address
	displayAddr := utils.GetDisplayAddress(address)

	// Display credentials (shown once)
	fmt.Println("\n========================================")
	fmt.Println("ZIPCODES API - ADMIN CREDENTIALS")
	fmt.Println("========================================")
	fmt.Println("WEB UI LOGIN:")
	fmt.Printf("  URL:      http://%s:%s/admin\n", displayAddr, port)
	fmt.Printf("  Username: %s\n", username)
	fmt.Printf("  Password: %s\n", password)
	fmt.Println("\nAPI TOKEN:")
	fmt.Printf("  Header:   Authorization: Bearer %s\n", token)
	fmt.Printf("  Token:    %s\n", token)
	if configDir != "" {
		fmt.Printf("\nCredentials saved to: %s/admin_credentials\n", configDir)
	}
	fmt.Println("\n⚠️  Save these credentials securely!")
	fmt.Println("They will not be shown again.")
	fmt.Println("========================================\n")

	return nil
}

// insertAdminDefaultSettings adds default server settings
func insertAdminDefaultSettings(db *sql.DB) error {
	defaults := []struct {
		key         string
		value       string
		typ         string
		category    string
		description string
	}{
		{"server.title", "Zipcodes", "string", "server", "Application display name"},
		{"server.tagline", "US Postal Code Lookup API", "string", "server", "Short subtitle/slogan"},
		{"server.description", "Fast and accurate US zipcode lookup API with 340,000+ zipcodes, GeoIP integration, and modern web interface.", "string", "server", "Full description"},
		{"server.address", "0.0.0.0", "string", "server", "Listen address"},
		{"server.http_port", "64080", "number", "server", "HTTP port"},
		{"server.https_enabled", "false", "boolean", "server", "Enable HTTPS"},
		{"server.timezone", "UTC", "string", "server", "Server timezone"},
		{"server.date_format", "US", "string", "server", "Date format (US, EU, ISO)"},
		{"server.time_format", "12-hour", "string", "server", "Time format (12-hour, 24-hour)"},
		{"proxy.enabled", "true", "boolean", "proxy", "Enable reverse proxy support"},
		{"proxy.trust_headers", "true", "boolean", "proxy", "Trust proxy headers"},
		{"features.api_enabled", "true", "boolean", "features", "Enable API endpoints"},
	}

	for _, setting := range defaults {
		_, err := db.Exec(`
			INSERT OR IGNORE INTO settings (key, value, type, category, description)
			VALUES (?, ?, ?, ?, ?)
		`, setting.key, setting.value, setting.typ, setting.category, setting.description)
		if err != nil {
			return err
		}
	}

	return nil
}

// initializeAdminCredentials creates admin credentials on first run
func initializeAdminCredentials(db *sql.DB) error {
	// Check if admin already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM admin_credentials").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// Admin already exists
		return nil
	}

	// Get from environment or use defaults
	username := os.Getenv("ADMIN_USER")
	if username == "" {
		username = "administrator"
	}

	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		// Generate random password
		password = generateRandomString(16)
	}

	token := os.Getenv("ADMIN_TOKEN")
	if token == "" {
		// Generate random token
		token = generateRandomString(64)
	}

	// Hash password and token
	passwordHash := hashString(password)
	tokenHash := hashString(token)

	// Insert admin credentials
	_, err = db.Exec(`
		INSERT INTO admin_credentials (id, username, password_hash, token_hash)
		VALUES (1, ?, ?, ?)
	`, username, passwordHash, tokenHash)
	if err != nil {
		return err
	}

	// Credentials will be displayed later after port is determined
	// Save password and token to environment for later display
	os.Setenv("ADMIN_PASSWORD", password)
	os.Setenv("ADMIN_TOKEN", token)

	return nil
}

// writeCredentialsFileWithPort writes credentials to a file with proper URL including port
func writeCredentialsFileWithPort(configDir, username, password, token, port, address string) error {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	credFile := filepath.Join(configDir, "admin_credentials")

	// Get display address
	displayAddr := utils.GetDisplayAddress(address)

	content := fmt.Sprintf(`ZIPCODES API - ADMIN CREDENTIALS
========================================
WEB UI LOGIN:
  URL:      http://%s:%s/admin
  Username: %s
  Password: %s

API TOKEN:
  URL:      http://%s:%s/api/v1/admin
  Header:   Authorization: Bearer %s
  Token:    %s

Created: %s
========================================

⚠️  Keep these credentials secure!
They will not be shown again.
`, displayAddr, port, username, password, displayAddr, port, token, token, time.Now().Format("2006-01-02 15:04:05"))

	// Write file with 0600 permissions (owner read/write only)
	if err := os.WriteFile(credFile, []byte(content), 0600); err != nil {
		return err
	}

	return nil
}

// generateRandomString generates a random hex string
func generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// hashString creates a SHA-256 hash
func hashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// VerifyAdminPassword verifies admin password
func VerifyAdminPassword(db *sql.DB, username, password string) bool {
	var storedHash string
	err := db.QueryRow(`
		SELECT password_hash FROM admin_credentials
		WHERE username = ?
	`, username).Scan(&storedHash)
	if err != nil {
		return false
	}

	passwordHash := hashString(password)
	return passwordHash == storedHash
}

// VerifyAdminToken verifies admin API token
func VerifyAdminToken(db *sql.DB, token string) bool {
	var storedHash string
	err := db.QueryRow(`
		SELECT token_hash FROM admin_credentials LIMIT 1
	`).Scan(&storedHash)
	if err != nil {
		return false
	}

	tokenHash := hashString(token)
	return tokenHash == storedHash
}
