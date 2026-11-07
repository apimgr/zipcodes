package geoip

import (
	"fmt"
	"net"
	"sync"

	"github.com/oschwald/maxminddb-golang"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// GeoIP owns readers for the CC BY 4.0 ip-location-db databases.
type GeoIP struct {
	cityIPv4DB *maxminddb.Reader
	cityIPv6DB *maxminddb.Reader
	countryDB *maxminddb.Reader
	asnDB *maxminddb.Reader
	mu sync.RWMutex
}

type cityRecord struct {
	City string `maxminddb:"city"`
	CountryCode string `maxminddb:"country_code"`
	State1 string `maxminddb:"state1"`
	State2 string `maxminddb:"state2"`
	Postcode string `maxminddb:"postcode"`
	Latitude float64 `maxminddb:"latitude"`
	Longitude float64 `maxminddb:"longitude"`
	Timezone string `maxminddb:"timezone"`
}

type countryRecord struct {
	CountryCode string `maxminddb:"country_code"`
}

type asnRecord struct {
	Number uint `maxminddb:"autonomous_system_number"`
	Organization string `maxminddb:"autonomous_system_organization"`
}

// Location represents the available approximate location for an IP address.
type Location struct {
	IP string `json:"ip"`
	Country string `json:"country"`
	CountryCode string `json:"country_code"`
	City string `json:"city"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone string `json:"timezone"`
	ASN uint `json:"asn,omitempty"`
	ASNOrg string `json:"asn_org,omitempty"`
}

var (
	instance *GeoIP
	instanceMu sync.RWMutex
)

// Initialize opens the configured MMDB files and replaces the active readers.
func Initialize(cityIPv4DBPath, cityIPv6DBPath, countryDBPath, asnDBPath string) error {
	readers := make([]*maxminddb.Reader, 4)
	paths := []string{cityIPv4DBPath, cityIPv6DBPath, countryDBPath, asnDBPath}
	for index, path := range paths {
		if path == "" { continue }
		reader, err := maxminddb.Open(path)
		if err != nil {
			for _, opened := range readers { if opened != nil { _ = opened.Close() } }
			return fmt.Errorf("failed to open GeoIP database %q: %w", path, err)
		}
		readers[index] = reader
	}
	if readers[0] == nil && readers[1] == nil && readers[2] == nil && readers[3] == nil {
		return fmt.Errorf("no GeoIP databases were configured")
	}
	initialized := &GeoIP{cityIPv4DB: readers[0], cityIPv6DB: readers[1], countryDB: readers[2], asnDB: readers[3]}
	instanceMu.Lock()
	previous := instance
	instance = initialized
	instanceMu.Unlock()
	if previous != nil { _ = previous.Close() }
	return nil
}

// GetInstance returns the active GeoIP reader set.
func GetInstance() *GeoIP {
	instanceMu.RLock()
	defer instanceMu.RUnlock()
	return instance
}

// Lookup queries the local IP-location databases without making network requests.
func (g *GeoIP) Lookup(ip string) (*Location, error) {
	if g == nil { return nil, fmt.Errorf("GeoIP not initialized") }
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil { return nil, fmt.Errorf("invalid IP address") }
	if !isPublicAddress(parsedIP) { return &Location{IP: parsedIP.String()}, nil }
	g.mu.RLock()
	defer g.mu.RUnlock()
	location := &Location{IP: parsedIP.String()}
	var cityDB *maxminddb.Reader
	if parsedIP.To4() != nil { cityDB = g.cityIPv4DB } else { cityDB = g.cityIPv6DB }
	if cityDB != nil {
		var record cityRecord
		if err := cityDB.Lookup(parsedIP, &record); err == nil {
			location.City = record.City
			location.CountryCode = record.CountryCode
			location.Latitude = record.Latitude
			location.Longitude = record.Longitude
			location.Timezone = record.Timezone
		}
	}
	if location.CountryCode == "" && g.countryDB != nil {
		var record countryRecord
		if err := g.countryDB.Lookup(parsedIP, &record); err == nil { location.CountryCode = record.CountryCode }
	}
	if location.CountryCode != "" {
		if region, err := language.ParseRegion(location.CountryCode); err == nil {
			location.Country = display.English.Regions().Name(region)
		}
	}
	if g.asnDB != nil {
		var record asnRecord
		if err := g.asnDB.Lookup(parsedIP, &record); err == nil {
			location.ASN = record.Number
			location.ASNOrg = record.Organization
		}
	}
	return location, nil
}

// Reload replaces all active database readers with the updated files.
func (g *GeoIP) Reload(cityIPv4DBPath, cityIPv6DBPath, countryDBPath, asnDBPath string) error {
	return Initialize(cityIPv4DBPath, cityIPv6DBPath, countryDBPath, asnDBPath)
}

// Close closes all database readers owned by this instance.
func (g *GeoIP) Close() error {
	if g == nil { return nil }
	g.mu.Lock()
	defer g.mu.Unlock()
	var firstErr error
	for _, reader := range []*maxminddb.Reader{g.cityIPv4DB, g.cityIPv6DB, g.countryDB, g.asnDB} {
		if reader != nil {
			if err := reader.Close(); err != nil && firstErr == nil { firstErr = err }
		}
	}
	return firstErr
}

// LookupIP is a convenience function to lookup an IP using the global instance
func LookupIP(ip string) (*Location, error) {
	instanceMu.RLock()
	defer instanceMu.RUnlock()
	if instance == nil { return nil, fmt.Errorf("GeoIP not initialized") }
	return instance.Lookup(ip)
}
