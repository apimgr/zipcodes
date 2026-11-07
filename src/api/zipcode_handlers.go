package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/apimgr/zipcodes/src/config"
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
	if len(query) > 80 {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]string{"code": "QUERY_TOO_LONG", "message": "query must be no longer than 80 characters"}})
		return
	}
	if query == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "MISSING_PARAMETER", "message": "query parameter 'q' is required"},
		})
		return
	}

	// Try to parse as zipcode number
	if zipCode, err := strconv.Atoi(query); err == nil && len(query) == 5 && isNumeric(query) {
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
		respondList(w, r, results)
		return
	}

	// Try as a two-letter state abbreviation.
	if len(query) == 2 && isAlpha(query) {
		results, err := db.SearchByState(query)
		if err != nil {
			respondError(w, err)
			return
		}
		respondList(w, r, results)
		return
	}

	// Try as city name
	if len(query) > 2 && !isNumeric(query) {
		results, err := db.SearchByCity(query)
		if err != nil {
			respondError(w, err)
			return
		}
		respondList(w, r, results)
		return
	}

	// Try as zipcode prefix
	if isNumeric(query) {
		results, err := db.SearchByPrefix(query)
		if err != nil {
			respondError(w, err)
			return
		}
		respondList(w, r, results)
		return
	}

	respondJSON(w, http.StatusBadRequest, map[string]interface{}{
		"success": false,
		"error":   map[string]string{"code": "INVALID_QUERY", "message": "invalid query format"},
	})
}

// GetByZipCodeHandler handles GET /api/v1/zipcodes/:code
func GetByZipCodeHandler(w http.ResponseWriter, r *http.Request) {
	codeStr := chi.URLParam(r, "code")
	code, err := strconv.Atoi(codeStr)
	if err != nil || len(codeStr) != 5 || !isNumeric(codeStr) {
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

// GetByZipCodeTextHandler handles GET /api/v1/zipcodes/:code.txt
func GetByZipCodeTextHandler(w http.ResponseWriter, r *http.Request) {
	codeStr := chi.URLParam(r, "code")
	code, err := strconv.Atoi(codeStr)
	if err != nil || len(codeStr) != 5 || !isNumeric(codeStr) {
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
	if _, err := w.Write([]byte(response)); err != nil {
		return
	}
}

// GetByCityHandler handles GET /api/v1/zipcodes/city/:city
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

	respondList(w, r, results)
}

// GetByStateHandler handles GET /api/v1/zipcodes/state/:state
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

	respondList(w, r, results)
}

// NearbyHandler returns ZIP codes within a bounded distance of a ZIP code or coordinate.
func NearbyHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	radius, err := strconv.ParseFloat(query.Get("radius"), 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "INVALID_RADIUS", "message": "radius in miles is required"},
		})
		return
	}
	maxRadius := config.Get().Server.MaxRadiusMiles
	if maxRadius <= 0 || maxRadius > 100 {
		maxRadius = 100
	}
	if radius <= 0 || radius > maxRadius {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "INVALID_RADIUS", "message": "radius must be greater than 0 and within the configured maximum"},
		})
		return
	}

	var latitude, longitude float64
	if zipCode := query.Get("zip"); zipCode != "" {
		if len(zipCode) != 5 || !isNumeric(zipCode) {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   map[string]string{"code": "INVALID_ZIPCODE", "message": "zip must be a five-digit US ZIP code"},
			})
			return
		}
		code, _ := strconv.Atoi(zipCode)
		origin, err := db.SearchByZipCode(code)
		if err != nil {
			respondError(w, err)
			return
		}
		if origin == nil {
			respondJSON(w, http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   map[string]string{"code": "NOT_FOUND", "message": "zipcode not found"},
			})
			return
		}
		latitude, err = strconv.ParseFloat(origin.GetLatitudeString(), 64)
		if err == nil {
			longitude, err = strconv.ParseFloat(origin.GetLongitudeString(), 64)
		}
		if err != nil {
			respondJSON(w, http.StatusUnprocessableEntity, map[string]interface{}{
				"success": false,
				"error":   map[string]string{"code": "MISSING_COORDINATES", "message": "zipcode has no coordinates for radius search"},
			})
			return
		}
	} else {
		latitude, err = strconv.ParseFloat(query.Get("lat"), 64)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   map[string]string{"code": "INVALID_LATITUDE", "message": "lat must be a valid coordinate"},
			})
			return
		}
		longitude, err = strconv.ParseFloat(query.Get("lng"), 64)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   map[string]string{"code": "INVALID_LONGITUDE", "message": "lng must be a valid coordinate"},
			})
			return
		}
	}

	results, err := db.SearchWithinRadius(latitude, longitude, radius)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "INVALID_SEARCH", "message": "radius search parameters are invalid"},
		})
		return
	}
	respondList(w, r, results)
}

