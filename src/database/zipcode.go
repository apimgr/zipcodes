package database

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// Zipcode represents a US zipcode record
type Zipcode struct {
	Zip       string      `json:"zip"`
	State     string      `json:"state"`
	StateFull string      `json:"state_full"`
	City      string      `json:"city"`
	County    string      `json:"county"`
	ZipCode   int         `json:"zip_code"`
	Latitude  interface{} `json:"latitude"`
	Longitude interface{} `json:"longitude"`
	Timezone  string      `json:"timezone"`
	AreaCodes []string    `json:"area_codes"`
}

// MarshalJSON emits ZIP codes and coordinates in their public API forms.
func (z Zipcode) MarshalJSON() ([]byte, error) {
	type zipcodeJSON struct {
		Zip       string   `json:"zip"`
		State     string   `json:"state"`
		StateFull string   `json:"state_full"`
		City      string   `json:"city"`
		County    string   `json:"county"`
		ZipCode   int      `json:"zip_code"`
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
		Timezone  string   `json:"timezone"`
		AreaCodes []string `json:"area_codes"`
	}
	zip := z.Zip
	if zip == "" {
		zip = fmt.Sprintf("%05d", z.ZipCode)
	}
	latitude := coordinatePointer(z.Latitude)
	longitude := coordinatePointer(z.Longitude)
	areaCodes := z.AreaCodes
	if areaCodes == nil {
		areaCodes = []string{}
	}
	return json.Marshal(zipcodeJSON{zip, z.State, z.StateFull, z.City, z.County, z.ZipCode, latitude, longitude, z.Timezone, areaCodes})
}

func coordinatePointer(value interface{}) *float64 {
	if value == nil {
		return nil
	}
	coordinate, err := strconv.ParseFloat(fmt.Sprint(value), 64)
	if err != nil {
		return nil
	}
	return &coordinate
}

