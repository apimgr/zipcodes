package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Zipcode represents a US zipcode record
type Zipcode struct {
	State     string  `json:"state"`
	City      string  `json:"city"`
	County    string  `json:"county"`
	ZipCode   int     `json:"zip_code"`
	Latitude  string  `json:"latitude"`
	Longitude string  `json:"longitude"`
}

// DB holds the database connection
type DB struct {
	conn *sql.DB
}

// Initialize creates and initializes the database
func Initialize(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}

	// Create schema
	if err := db.createSchema(); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return db, nil
}

// createSchema creates the database tables
func (db *DB) createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS zipcodes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		state TEXT NOT NULL,
		city TEXT NOT NULL,
		county TEXT,
		zip_code INTEGER NOT NULL UNIQUE,
		latitude TEXT,
		longitude TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_zip_code ON zipcodes(zip_code);
	CREATE INDEX IF NOT EXISTS idx_city ON zipcodes(city);
	CREATE INDEX IF NOT EXISTS idx_state ON zipcodes(state);
	CREATE INDEX IF NOT EXISTS idx_state_city ON zipcodes(state, city);
	`

	_, err := db.conn.Exec(schema)
	return err
}

// LoadFromJSON loads zipcode data from embedded JSON bytes
func (db *DB) LoadFromJSON(data []byte) error {
	// Check if data already loaded
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM zipcodes").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing data: %w", err)
	}

	if count > 0 {
		fmt.Printf("Database already contains %d zipcodes, skipping load\n", count)
		return nil
	}

	// Parse JSON
	var zipcodes []Zipcode
	if err := json.Unmarshal(data, &zipcodes); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Begin transaction
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare statement
	stmt, err := tx.Prepare(`
		INSERT INTO zipcodes (state, city, county, zip_code, latitude, longitude)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert data
	for i, zc := range zipcodes {
		_, err := stmt.Exec(zc.State, zc.City, zc.County, zc.ZipCode, zc.Latitude, zc.Longitude)
		if err != nil {
			return fmt.Errorf("failed to insert zipcode at index %d: %w", i, err)
		}

		if (i+1)%10000 == 0 {
			fmt.Printf("Loaded %d zipcodes...\n", i+1)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("Successfully loaded %d zipcodes\n", len(zipcodes))
	return nil
}

// SearchByZipCode finds a zipcode by its code
func (db *DB) SearchByZipCode(zipCode int) (*Zipcode, error) {
	var zc Zipcode
	err := db.conn.QueryRow(`
		SELECT state, city, county, zip_code, latitude, longitude
		FROM zipcodes WHERE zip_code = ?
	`, zipCode).Scan(&zc.State, &zc.City, &zc.County, &zc.ZipCode, &zc.Latitude, &zc.Longitude)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &zc, nil
}

// SearchByCity finds zipcodes by city name
func (db *DB) SearchByCity(city string) ([]Zipcode, error) {
	rows, err := db.conn.Query(`
		SELECT state, city, county, zip_code, latitude, longitude
		FROM zipcodes WHERE LOWER(city) = LOWER(?)
		ORDER BY state, zip_code
	`, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.scanZipcodes(rows)
}

// SearchByState finds zipcodes by state
func (db *DB) SearchByState(state string) ([]Zipcode, error) {
	rows, err := db.conn.Query(`
		SELECT state, city, county, zip_code, latitude, longitude
		FROM zipcodes WHERE UPPER(state) = UPPER(?)
		ORDER BY city, zip_code
		LIMIT 1000
	`, state)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.scanZipcodes(rows)
}

// SearchByStateAndCity finds zipcodes by state and city
func (db *DB) SearchByStateAndCity(state, city string) ([]Zipcode, error) {
	rows, err := db.conn.Query(`
		SELECT state, city, county, zip_code, latitude, longitude
		FROM zipcodes WHERE UPPER(state) = UPPER(?) AND LOWER(city) = LOWER(?)
		ORDER BY zip_code
	`, state, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.scanZipcodes(rows)
}

// SearchByPrefix finds zipcodes by prefix (e.g., "94" matches 94000-94999)
func (db *DB) SearchByPrefix(prefix string) ([]Zipcode, error) {
	rows, err := db.conn.Query(`
		SELECT state, city, county, zip_code, latitude, longitude
		FROM zipcodes WHERE CAST(zip_code AS TEXT) LIKE ?
		ORDER BY zip_code
		LIMIT 500
	`, prefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.scanZipcodes(rows)
}

// AutoComplete provides autocomplete suggestions
func (db *DB) AutoComplete(query string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return []string{}, nil
	}

	rows, err := db.conn.Query(`
		SELECT DISTINCT city || ', ' || state as suggestion
		FROM zipcodes
		WHERE LOWER(city) LIKE LOWER(?) OR UPPER(state) LIKE UPPER(?)
		ORDER BY city
		LIMIT ?
	`, query+"%", query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suggestions []string
	for rows.Next() {
		var suggestion string
		if err := rows.Scan(&suggestion); err != nil {
			return nil, err
		}
		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}

// GetStats returns database statistics
func (db *DB) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total zipcodes
	var total int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM zipcodes").Scan(&total)
	if err != nil {
		return nil, err
	}
	stats["total_zipcodes"] = total

	// Total states
	var states int
	err = db.conn.QueryRow("SELECT COUNT(DISTINCT state) FROM zipcodes").Scan(&states)
	if err != nil {
		return nil, err
	}
	stats["total_states"] = states

	// Total cities
	var cities int
	err = db.conn.QueryRow("SELECT COUNT(DISTINCT city) FROM zipcodes").Scan(&cities)
	if err != nil {
		return nil, err
	}
	stats["total_cities"] = cities

	return stats, nil
}

// scanZipcodes is a helper to scan multiple zipcode rows
func (db *DB) scanZipcodes(rows *sql.Rows) ([]Zipcode, error) {
	var zipcodes []Zipcode
	for rows.Next() {
		var zc Zipcode
		if err := rows.Scan(&zc.State, &zc.City, &zc.County, &zc.ZipCode, &zc.Latitude, &zc.Longitude); err != nil {
			return nil, err
		}
		zipcodes = append(zipcodes, zc)
	}
	return zipcodes, rows.Err()
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}
