package database

import (
	"database/sql"
	"fmt"
)

// InitializeSchema creates all required tables according to Universal Server Template spec
func InitializeSchema(db *sql.DB) error {
	schema := `
	-- Users table (exact structure required by spec)
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		display_name TEXT,
		avatar_url TEXT,
		bio TEXT,
		role TEXT NOT NULL CHECK (role IN ('administrator', 'user', 'guest')) DEFAULT 'user',
		status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'pending')),
		timezone TEXT DEFAULT 'UTC',
		language TEXT DEFAULT 'en',
		theme TEXT DEFAULT 'dark',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_login DATETIME,
		failed_login_attempts INTEGER DEFAULT 0,
		locked_until DATETIME,
		metadata TEXT
	);

	-- Sessions table (exact structure required by spec)
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
		user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token TEXT UNIQUE NOT NULL,
		ip_address TEXT NOT NULL,
		user_agent TEXT,
		device_name TEXT,
		expires_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_activity DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		remember_me INTEGER DEFAULT 0,
		is_active INTEGER DEFAULT 1
	);

	-- Tokens table (exact structure required by spec)
	CREATE TABLE IF NOT EXISTS tokens (
		id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
		user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		token_hash TEXT UNIQUE NOT NULL,
		last_used DATETIME,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		revoked_at DATETIME
	);

	-- Settings table (exact structure required by spec)
	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('string', 'number', 'boolean', 'json')),
		category TEXT NOT NULL,
		description TEXT,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_by TEXT REFERENCES users(id)
	);

	-- Audit log table (exact structure required by spec)
	CREATE TABLE IF NOT EXISTS audit_log (
		id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
		user_id TEXT REFERENCES users(id),
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

	-- Scheduled tasks table (exact structure required by spec)
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
	CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
	CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);
	CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
	CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);
	CREATE INDEX IF NOT EXISTS idx_tokens_token_hash ON tokens(token_hash);
	CREATE INDEX IF NOT EXISTS idx_audit_log_user_id ON audit_log(user_id);
	CREATE INDEX IF NOT EXISTS idx_audit_log_timestamp ON audit_log(timestamp);
	CREATE INDEX IF NOT EXISTS idx_settings_category ON settings(category);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Insert default settings
	if err := insertDefaultSettings(db); err != nil {
		return fmt.Errorf("failed to insert default settings: %w", err)
	}

	return nil
}

// insertDefaultSettings adds default server settings
func insertDefaultSettings(db *sql.DB) error {
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
		{"server.timezone", "UTC", "string", "server", "Server timezone"},
		{"server.date_format", "US", "string", "server", "Date format (US, EU, ISO)"},
		{"server.time_format", "12-hour", "string", "server", "Time format (12-hour, 24-hour)"},
		{"security.session_timeout", "43200", "number", "security", "Session timeout in minutes (30 days)"},
		{"security.max_login_attempts", "5", "number", "security", "Maximum login attempts before lockout"},
		{"security.lockout_duration", "15", "number", "security", "Lockout duration in minutes"},
		{"security.password_min_length", "8", "number", "security", "Minimum password length"},
		{"features.registration_enabled", "false", "boolean", "features", "Allow user registration"},
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

// GetDB wraps the database connection with helper methods
type AppDB struct {
	*DB
	conn *sql.DB
}

// NewAppDB creates a new application database wrapper
func NewAppDB(dbPath string) (*AppDB, error) {
	// Initialize zipcode database
	zipcodeDB, err := Initialize(dbPath)
	if err != nil {
		return nil, err
	}

	// Initialize admin-only auth schema on same database
	if err := InitializeAdminSchema(zipcodeDB.conn); err != nil {
		return nil, err
	}

	return &AppDB{
		DB:   zipcodeDB,
		conn: zipcodeDB.conn,
	}, nil
}

// GetConn returns the underlying database connection
func (db *AppDB) GetConn() *sql.DB {
	return db.conn
}
