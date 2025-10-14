package geoip

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	githubAPIURL   = "https://api.github.com/repos/P3TERX/GeoLite.mmdb/releases/latest"
	githubRepoURL  = "https://github.com/P3TERX/GeoLite.mmdb"
	defaultTimeout = 300 * time.Second // 5 minutes for large downloads
)

// DatabaseFiles holds paths to downloaded database files
type DatabaseFiles struct {
	CityDB    string
	CountryDB string
	ASNDB     string
}

// Release represents a GitHub release
type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// DownloadDatabases downloads the latest GeoLite.mmdb databases from GitHub
func DownloadDatabases(dataDir string) (*DatabaseFiles, error) {
	// Create data directory if it doesn't exist
	geoipDir := filepath.Join(dataDir, "geoip")
	if err := os.MkdirAll(geoipDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create geoip directory: %w", err)
	}

	// Get latest release info
	release, err := getLatestRelease()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %w", err)
	}

	dbFiles := &DatabaseFiles{}

	// Download each database
	for _, asset := range release.Assets {
		name := asset.Name
		url := asset.BrowserDownloadURL

		// Determine database type
		var dbPath string
		if strings.Contains(name, "GeoLite2-City") {
			dbPath = filepath.Join(geoipDir, "GeoLite2-City.mmdb")
		} else if strings.Contains(name, "GeoLite2-Country") {
			dbPath = filepath.Join(geoipDir, "GeoLite2-Country.mmdb")
		} else if strings.Contains(name, "GeoLite2-ASN") {
			dbPath = filepath.Join(geoipDir, "GeoLite2-ASN.mmdb")
		} else {
			continue // Skip non-database files
		}

		// Download and extract
		fmt.Printf("Downloading %s...\n", name)
		if err := downloadAndExtract(url, dbPath); err != nil {
			return nil, fmt.Errorf("failed to download %s: %w", name, err)
		}

		// Set the appropriate database path
		if strings.Contains(name, "City") {
			dbFiles.CityDB = dbPath
		} else if strings.Contains(name, "Country") {
			dbFiles.CountryDB = dbPath
		} else if strings.Contains(name, "ASN") {
			dbFiles.ASNDB = dbPath
		}

		fmt.Printf("Downloaded and extracted: %s\n", dbPath)
	}

	return dbFiles, nil
}

// getLatestRelease fetches the latest release information from GitHub
func getLatestRelease() (*Release, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", githubAPIURL, nil)
	if err != nil {
		return nil, err
	}

	// Set user agent to avoid rate limiting
	req.Header.Set("User-Agent", "zipcodes-geoip-updater")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// downloadAndExtract downloads a tar.gz file and extracts the .mmdb file
func downloadAndExtract(url, outputPath string) error {
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

	// Create gzip reader
	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzReader)

	// Extract .mmdb file
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar read error: %w", err)
		}

		// Look for .mmdb file
		if strings.HasSuffix(header.Name, ".mmdb") {
			// Create output file
			outFile, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer outFile.Close()

			// Copy data
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("no .mmdb file found in archive")
}

// CheckForUpdates checks if there are newer databases available
func CheckForUpdates(currentVersion string) (bool, string, error) {
	release, err := getLatestRelease()
	if err != nil {
		return false, "", err
	}

	// Compare versions
	if release.TagName != currentVersion {
		return true, release.TagName, nil
	}

	return false, currentVersion, nil
}

// GetDatabasePaths returns the paths to the database files
func GetDatabasePaths(dataDir string) *DatabaseFiles {
	geoipDir := filepath.Join(dataDir, "geoip")
	return &DatabaseFiles{
		CityDB:    filepath.Join(geoipDir, "GeoLite2-City.mmdb"),
		CountryDB: filepath.Join(geoipDir, "GeoLite2-Country.mmdb"),
		ASNDB:     filepath.Join(geoipDir, "GeoLite2-ASN.mmdb"),
	}
}

// DatabasesExist checks if all required databases exist
func DatabasesExist(dataDir string) bool {
	paths := GetDatabasePaths(dataDir)

	cityExists := fileExists(paths.CityDB)
	countryExists := fileExists(paths.CountryDB)
	asnExists := fileExists(paths.ASNDB)

	// At least City or Country database should exist
	return (cityExists || countryExists) && asnExists
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