// AutoCompleteHandler handles GET /api/v1/zipcodes/autocomplete
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

// StatsHandler handles GET /api/v1/zipcodes/stats
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
	var records interface{}
	if err := json.Unmarshal(zipcodesJSON, &records); err != nil {
		respondError(w, err)
		return
	}
	w.Header().Set("Content-Disposition", "inline; filename=\"zipcodes.json\"")
	respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": records})
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	if result, ok := data.(map[string]interface{}); ok {
		success, hasSuccess := result["success"].(bool)
		if hasSuccess {
			delete(result, "success")
			result["ok"] = success
			if !success {
				errorValue, _ := result["error"].(map[string]string)
				code := errorValue["code"]
				message := errorValue["message"]
				if code == "" {
					code = "REQUEST_FAILED"
				}
				if message == "" {
					message = "The request could not be completed."
				}
				delete(result, "error")
				result["error"] = code
				result["message"] = message
				if _, ok := result["details"]; !ok {
					result["details"] = map[string]interface{}{}
				}
			}
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(data)
}

func respondError(w http.ResponseWriter, err error) {
	respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
		"success": false,
		"error":   map[string]string{"code": "INTERNAL_ERROR", "message": "The request could not be completed."},
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

func isAlpha(value string) bool {
	if value == "" {
		return false
	}
	for _, character := range value {
		if (character < 'A' || character > 'Z') && (character < 'a' || character > 'z') {
			return false
		}
	}
	return true
}

func formatZipcodeText(zc *database.Zipcode) string {
	var sb strings.Builder

	sb.WriteString("Zip Code: ")
	sb.WriteString(zc.Zip)
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

	lat := zc.GetLatitudeString()
	lon := zc.GetLongitudeString()
	if lat != "" && lon != "" {
		sb.WriteString("Coordinates: ")
		sb.WriteString(lat)
		sb.WriteString(", ")
		sb.WriteString(lon)
		sb.WriteString("\n")
	}

	return sb.String()
}

// RandomHandler returns a random zipcode
func RandomHandler(w http.ResponseWriter, r *http.Request) {
	result, err := db.GetRandom()
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// RandomTextHandler returns a random zipcode as plain text
func RandomTextHandler(w http.ResponseWriter, r *http.Request) {
	result, err := db.GetRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := io.WriteString(w, formatZipcodeText(result)); err != nil {
		return
	}
}

// StatsTextHandler returns stats as plain text
func StatsTextHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := db.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	for _, line := range []string{
		fmt.Sprintf("Total Zipcodes: %v\n", stats["total_zipcodes"]),
		fmt.Sprintf("Total States: %v\n", stats["total_states"]),
		fmt.Sprintf("Total Cities: %v\n", stats["total_cities"]),
	} {
		if _, err := io.WriteString(w, line); err != nil {
			return
		}
	}
}

// CountHandler returns the total count of zipcodes
func CountHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := db.GetStats()
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"count":   stats["total_zipcodes"],
	})
}

// CountTextHandler returns the total count as plain text
func CountTextHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := db.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := fmt.Fprintf(w, "%v\n", stats["total_zipcodes"]); err != nil {
		return
	}
}

func respondList[T any](w http.ResponseWriter, r *http.Request, records []T) {
	page := 1
	limit := 250
	query := r.URL.Query()
	if value := query.Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]string{"code": "INVALID_PAGE", "message": "page must be a positive integer"}})
			return
		}
		page = parsed
	}
	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			respondJSON(w, http.StatusBadRequest, map[string]interface{}{"success": false, "error": map[string]string{"code": "INVALID_LIMIT", "message": "limit must be a positive integer"}})
			return
		}
		limit = parsed
		if limit > 250 {
			limit = 250
		}
	}
	total := len(records)
	pages := (total + limit - 1) / limit
	start := total
	if page <= total/limit+1 {
		start = (page - 1) * limit
		if start > total {
			start = total
		}
	}
	end := start + limit
	if end < start || end > total {
		end = total
	}
	if records == nil {
		records = []T{}
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"data":       records[start:end],
		"pagination": map[string]int{"page": page, "limit": limit, "total": total, "pages": pages},
	})
}
