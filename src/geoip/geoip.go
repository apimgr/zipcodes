package geoip

import (
	"fmt"
	"net"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

// GeoIP holds the GeoIP2 database readers
type GeoIP struct {
	cityIPv4DB *geoip2.Reader // City database for IPv4 addresses
	cityIPv6DB *geoip2.Reader // City database for IPv6 addresses
	countryDB  *geoip2.Reader // Country database (combined IPv4/IPv6)
	asnDB      *geoip2.Reader // ASN database (combined IPv4/IPv6)
	mu         sync.RWMutex
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
func Initialize(cityIPv4DBPath, cityIPv6DBPath, countryDBPath, asnDBPath string) error {
	var err error
	once.Do(func() {
		instance = &GeoIP{}

		// Load City IPv4 database
		if cityIPv4DBPath != "" {
			instance.cityIPv4DB, err = geoip2.Open(cityIPv4DBPath)
			if err != nil {
				err = fmt.Errorf("failed to open city IPv4 database: %w", err)
				return
			}
		}

		// Load City IPv6 database
		if cityIPv6DBPath != "" {
			instance.cityIPv6DB, err = geoip2.Open(cityIPv6DBPath)
			if err != nil {
				err = fmt.Errorf("failed to open city IPv6 database: %w", err)
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

	// Determine which city database to use based on IP version
	var cityDB *geoip2.Reader
	if parsedIP.To4() != nil {
		// IPv4 address
		cityDB = g.cityIPv4DB
	} else {
		// IPv6 address
		cityDB = g.cityIPv6DB
	}

	// Try City database first (has most detailed info)
	if cityDB != nil {
		record, err := cityDB.City(parsedIP)
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
func (g *GeoIP) Reload(cityIPv4DBPath, cityIPv6DBPath, countryDBPath, asnDBPath string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Close existing databases
	if g.cityIPv4DB != nil {
		g.cityIPv4DB.Close()
	}
	if g.cityIPv6DB != nil {
		g.cityIPv6DB.Close()
	}
	if g.countryDB != nil {
		g.countryDB.Close()
	}
	if g.asnDB != nil {
		g.asnDB.Close()
	}

	// Reload databases
	var err error
	if cityIPv4DBPath != "" {
		g.cityIPv4DB, err = geoip2.Open(cityIPv4DBPath)
		if err != nil {
			return fmt.Errorf("failed to reload city IPv4 database: %w", err)
		}
	}

	if cityIPv6DBPath != "" {
		g.cityIPv6DB, err = geoip2.Open(cityIPv6DBPath)
		if err != nil {
			return fmt.Errorf("failed to reload city IPv6 database: %w", err)
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

	if g.cityIPv4DB != nil {
		g.cityIPv4DB.Close()
	}
	if g.cityIPv6DB != nil {
		g.cityIPv6DB.Close()
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
