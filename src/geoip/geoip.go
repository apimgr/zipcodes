package geoip

import (
	"fmt"
	"net"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

// GeoIP holds the GeoIP2 database readers
type GeoIP struct {
	cityDB    *geoip2.Reader
	countryDB *geoip2.Reader
	asnDB     *geoip2.Reader
	mu        sync.RWMutex
}

// Location represents a geographical location
type Location struct {
	IP          string  `json:"ip"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
	ASN         uint    `json:"asn,omitempty"`
	ASNOrg      string  `json:"asn_org,omitempty"`
}

var (
	instance *GeoIP
	once     sync.Once
)

// Initialize creates the GeoIP instance with database paths
func Initialize(cityDBPath, countryDBPath, asnDBPath string) error {
	var err error
	once.Do(func() {
		instance = &GeoIP{}

		// Load City database
		if cityDBPath != "" {
			instance.cityDB, err = geoip2.Open(cityDBPath)
			if err != nil {
				err = fmt.Errorf("failed to open city database: %w", err)
				return
			}
		}

		// Load Country database
		if countryDBPath != "" {
			instance.countryDB, err = geoip2.Open(countryDBPath)
			if err != nil {
				err = fmt.Errorf("failed to open country database: %w", err)
				return
			}
		}

		// Load ASN database
		if asnDBPath != "" {
			instance.asnDB, err = geoip2.Open(asnDBPath)
			if err != nil {
				err = fmt.Errorf("failed to open ASN database: %w", err)
				return
			}
		}
	})

	return err
}

// GetInstance returns the GeoIP singleton instance
func GetInstance() *GeoIP {
	return instance
}

// Lookup performs a GeoIP lookup for the given IP address
func (g *GeoIP) Lookup(ip string) (*Location, error) {
	if g == nil {
		return nil, fmt.Errorf("GeoIP not initialized")
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	location := &Location{
		IP: ip,
	}

	// Try City database first (has most detailed info)
	if g.cityDB != nil {
		record, err := g.cityDB.City(parsedIP)
		if err == nil {
			location.Country = record.Country.Names["en"]
			location.CountryCode = record.Country.IsoCode
			location.City = record.City.Names["en"]
			location.Latitude = record.Location.Latitude
			location.Longitude = record.Location.Longitude
			location.Timezone = record.Location.TimeZone
		}
	} else if g.countryDB != nil {
		// Fallback to Country database
		record, err := g.countryDB.Country(parsedIP)
		if err == nil {
			location.Country = record.Country.Names["en"]
			location.CountryCode = record.Country.IsoCode
		}
	}

	// Get ASN information if available
	if g.asnDB != nil {
		record, err := g.asnDB.ASN(parsedIP)
		if err == nil {
			location.ASN = record.AutonomousSystemNumber
			location.ASNOrg = record.AutonomousSystemOrganization
		}
	}

	return location, nil
}

// Reload reloads the GeoIP databases (for updates)
func (g *GeoIP) Reload(cityDBPath, countryDBPath, asnDBPath string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Close existing databases
	if g.cityDB != nil {
		g.cityDB.Close()
	}
	if g.countryDB != nil {
		g.countryDB.Close()
	}
	if g.asnDB != nil {
		g.asnDB.Close()
	}

	// Reload databases
	var err error
	if cityDBPath != "" {
		g.cityDB, err = geoip2.Open(cityDBPath)
		if err != nil {
			return fmt.Errorf("failed to reload city database: %w", err)
		}
	}

	if countryDBPath != "" {
		g.countryDB, err = geoip2.Open(countryDBPath)
		if err != nil {
			return fmt.Errorf("failed to reload country database: %w", err)
		}
	}

	if asnDBPath != "" {
		g.asnDB, err = geoip2.Open(asnDBPath)
		if err != nil {
			return fmt.Errorf("failed to reload ASN database: %w", err)
		}
	}

	return nil
}

// Close closes all database readers
func (g *GeoIP) Close() error {
	if g == nil {
		return nil
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	if g.cityDB != nil {
		g.cityDB.Close()
	}
	if g.countryDB != nil {
		g.countryDB.Close()
	}
	if g.asnDB != nil {
		g.asnDB.Close()
	}

	return nil
}

// LookupIP is a convenience function to lookup an IP using the global instance
func LookupIP(ip string) (*Location, error) {
	if instance == nil {
		return nil, fmt.Errorf("GeoIP not initialized")
	}
	return instance.Lookup(ip)
}
