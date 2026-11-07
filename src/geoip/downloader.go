package geoip

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/oschwald/maxminddb-golang"
)

const (
	// ip-location-db databases use the canonical CDN and release URLs.
	cityIPv4URL    = "https://github.com/sapics/ip-location-db/releases/download/latest/dbip-city-ipv4.mmdb"
	cityIPv6URL    = "https://github.com/sapics/ip-location-db/releases/download/latest/dbip-city-ipv6.mmdb"
	countryURL     = "https://cdn.jsdelivr.net/npm/@ip-location-db/geo-whois-asn-country-mmdb/geo-whois-asn-country.mmdb"
	asnURL         = "https://cdn.jsdelivr.net/npm/@ip-location-db/asn-mmdb/asn.mmdb"
	defaultTimeout = 300 * time.Second // 5 minutes for large downloads
)

// DatabaseFiles holds paths to downloaded database files
type DatabaseFiles struct {
	CityIPv4DB string
	CityIPv6DB string
	CountryDB  string
	ASNDB      string
}

// DownloadDatabases downloads the latest ip-location-db databases.
func DownloadDatabases(dataDir string) (*DatabaseFiles, error) {
	return DownloadDatabasesInDir(filepath.Join(dataDir, "security", "geoip"))
}

// DownloadDatabasesInDir downloads databases into the configured GeoIP directory.
func DownloadDatabasesInDir(geoipDir string) (*DatabaseFiles, error) {
	return DownloadSelectedDatabases(geoipDir, true, true, true)
}

// DownloadSelectedDatabases downloads only enabled database categories.
func DownloadSelectedDatabases(geoipDir string, cityEnabled, countryEnabled, asnEnabled bool) (*DatabaseFiles, error) {
	if !cityEnabled && !countryEnabled && !asnEnabled {
		return nil, fmt.Errorf("no GeoIP database categories are enabled")
	}
	if err := os.MkdirAll(geoipDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create geoip directory: %w", err)
	}

	dbFiles := &DatabaseFiles{
		CityIPv4DB: filepath.Join(geoipDir, "dbip-city-ipv4.mmdb"),
		CityIPv6DB: filepath.Join(geoipDir, "dbip-city-ipv6.mmdb"),
		CountryDB:  filepath.Join(geoipDir, "geo-whois-asn-country.mmdb"),
		ASNDB:      filepath.Join(geoipDir, "asn.mmdb"),
	}

	databases := make(map[string]string)
	if cityEnabled {
		databases[dbFiles.CityIPv4DB] = cityIPv4URL
		databases[dbFiles.CityIPv6DB] = cityIPv6URL
	}
	if countryEnabled {
		databases[dbFiles.CountryDB] = countryURL
	}
	if asnEnabled {
		databases[dbFiles.ASNDB] = asnURL
	}

	// Download each database
	for dbPath, url := range databases {
		filename := filepath.Base(dbPath)
		fmt.Printf("Downloading %s...\n", filename)
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			return nil, err
		}
		if err := downloadFile(url, dbPath); err != nil {
			return nil, fmt.Errorf("failed to download %s: %w", filename, err)
		}
		fmt.Printf("Downloaded: %s\n", filename)
	}

	return dbFiles, nil
}

// downloadFile downloads a file from a URL and saves it to the specified path
func downloadFile(url, destination string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	if resp.ContentLength > 500<<20 {
		return fmt.Errorf("download exceeds maximum allowed size")
	}
	outFile, err := os.CreateTemp(filepath.Dir(destination), ".geoip-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	temporaryPath := outFile.Name()
	defer func() { _ = os.Remove(temporaryPath) }()
	if err := outFile.Chmod(0600); err != nil {
		_ = outFile.Close()
		return err
	}
	written, err := io.Copy(outFile, io.LimitReader(resp.Body, (500<<20)+1))
	if err != nil {
		_ = outFile.Close()
		return fmt.Errorf("failed to write file: %w", err)
	}
	if written > 500<<20 {
		_ = outFile.Close()
		return fmt.Errorf("download exceeds maximum allowed size")
	}
	if err := outFile.Sync(); err != nil {
		_ = outFile.Close()
		return err
	}
	if err := outFile.Close(); err != nil {
		return err
	}
	backupPath := temporaryPath + ".previous"
	hadPrevious := false
	if err := os.Rename(destination, backupPath); err == nil {
		hadPrevious = true
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to preserve existing database: %w", err)
	}
	if err := os.Rename(temporaryPath, destination); err != nil {
		if hadPrevious {
			_ = os.Rename(backupPath, destination)
		}
		return fmt.Errorf("failed to install downloaded file: %w", err)
	}
	if hadPrevious {
		_ = os.Remove(backupPath)
	}
	return nil
}

// GetDatabasePaths returns the paths to the database files
func GetDatabasePaths(dataDir string) *DatabaseFiles {
	return GetDatabasePathsInDir(filepath.Join(dataDir, "security", "geoip"))
}

// GetDatabasePathsInDir returns paths under the configured GeoIP directory.
func GetDatabasePathsInDir(geoipDir string) *DatabaseFiles {
	return &DatabaseFiles{
		CityIPv4DB: filepath.Join(geoipDir, "dbip-city-ipv4.mmdb"),
		CityIPv6DB: filepath.Join(geoipDir, "dbip-city-ipv6.mmdb"),
		CountryDB:  filepath.Join(geoipDir, "geo-whois-asn-country.mmdb"),
		ASNDB:      filepath.Join(geoipDir, "asn.mmdb"),
	}
}

// DatabasesExist checks if all required databases exist
func DatabasesExist(dataDir string) bool {
	return DatabasesExistInDir(filepath.Join(dataDir, "security", "geoip"))
}

// DatabasesExistInDir reports whether all configured databases exist.
func DatabasesExistInDir(geoipDir string) bool {
	return DatabasesExistInDirFor(geoipDir, true, true, true)
}

// DatabasesExistInDirFor checks whether all enabled categories have usable files.
func DatabasesExistInDirFor(geoipDir string, cityEnabled, countryEnabled, asnEnabled bool) bool {
	paths := GetDatabasePathsInDir(geoipDir)
	if cityEnabled && (!validMMDB(paths.CityIPv4DB) || !validMMDB(paths.CityIPv6DB)) {
		return false
	}
	if countryEnabled && !validMMDB(paths.CountryDB) {
		return false
	}
	if asnEnabled && !validMMDB(paths.ASNDB) {
		return false
	}
	return cityEnabled || countryEnabled || asnEnabled
}

// validMMDB confirms a required database can be opened before it is reused.
func validMMDB(path string) bool {
	reader, err := maxminddb.Open(path)
	if err != nil {
		return false
	}
	return reader.Close() == nil
}
