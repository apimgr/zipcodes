package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apimgr/zipcodes/src/database"
	"github.com/go-chi/chi/v5"
)

func setupTestAPI(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.Initialize(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.LoadFromJSON([]byte(`[{"state":"CA","city":"Beverly Hills","county":"Los Angeles","zip_code":90210,"latitude":34.0901,"longitude":-118.4065},{"state":"CA","city":"Los Angeles","county":"Los Angeles","zip_code":90001,"latitude":33.9731,"longitude":-118.2479},{"state":"NY","city":"New York","county":"New York","zip_code":10001,"latitude":40.7506,"longitude":-73.9972}]`)); err != nil {
		t.Fatal(err)
	}
	SetDatabase(db)
	SetZipcodesJSON([]byte(`[{"zip_code":90210}]`))
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestSearchHandler(t *testing.T) {
	setupTestAPI(t)
	for _, tc := range []struct {
		query    string
		status   int
		contains string
	}{
		{"", http.StatusBadRequest, "MISSING_PARAMETER"},
		{"90210", http.StatusOK, "Beverly Hills"},
		{"99999", http.StatusNotFound, "NOT_FOUND"},
		{"Los Angeles, CA", http.StatusOK, "Los Angeles"},
		{"New York", http.StatusOK, "New York"},
		{"90", http.StatusOK, "90210"},
		{"!", http.StatusBadRequest, "INVALID_QUERY"},
	} {
		t.Run(tc.query, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/?q="+strings.ReplaceAll(tc.query, " ", "%20"), nil)
			response := httptest.NewRecorder()
			SearchHandler(response, request)
			if response.Code != tc.status || !strings.Contains(response.Body.String(), tc.contains) {
				t.Fatalf("response = %d %s", response.Code, response.Body.String())
			}
		})
	}
}

func TestLookupAndListHandlers(t *testing.T) {
	setupTestAPI(t)
	requestFor := func(key, value string) *http.Request {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add(key, value)
		return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	}
	for _, tc := range []struct {
		name       string
		handler    http.HandlerFunc
		key, value string
		status     int
		contains   string
	}{
		{"valid zip", GetByZipCodeHandler, "code", "90210", http.StatusOK, "Beverly Hills"},
		{"invalid zip", GetByZipCodeHandler, "code", "bad", http.StatusBadRequest, "INVALID_FORMAT"},
		{"missing zip", GetByZipCodeHandler, "code", "99999", http.StatusNotFound, "NOT_FOUND"},
		{"city", GetByCityHandler, "city", "Los Angeles", http.StatusOK, "Los Angeles"},
		{"missing city", GetByCityHandler, "city", "", http.StatusBadRequest, "MISSING_PARAMETER"},
		{"state", GetByStateHandler, "state", "CA", http.StatusOK, "90210"},
		{"missing state", GetByStateHandler, "state", "", http.StatusBadRequest, "MISSING_PARAMETER"},
		{"text zip", GetByZipCodeTextHandler, "code", "90210", http.StatusOK, "Coordinates:"},
		{"text zip invalid", GetByZipCodeTextHandler, "code", "bad", http.StatusBadRequest, "Invalid zipcode"},
		{"text zip absent", GetByZipCodeTextHandler, "code", "99999", http.StatusNotFound, "not found"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			tc.handler(response, requestFor(tc.key, tc.value))
			if response.Code != tc.status || !strings.Contains(response.Body.String(), tc.contains) {
				t.Fatalf("response = %d %s", response.Code, response.Body.String())
			}
		})
	}
}

func TestAutocompleteStatsAndRawHandlers(t *testing.T) {
	setupTestAPI(t)
	for _, tc := range []struct {
		handler        http.HandlerFunc
		path, contains string
	}{
		{AutoCompleteHandler, "/", "suggestions"},
		{AutoCompleteHandler, "/?q=Los&limit=1", "Los Angeles, CA"},
		{AutoCompleteHandler, "/?q=New&limit=not-a-number", "New York, NY"},
		{AutoCompleteHandler, "/?q=New&limit=99", "New York, NY"},
		{StatsHandler, "/", "total_zipcodes"},
		{StatsTextHandler, "/", "Total States: 2"},
		{CountHandler, "/", `"count": 3`},
		{CountTextHandler, "/", "3"},
		{RandomHandler, "/", "zip_code"},
		{RandomTextHandler, "/", "Zip Code:"},
		{RawJSONHandler, "/", "90210"},
	} {
		response := httptest.NewRecorder()
		tc.handler(response, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), tc.contains) {
			t.Errorf("%s response = %d %s", tc.path, response.Code, response.Body.String())
		}
	}
	response := httptest.NewRecorder()
	RawJSONHandler(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Header().Get("Content-Disposition") != `inline; filename="zipcodes.json"` {
		t.Fatalf("content disposition = %q", response.Header().Get("Content-Disposition"))
	}
	if isNumeric("123") != true || isNumeric("") != true || isNumeric("12a") {
		t.Fatal("isNumeric returned an unexpected result")
	}
	if got := formatZipcodeText(&database.Zipcode{ZipCode: 1, City: "Town", State: "ST"}); strings.Contains(got, "Coordinates:") || strings.Contains(got, "County:") {
		t.Fatalf("optional fields were rendered: %q", got)
	}
}

func TestDatabaseErrorResponses(t *testing.T) {
	db := setupTestAPI(t)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	for _, handler := range []http.HandlerFunc{SearchHandler, StatsHandler, CountHandler, RandomHandler} {
		response := httptest.NewRecorder()
		handler(response, httptest.NewRequest(http.MethodGet, "/?q=90210", nil))
		if response.Code != http.StatusInternalServerError {
			t.Errorf("status = %d", response.Code)
		}
	}
}

func TestNearbyHandler(t *testing.T) {
	setupTestAPI(t)
	for _, tc := range []struct {
		name, path string
		status     int
		contains   string
	}{
		{"missing radius", "/?lat=34&lng=-118", http.StatusBadRequest, "INVALID_RADIUS"},
		{"invalid radius", "/?lat=34&lng=-118&radius=0", http.StatusBadRequest, "INVALID_RADIUS"},
		{"invalid zip", "/?zip=9021&radius=5", http.StatusBadRequest, "INVALID_ZIPCODE"},
		{"unknown zip", "/?zip=99999&radius=5", http.StatusNotFound, "NOT_FOUND"},
		{"invalid latitude", "/?lat=nope&lng=-118&radius=5", http.StatusBadRequest, "INVALID_LATITUDE"},
		{"invalid longitude", "/?lat=34&lng=nope&radius=5", http.StatusBadRequest, "INVALID_LONGITUDE"},
		{"coordinates", "/?lat=34.09&lng=-118.40&radius=10", http.StatusOK, "Beverly Hills"},
		{"zip origin", "/?zip=90210&radius=10", http.StatusOK, "Beverly Hills"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			NearbyHandler(response, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if response.Code != tc.status || !strings.Contains(response.Body.String(), tc.contains) {
				t.Fatalf("response = %d %s", response.Code, response.Body.String())
			}
		})
	}
}
