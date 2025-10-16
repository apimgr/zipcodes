package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apimgr/zipcodes/src/database"
	"github.com/go-chi/chi/v5"
)

var db *database.DB
var zipcodesJSON []byte

// SetDatabase sets the database instance for handlers
func SetDatabase(database *database.DB) {
	db = database
}

// SetZipcodesJSON sets the embedded JSON data for raw JSON endpoint
func SetZipcodesJSON(data []byte) {
	zipcodesJSON = data
}

// SearchHandler handles zipcode search requests
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "MISSING_PARAMETER", "message": "query parameter 'q' is required"},
		})
		return
	}

	// Try to parse as zipcode number
	if zipCode, err := strconv.Atoi(query); err == nil {
		result, err := db.SearchByZipCode(zipCode)
		if err != nil {
			respondError(w, err)
			return
		}
		if result == nil {
			respondJSON(w, http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   map[string]string{"code": "NOT_FOUND", "message": "zipcode not found"},
			})
			return
		}
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    result,
		})
		return
	}

	// Try state, city format
	parts := strings.Split(query, ",")
	if len(parts) == 2 {
		state := strings.TrimSpace(parts[1])
		city := strings.TrimSpace(parts[0])
		results, err := db.SearchByStateAndCity(state, city)
		if err != nil {
			respondError(w, err)
			return
		}
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"count":   len(results),
			"data":    results,
		})
		return
	}

	// Try as city name
	if len(query) > 2 && !isNumeric(query) {
		results, err := db.SearchByCity(query)
		if err != nil {
			respondError(w, err)
			return
		}
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"count":   len(results),
			"data":    results,
		})
		return
	}

	// Try as zipcode prefix
	if isNumeric(query) {
		results, err := db.SearchByPrefix(query)
		if err != nil {
			respondError(w, err)
			return
		}
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"count":   len(results),
			"data":    results,
		})
		return
	}

	respondJSON(w, http.StatusBadRequest, map[string]interface{}{
		"success": false,
		"error":   map[string]string{"code": "INVALID_QUERY", "message": "invalid query format"},
	})
}

// GetByZipCodeHandler handles GET /api/v1/zipcode/:code
func GetByZipCodeHandler(w http.ResponseWriter, r *http.Request) {
	codeStr := chi.URLParam(r, "code")
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "INVALID_FORMAT", "message": "invalid zipcode format"},
		})
		return
	}

	result, err := db.SearchByZipCode(code)
	if err != nil {
		respondError(w, err)
		return
	}

	if result == nil {
		respondJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "NOT_FOUND", "message": "zipcode not found"},
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// GetByZipCodeTextHandler handles GET /api/v1/zipcode/:code.txt
func GetByZipCodeTextHandler(w http.ResponseWriter, r *http.Request) {
	codeStr := chi.URLParam(r, "code")
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		http.Error(w, "Invalid zipcode format", http.StatusBadRequest)
		return
	}

	result, err := db.SearchByZipCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result == nil {
		http.Error(w, "Zipcode not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	response := formatZipcodeText(result)
	w.Write([]byte(response))
}

// GetByCityHandler handles GET /api/v1/zipcode/city/:city
func GetByCityHandler(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	if city == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "MISSING_PARAMETER", "message": "city is required"},
		})
		return
	}

	results, err := db.SearchByCity(city)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"count":   len(results),
		"data":    results,
	})
}

// GetByStateHandler handles GET /api/v1/zipcode/state/:state
func GetByStateHandler(w http.ResponseWriter, r *http.Request) {
	state := chi.URLParam(r, "state")
	if state == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "MISSING_PARAMETER", "message": "state is required"},
		})
		return
	}

	results, err := db.SearchByState(state)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"count":   len(results),
		"data":    results,
	})
}

// AutoCompleteHandler handles GET /api/v1/zipcode/autocomplete
func AutoCompleteHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success":     true,
			"suggestions": []string{},
		})
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	suggestions, err := db.AutoComplete(query, limit)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":     true,
		"suggestions": suggestions,
	})
}

// StatsHandler handles GET /api/v1/zipcode/stats
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := db.GetStats()
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// RawJSONHandler serves the raw zipcodes.json file from embedded data
func RawJSONHandler(w http.ResponseWriter, r *http.Request) {
	// Serve embedded JSON
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "inline; filename=\"zipcodes.json\"")
	w.WriteHeader(http.StatusOK)
	w.Write(zipcodesJSON)
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Wrap response with timestamp if not already present
	if m, ok := data.(map[string]interface{}); ok {
		if _, hasTimestamp := m["timestamp"]; !hasTimestamp {
			m["timestamp"] = time.Now().Format(time.RFC3339)
		}
	}

	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, err error) {
	respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
		"success":   false,
		"error":     map[string]string{"message": err.Error()},
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func formatZipcodeText(zc *database.Zipcode) string {
	var sb strings.Builder

	sb.WriteString("Zip Code: ")
	sb.WriteString(strconv.Itoa(zc.ZipCode))
	sb.WriteString("\n")

	sb.WriteString("City: ")
	sb.WriteString(zc.City)
	sb.WriteString("\n")

	sb.WriteString("State: ")
	sb.WriteString(zc.State)
	sb.WriteString("\n")

	if zc.County != "" {
		sb.WriteString("County: ")
		sb.WriteString(zc.County)
		sb.WriteString("\n")
	}

	if zc.Latitude != "" && zc.Longitude != "" {
		sb.WriteString("Coordinates: ")
		sb.WriteString(zc.Latitude)
		sb.WriteString(", ")
		sb.WriteString(zc.Longitude)
		sb.WriteString("\n")
	}

	return sb.String()
}
