package geoip

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	// sapics/ip-location-db databases via jsdelivr CDN (daily updates)
	cityIPv4URL  = "https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv4.mmdb"
	cityIPv6URL  = "https://cdn.jsdelivr.net/npm/@ip-location-db/geolite2-city-mmdb/geolite2-city-ipv6.mmdb"
	countryURL   = "https://cdn.jsdelivr.net/npm/@ip-location-db/geo-whois-asn-country-mmdb/geo-whois-asn-country.mmdb"
	asnURL       = "https://cdn.jsdelivr.net/npm/@ip-location-db/asn-mmdb/asn.mmdb"
	defaultTimeout = 300 * time.Second // 5 minutes for large downloads
)

// DatabaseFiles holds paths to downloaded database files
type DatabaseFiles struct {
	CityIPv4DB string
	CityIPv6DB string
	CountryDB  string
	ASNDB      string
}

// DownloadDatabases downloads the latest GeoIP databases from sapics/ip-location-db via jsdelivr CDN
func DownloadDatabases(dataDir string) (*DatabaseFiles, error) {
	// Create data directory if it doesn't exist
	geoipDir := filepath.Join(dataDir, "geoip")
	if err := os.MkdirAll(geoipDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create geoip directory: %w", err)
	}

	dbFiles := &DatabaseFiles{
		CityIPv4DB: filepath.Join(geoipDir, "geolite2-city-ipv4.mmdb"),
		CityIPv6DB: filepath.Join(geoipDir, "geolite2-city-ipv6.mmdb"),
		CountryDB:  filepath.Join(geoipDir, "geo-whois-asn-country.mmdb"),
		ASNDB:      filepath.Join(geoipDir, "asn.mmdb"),
	}

	databases := map[string]string{
		dbFiles.CityIPv4DB: cityIPv4URL,
		dbFiles.CityIPv6DB: cityIPv6URL,
		dbFiles.CountryDB:  countryURL,
		dbFiles.ASNDB:      asnURL,
	}

	// Download each database
	for dbPath, url := range databases {
		filename := filepath.Base(dbPath)
		fmt.Printf("Downloading %s...\n", filename)
		if err := downloadFile(url, dbPath); err != nil {
			return nil, fmt.Errorf("failed to download %s: %w", filename, err)
		}
		fmt.Printf("Downloaded: %s\n", filename)
	}

	return dbFiles, nil
}

// downloadFile downloads a file from a URL and saves it to the specified path
func downloadFile(url, filepath string) error {
	client := &http.Client{Timeout: defaultTimeout}

	// Download file
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create output file
	outFile, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Copy data
	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// CheckForUpdates checks if there are newer databases available
// Note: sapics databases are updated daily via jsdelivr CDN
func CheckForUpdates(currentVersion string) (bool, string, error) {
	// sapics databases are updated daily, so we can't easily check versions
	// Return false for now - caller should update based on time intervals
	return false, currentVersion, nil
}

// GetDatabasePaths returns the paths to the database files
func GetDatabasePaths(dataDir string) *DatabaseFiles {
	geoipDir := filepath.Join(dataDir, "geoip")
	return &DatabaseFiles{
		CityIPv4DB: filepath.Join(geoipDir, "geolite2-city-ipv4.mmdb"),
		CityIPv6DB: filepath.Join(geoipDir, "geolite2-city-ipv6.mmdb"),
		CountryDB:  filepath.Join(geoipDir, "geo-whois-asn-country.mmdb"),
		ASNDB:      filepath.Join(geoipDir, "asn.mmdb"),
	}
}

// DatabasesExist checks if all required databases exist
func DatabasesExist(dataDir string) bool {
	paths := GetDatabasePaths(dataDir)

	cityIPv4Exists := fileExists(paths.CityIPv4DB)
	cityIPv6Exists := fileExists(paths.CityIPv6DB)
	countryExists := fileExists(paths.CountryDB)
	asnExists := fileExists(paths.ASNDB)

	// At least one city database and country database should exist, plus ASN
	return (cityIPv4Exists || cityIPv6Exists) && countryExists && asnExists
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
