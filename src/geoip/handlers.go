package geoip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"
)

// MeHandler returns a request-scoped lookup for the caller's address.
func MeHandler(w http.ResponseWriter, r *http.Request) {
	ip := getClientIP(r)
	location := &Location{IP: ip}
	if parsed := net.ParseIP(ip); parsed != nil && isPublicAddress(parsed) {
		if found, err := LookupIP(ip); err == nil && found != nil {
			location = found
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(map[string]interface{}{"ok": true, "data": location})
}

// LookupHandler handles GeoIP lookup requests.
func LookupHandler(w http.ResponseWriter, r *http.Request) {
	// Get IP from query parameter or use client IP
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		ip = getClientIP(r)
	}

	// Perform lookup
	location, err := LookupIP(ip)
	if err != nil {
		writeGeoIPError(w, http.StatusServiceUnavailable)
		return
	}

	writeGeoIPJSON(w, http.StatusOK, map[string]interface{}{"ok": true, "data": location})
}

// LookupTextHandler handles GeoIP lookup requests with plain text response
func LookupTextHandler(w http.ResponseWriter, r *http.Request) {
	// Get IP from query parameter or use client IP
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		ip = getClientIP(r)
	}

	// Perform lookup
	location, err := LookupIP(ip)
	if err != nil {
		http.Error(w, "GeoIP lookup is unavailable", http.StatusServiceUnavailable)
		return
	}

	// Return text response
	w.Header().Set("Content-Type", "text/plain")
	response := formatTextResponse(location)
	if _, err := w.Write([]byte(response)); err != nil {
		return
	}
}

// BatchLookupHandler handles batch GeoIP lookups
func BatchLookupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		IPs []string `json:"ips"`
	}

	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&request); err != nil {
		writeGeoIPError(w, http.StatusBadRequest)
		return
	}

	// Limit batch size
	if len(request.IPs) > 100 {
		writeGeoIPError(w, http.StatusBadRequest)
		return
	}

	// Perform lookups
	results := make([]*Location, 0, len(request.IPs))
	for _, ip := range request.IPs {
		location, err := LookupIP(ip)
		if err != nil {
			results = append(results, &Location{IP: ip})
			continue
		}
		results = append(results, location)
	}

	writeGeoIPJSON(w, http.StatusOK, map[string]interface{}{
		"ok":   true,
		"data": map[string]interface{}{"count": len(results), "results": results},
	})
}

// getClientIP extracts the real client IP from the request
func getClientIP(r *http.Request) string {
	// Use RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

// formatTextResponse formats a Location as plain text
func formatTextResponse(loc *Location) string {
	var sb strings.Builder

	sb.WriteString("IP: " + loc.IP + "\n")

	if loc.Country != "" {
		sb.WriteString("Country: " + loc.Country)
		if loc.CountryCode != "" {
			sb.WriteString(" (" + loc.CountryCode + ")")
		}
		sb.WriteString("\n")
	}

	if loc.City != "" {
		sb.WriteString("City: " + loc.City + "\n")
	}

	if loc.Latitude != 0 || loc.Longitude != 0 {
		sb.WriteString("Coordinates: ")
		sb.WriteString(formatFloat(loc.Latitude))
		sb.WriteString(", ")
		sb.WriteString(formatFloat(loc.Longitude))
		sb.WriteString("\n")
	}

	if loc.Timezone != "" {
		sb.WriteString("Timezone: " + loc.Timezone + "\n")
	}

	if loc.ASN != 0 {
		sb.WriteString("ASN: ")
		sb.WriteString(formatUint(loc.ASN))
		if loc.ASNOrg != "" {
			sb.WriteString(" (" + loc.ASNOrg + ")")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// formatFloat formats a float64 to string
func formatFloat(f float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.6f", f), "0"), ".")
}

// formatUint formats a uint to string
func formatUint(u uint) string {
	return fmt.Sprintf("%d", u)
}

var nonPublicPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"), netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"), netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("169.254.0.0/16"), netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.0.0.0/24"), netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.88.99.0/24"), netip.MustParsePrefix("192.168.0.0/16"),
	netip.MustParsePrefix("198.18.0.0/15"), netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"), netip.MustParsePrefix("224.0.0.0/4"),
	netip.MustParsePrefix("240.0.0.0/4"), netip.MustParsePrefix("::/128"),
	netip.MustParsePrefix("::1/128"), netip.MustParsePrefix("100::/64"),
	netip.MustParsePrefix("2001:db8::/32"), netip.MustParsePrefix("2001:10::/28"),
	netip.MustParsePrefix("fc00::/7"), netip.MustParsePrefix("fe80::/10"),
	netip.MustParsePrefix("ff00::/8"),
}

func isPublicAddress(ip net.IP) bool {
	address, ok := netip.AddrFromSlice(ip)
	if !ok {
		return false
	}
	address = address.Unmap()
	if !address.IsGlobalUnicast() {
		return false
	}
	for _, prefix := range nonPublicPrefixes {
		if prefix.Contains(address) {
			return false
		}
	}
	return true
}

func writeGeoIPJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(data)
}

func writeGeoIPError(w http.ResponseWriter, status int) {
	writeGeoIPJSON(w, status, map[string]interface{}{
		"ok":      false,
		"error":   "GEOIP_LOOKUP_FAILED",
		"message": "The GeoIP lookup could not be completed.",
		"details": map[string]interface{}{},
	})
}
