package geoip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// LookupHandler handles GeoIP lookup requests
func LookupHandler(w http.ResponseWriter, r *http.Request) {
	// Get IP from query parameter or use client IP
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		ip = getClientIP(r)
	}

	// Perform lookup
	location, err := LookupIP(ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return text response
	w.Header().Set("Content-Type", "text/plain")
	response := formatTextResponse(location)
	w.Write([]byte(response))
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

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Limit batch size
	if len(request.IPs) > 100 {
		http.Error(w, "Maximum 100 IPs per request", http.StatusBadRequest)
		return
	}

	// Perform lookups
	results := make([]*Location, 0, len(request.IPs))
	for _, ip := range request.IPs {
		location, err := LookupIP(ip)
		if err != nil {
			// Include error in response but continue
			results = append(results, &Location{
				IP:      ip,
				Country: "Error: " + err.Error(),
			})
			continue
		}
		results = append(results, location)
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   len(results),
		"results": results,
	})
}

// getClientIP extracts the real client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

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