// GetLatitudeString returns latitude as string
func (z *Zipcode) GetLatitudeString() string {
	if z.Latitude == nil {
		return ""
	}
	switch v := z.Latitude.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.6f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetLongitudeString returns longitude as string
func (z *Zipcode) GetLongitudeString() string {
	if z.Longitude == nil {
		return ""
	}
	switch v := z.Longitude.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.6f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// DB holds the database connection
type DB struct {
	conn        *sql.DB
	rateLimitMu sync.Mutex
}

// Initialize creates and initializes the database
func Initialize(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	conn.SetMaxOpenConns(5)
	conn.SetMaxIdleConns(2)
	conn.SetConnMaxLifetime(5 * time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db := &DB{conn: conn}

	// Create schema
	if err := db.createSchema(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return db, nil
}

// createSchema creates the database tables
func (db *DB) createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS dataset_metadata (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS rate_limits (
		client_key TEXT NOT NULL,
		class TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		expires_at INTEGER NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_rate_limits_key_class_time
		ON rate_limits(client_key, class, created_at);
	CREATE INDEX IF NOT EXISTS idx_rate_limits_expiry
		ON rate_limits(expires_at);

	CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		type TEXT NOT NULL DEFAULT 'string',
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	);
	CREATE INDEX IF NOT EXISTS idx_config_key_prefix ON config(key);
	CREATE TABLE IF NOT EXISTS config_meta (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		version INTEGER NOT NULL DEFAULT 1,
		updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	);
	INSERT OR IGNORE INTO config_meta(id, version) VALUES(1, 1);
	-- SQLite permits one change event per trigger, so config version changes use three triggers.
	CREATE TRIGGER IF NOT EXISTS config_version_bump_insert
	AFTER INSERT ON config
	BEGIN
		UPDATE config_meta SET version = version + 1, updated_at = strftime('%s', 'now') WHERE id = 1;
	END;
	CREATE TRIGGER IF NOT EXISTS config_version_bump_update
	AFTER UPDATE ON config
	BEGIN
		UPDATE config_meta SET version = version + 1, updated_at = strftime('%s', 'now') WHERE id = 1;
	END;
	CREATE TRIGGER IF NOT EXISTS config_version_bump_delete
	AFTER DELETE ON config
	BEGIN
		UPDATE config_meta SET version = version + 1, updated_at = strftime('%s', 'now') WHERE id = 1;
	END;
	CREATE TABLE IF NOT EXISTS audit_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
		level TEXT NOT NULL DEFAULT 'info',
		category TEXT NOT NULL,
		action TEXT NOT NULL,
		actor_ip TEXT,
		target_type TEXT,
		target_id TEXT,
		details TEXT,
		success INTEGER NOT NULL DEFAULT 1
	);
	CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON audit_log(timestamp);
	CREATE INDEX IF NOT EXISTS idx_audit_category ON audit_log(category);
	CREATE INDEX IF NOT EXISTS idx_audit_actor_ip ON audit_log(actor_ip);
	CREATE INDEX IF NOT EXISTS idx_audit_target ON audit_log(target_type, target_id);

	CREATE TABLE IF NOT EXISTS scheduler_tasks (
		name TEXT PRIMARY KEY,
		interval_seconds INTEGER NOT NULL,
		next_run INTEGER NOT NULL,
		enabled INTEGER NOT NULL DEFAULT 1,
		updated_at INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS scheduler_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		started_at INTEGER NOT NULL,
		finished_at INTEGER NOT NULL,
		status TEXT NOT NULL,
		error TEXT NOT NULL DEFAULT ''
	);
	CREATE INDEX IF NOT EXISTS idx_scheduler_history_name_time
		ON scheduler_history(name, started_at);

	CREATE TABLE IF NOT EXISTS zipcodes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		state TEXT NOT NULL,
		city TEXT NOT NULL,
		county TEXT,
		zip_code INTEGER NOT NULL UNIQUE,
		latitude TEXT,
		longitude TEXT,
		state_full TEXT NOT NULL DEFAULT '',
		timezone TEXT NOT NULL DEFAULT '',
		area_codes TEXT NOT NULL DEFAULT '[]',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_zip_code ON zipcodes(zip_code);
	CREATE INDEX IF NOT EXISTS idx_city ON zipcodes(city);
	CREATE INDEX IF NOT EXISTS idx_state ON zipcodes(state);
	CREATE INDEX IF NOT EXISTS idx_state_city ON zipcodes(state, city);
	CREATE INDEX IF NOT EXISTS idx_zipcodes_latitude_numeric ON zipcodes(CAST(latitude AS REAL));
	CREATE INDEX IF NOT EXISTS idx_zipcodes_longitude_numeric ON zipcodes(CAST(longitude AS REAL));
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := db.conn.ExecContext(ctx, schema); err != nil {
		return err
	}
	for _, statement := range []string{
		`ALTER TABLE zipcodes ADD COLUMN state_full TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE zipcodes ADD COLUMN timezone TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE zipcodes ADD COLUMN area_codes TEXT NOT NULL DEFAULT '[]'`,
	} {
		if _, err := db.conn.ExecContext(ctx, statement); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column") {
			return err
		}
	}
	return nil
}

// LoadFromJSON loads zipcode data from embedded JSON bytes
func (db *DB) LoadFromJSON(data []byte) error {

	// Parse JSON
	var zipcodes []Zipcode
	if err := json.Unmarshal(data, &zipcodes); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	if len(zipcodes) == 0 {
		return fmt.Errorf("refusing to replace ZIP data with an empty dataset")
	}
	seenCodes := make(map[int]struct{}, len(zipcodes))
	for index := range zipcodes {
		code := zipcodes[index].Zip
		if code == "" {
			code = fmt.Sprintf("%05d", zipcodes[index].ZipCode)
		}
		if len(code) != 5 {
			return fmt.Errorf("invalid ZIP code at dataset index %d", index)
		}
		for _, digit := range code {
			if digit < '0' || digit > '9' {
				return fmt.Errorf("invalid ZIP code at dataset index %d", index)
			}
		}
		codeValue, err := strconv.Atoi(code)
		if err != nil || codeValue <= 0 || (zipcodes[index].ZipCode != 0 && zipcodes[index].ZipCode != codeValue) {
			return fmt.Errorf("invalid ZIP code at dataset index %d", index)
		}
		if _, exists := seenCodes[codeValue]; exists {
			return fmt.Errorf("duplicate ZIP code at dataset index %d", index)
		}
		seenCodes[codeValue] = struct{}{}
		zipcodes[index].ZipCode = codeValue
	}
	digest := sha256.Sum256(data)
	datasetHash := hex.EncodeToString(digest[:])
	var loadedHash string
	readCtx, readCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer readCancel()
	if err := db.conn.QueryRowContext(readCtx, "SELECT value FROM dataset_metadata WHERE key = ?", "zipcodes_sha256").Scan(&loadedHash); err == nil && loadedHash == datasetHash {
		return nil
	}

	// Bound the complete static data import transaction.
	loadCtx, loadCancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer loadCancel()
	tx, err := db.conn.BeginTx(loadCtx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec(`CREATE TEMP TABLE current_dataset_zipcodes (zip_code INTEGER PRIMARY KEY)`); err != nil {
		return fmt.Errorf("failed to stage dataset revision: %w", err)
	}
	stageStmt, err := tx.Prepare(`INSERT OR IGNORE INTO current_dataset_zipcodes(zip_code) VALUES(?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare dataset staging: %w", err)
	}
	defer func() { _ = stageStmt.Close() }()

	// Prepare statement
	stmt, err := tx.Prepare(`
		INSERT INTO zipcodes (state, city, county, zip_code, latitude, longitude, state_full, timezone, area_codes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(zip_code) DO UPDATE SET state=excluded.state, city=excluded.city,
		county=excluded.county, latitude=excluded.latitude, longitude=excluded.longitude,
		state_full=excluded.state_full, timezone=excluded.timezone, area_codes=excluded.area_codes
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	// Insert data
	for i, zc := range zipcodes {
		if _, err := stageStmt.Exec(zc.ZipCode); err != nil {
			return fmt.Errorf("failed to stage ZIP at dataset index %d: %w", i, err)
		}
		if zc.AreaCodes == nil {
			zc.AreaCodes = []string{}
		}
		areaCodes, err := json.Marshal(zc.AreaCodes)
		if err != nil {
			return fmt.Errorf("failed to encode area codes at index %d: %w", i, err)
		}
		_, err = stmt.Exec(zc.State, zc.City, zc.County, zc.ZipCode, zc.GetLatitudeString(), zc.GetLongitudeString(), zc.StateFull, zc.Timezone, string(areaCodes))
		if err != nil {
			return fmt.Errorf("failed to insert zipcode at index %d: %w", i, err)
		}

		if (i+1)%10000 == 0 {
			fmt.Printf("Loaded %d zipcodes...\n", i+1)
		}
	}

	if _, err := tx.Exec(`DELETE FROM zipcodes
		WHERE NOT EXISTS (SELECT 1 FROM current_dataset_zipcodes WHERE current_dataset_zipcodes.zip_code = zipcodes.zip_code)`); err != nil {
		return fmt.Errorf("failed to remove records absent from the current dataset: %w", err)
	}
	if _, err := tx.Exec(`DROP TABLE current_dataset_zipcodes`); err != nil {
		return fmt.Errorf("failed to finish dataset staging: %w", err)
	}
	if _, err := tx.Exec(`INSERT INTO dataset_metadata(key, value) VALUES(?, ?)
		ON CONFLICT(key) DO UPDATE SET value=excluded.value`, "zipcodes_sha256", datasetHash); err != nil {
		return fmt.Errorf("failed to record dataset revision: %w", err)
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
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var zc Zipcode
	var lat, lon, areaCodes string
	err := db.conn.QueryRowContext(queryCtx, `
		SELECT state, city, county, zip_code, latitude, longitude, state_full, timezone, area_codes
		FROM zipcodes WHERE zip_code = ?
	`, zipCode).Scan(&zc.State, &zc.City, &zc.County, &zc.ZipCode, &lat, &lon, &zc.StateFull, &zc.Timezone, &areaCodes)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	zc.Latitude = lat
	zc.Longitude = lon
	zc.Zip = fmt.Sprintf("%05d", zc.ZipCode)
	_ = json.Unmarshal([]byte(areaCodes), &zc.AreaCodes)
	return &zc, nil
}

// SearchByCity finds zipcodes by city name
func (db *DB) SearchByCity(city string) ([]Zipcode, error) {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.conn.QueryContext(queryCtx, `
		SELECT state, city, county, zip_code, latitude, longitude, state_full, timezone, area_codes
		FROM zipcodes WHERE LOWER(city) = LOWER(?)
		ORDER BY state, zip_code
	`, city)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	return db.scanZipcodes(rows)
}

// SearchByState finds zipcodes by state
func (db *DB) SearchByState(state string) ([]Zipcode, error) {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.conn.QueryContext(queryCtx, `
		SELECT state, city, county, zip_code, latitude, longitude, state_full, timezone, area_codes
		FROM zipcodes WHERE UPPER(state) = UPPER(?)
		ORDER BY city, zip_code
	`, state)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	return db.scanZipcodes(rows)
}

// SearchByStateAndCity finds zipcodes by state and city
func (db *DB) SearchByStateAndCity(state, city string) ([]Zipcode, error) {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.conn.QueryContext(queryCtx, `
		SELECT state, city, county, zip_code, latitude, longitude, state_full, timezone, area_codes
		FROM zipcodes WHERE UPPER(state) = UPPER(?) AND LOWER(city) = LOWER(?)
		ORDER BY zip_code
	`, state, city)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	return db.scanZipcodes(rows)
}

// SearchByPrefix finds zipcodes by prefix (e.g., "94" matches 94000-94999)
func (db *DB) SearchByPrefix(prefix string) ([]Zipcode, error) {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.conn.QueryContext(queryCtx, `
		SELECT state, city, county, zip_code, latitude, longitude, state_full, timezone, area_codes
		FROM zipcodes WHERE printf('%05d', zip_code) LIKE ?
		ORDER BY zip_code
	`, prefix+"%")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	return db.scanZipcodes(rows)
}

// AutoComplete provides autocomplete suggestions
func (db *DB) AutoComplete(query string, limit int) ([]string, error) {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if limit <= 0 {
		limit = 10
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return []string{}, nil
	}

	rows, err := db.conn.QueryContext(queryCtx, `
		SELECT DISTINCT city || ', ' || state as suggestion
		FROM zipcodes
		WHERE LOWER(city) LIKE LOWER(?) OR UPPER(state) LIKE UPPER(?)
		ORDER BY city
		LIMIT ?
	`, query+"%", query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

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
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stats := make(map[string]interface{})

	// Total zipcodes
	var total int
	err := db.conn.QueryRowContext(queryCtx, "SELECT COUNT(*) FROM zipcodes").Scan(&total)
	if err != nil {
		return nil, err
	}
	stats["total_zipcodes"] = total

	// Total states
	var states int
	err = db.conn.QueryRowContext(queryCtx, "SELECT COUNT(DISTINCT state) FROM zipcodes").Scan(&states)
	if err != nil {
		return nil, err
	}
	stats["total_states"] = states

	// Total cities
	var cities int
	err = db.conn.QueryRowContext(queryCtx, "SELECT COUNT(DISTINCT city) FROM zipcodes").Scan(&cities)
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
		var lat, lon, areaCodes string
		if err := rows.Scan(&zc.State, &zc.City, &zc.County, &zc.ZipCode, &lat, &lon, &zc.StateFull, &zc.Timezone, &areaCodes); err != nil {
			return nil, err
		}
		zc.Latitude = lat
		zc.Longitude = lon
		zc.Zip = fmt.Sprintf("%05d", zc.ZipCode)
		_ = json.Unmarshal([]byte(areaCodes), &zc.AreaCodes)
		zipcodes = append(zipcodes, zc)
	}
	return zipcodes, rows.Err()
}

// GetRandom returns a random zipcode
func (db *DB) GetRandom() (*Zipcode, error) {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var zc Zipcode
	var lat, lon, areaCodes string
	err := db.conn.QueryRowContext(queryCtx, `
		SELECT state, city, county, zip_code, latitude, longitude, state_full, timezone, area_codes
		FROM zipcodes ORDER BY RANDOM() LIMIT 1
	`).Scan(&zc.State, &zc.City, &zc.County, &zc.ZipCode, &lat, &lon, &zc.StateFull, &zc.Timezone, &areaCodes)

	if err != nil {
		return nil, err
	}

	zc.Latitude = lat
	zc.Longitude = lon
	zc.Zip = fmt.Sprintf("%05d", zc.ZipCode)
	_ = json.Unmarshal([]byte(areaCodes), &zc.AreaCodes)
	return &zc, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// RecordRateLimitHit records abuse-control events without storing a raw IP address.
func (db *DB) RecordRateLimitHit(clientKey, class string, now time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.conn.ExecContext(ctx, `INSERT INTO audit_log(timestamp, level, category, action, actor_ip, target_type, target_id, details, success)
		VALUES(?, 'warning', 'security', 'rate_limit_hit', ?, 'rate_limit', ?, ?, 0)`,
		now.Unix(), clientKey, class, `{"class":"`+class+`"}`)
	return err
}

// ConsumeRateLimit records a request in the persistent sliding-window store.
func (db *DB) ConsumeRateLimit(clientKey, class string, limit int, window time.Duration, now time.Time) (time.Duration, error) {
	db.rateLimitMu.Lock()
	defer db.rateLimitMu.Unlock()
	if limit <= 0 || window <= 0 {
		return 0, fmt.Errorf("rate limit and window must be positive")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()
	nowUnix := now.UnixNano()
	if _, err := tx.ExecContext(ctx, `DELETE FROM rate_limits WHERE expires_at <= ?`, nowUnix); err != nil {
		return 0, err
	}
	cutoff := now.Add(-window).UnixNano()
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM rate_limits WHERE client_key = ? AND class = ? AND created_at > ?`, clientKey, class, cutoff).Scan(&count); err != nil {
		return 0, err
	}
	if count >= limit {
		var oldest int64
		if err := tx.QueryRowContext(ctx, `SELECT MIN(created_at) FROM rate_limits WHERE client_key = ? AND class = ? AND created_at > ?`, clientKey, class, cutoff).Scan(&oldest); err != nil {
			return 0, err
		}
		if err := tx.Commit(); err != nil {
			return 0, err
		}
		return time.Duration(oldest + int64(window) - nowUnix), nil
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO rate_limits(client_key, class, created_at, expires_at) VALUES(?, ?, ?, ?)`, clientKey, class, nowUnix, now.Add(window).UnixNano()); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return 0, nil
}

// ScheduledTaskState stores the persisted timing for a background task.
type ScheduledTaskState struct {
	Interval time.Duration
	NextRun  time.Time
	Enabled  bool
}

// GetScheduledTask returns a persisted task state, if one exists.
func (db *DB) GetScheduledTask(name string) (ScheduledTaskState, bool, error) {
	var state ScheduledTaskState
	var intervalSeconds, nextRun int64
	var enabled int
	err := db.conn.QueryRow(`SELECT interval_seconds, next_run, enabled FROM scheduler_tasks WHERE name = ?`, name).Scan(&intervalSeconds, &nextRun, &enabled)
	if err == sql.ErrNoRows {
		return state, false, nil
	}
	if err != nil {
		return state, false, err
	}
	state.Interval = time.Duration(intervalSeconds) * time.Second
	state.NextRun = time.Unix(nextRun, 0)
	state.Enabled = enabled != 0
	return state, true, nil
}

// SaveScheduledTask persists a task's next scheduled run.
func (db *DB) SaveScheduledTask(name string, interval time.Duration, nextRun time.Time, enabled bool) error {
	enabledValue := 0
	if enabled {
		enabledValue = 1
	}
	_, err := db.conn.Exec(`INSERT INTO scheduler_tasks(name, interval_seconds, next_run, enabled, updated_at) VALUES(?, ?, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET interval_seconds=excluded.interval_seconds, next_run=excluded.next_run, enabled=excluded.enabled, updated_at=excluded.updated_at`, name, int64(interval/time.Second), nextRun.Unix(), enabledValue, time.Now().Unix())
	return err
}

// RecordScheduledTaskRun stores a bounded execution history record.
func (db *DB) RecordScheduledTaskRun(name string, started, finished time.Time, runErr error) error {
	status, message := "success", ""
	if runErr != nil {
		status, message = "error", runErr.Error()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `INSERT INTO scheduler_history(name, started_at, finished_at, status, error) VALUES(?, ?, ?, ?, ?)`, name, started.Unix(), finished.Unix(), status, message); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM scheduler_history WHERE id NOT IN (SELECT id FROM scheduler_history ORDER BY id DESC LIMIT 1000)`); err != nil {
		return err
	}
	return tx.Commit()
}
