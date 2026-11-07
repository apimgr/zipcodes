package geoip

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeoIPLookupLifecycle(t *testing.T) {
	var nilGeoIP *GeoIP
	if _, err := nilGeoIP.Lookup("127.0.0.1"); err == nil {
		t.Fatal("nil GeoIP accepted lookup")
	}
	geo := &GeoIP{}
	if _, err := geo.Lookup("invalid"); err == nil {
		t.Fatal("Lookup accepted invalid IP")
	}
	for _, ip := range []string{"192.0.2.1", "2001:db8::1"} {
		location, err := geo.Lookup(ip)
		if err != nil || location.IP != ip {
			t.Fatalf("Lookup(%q) = %#v, %v", ip, location, err)
		}
	}
	if err := geo.Reload(filepath.Join(t.TempDir(), "missing.mmdb"), "", "", ""); err == nil {
		t.Fatal("Reload accepted missing database")
	}
	if err := geo.Close(); err != nil {
		t.Fatal(err)
	}
	if err := nilGeoIP.Close(); err != nil {
		t.Fatal(err)
	}
	if GetInstance() != nil {
		t.Fatal("unexpected GeoIP singleton")
	}
	if _, err := LookupIP("192.0.2.1"); err == nil {
		t.Fatal("LookupIP worked before initialization")
	}
}

func TestGeoIPHandlersAndFormatting(t *testing.T) {
	instance = nil
	request := httptest.NewRequest(http.MethodGet, "/?ip=invalid", nil)
	response := httptest.NewRecorder()
	LookupHandler(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("lookup status = %d", response.Code)
	}
	response = httptest.NewRecorder()
	LookupTextHandler(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("text lookup status = %d", response.Code)
	}
	if getClientIP(httptest.NewRequest(http.MethodGet, "/", nil)) == "" {
		t.Fatal("client IP was empty")
	}
	request = httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = "198.51.100.10:1234"
	request.Header.Set("X-Forwarded-For", " 192.0.2.1, 192.0.2.2")
	if got := getClientIP(request); got != "198.51.100.10" {
		t.Fatalf("untrusted X-Forwarded-For was honored: %q", got)
	}
	request.Header.Del("X-Forwarded-For")
	request.Header.Set("X-Real-IP", " 192.0.2.3 ")
	if got := getClientIP(request); got != "198.51.100.10" {
		t.Fatalf("untrusted X-Real-IP was honored: %q", got)
	}

	location := &Location{IP: "192.0.2.1", Country: "Example", CountryCode: "EX", City: "Town", Latitude: 1.25, Longitude: -2.5, Timezone: "UTC", ASN: 42, ASNOrg: "Example Net"}
	formatted := formatTextResponse(location)
	for _, field := range []string{"IP:", "Country: Example (EX)", "City: Town", "Coordinates: 1.25, -2.5", "Timezone: UTC", "ASN: 42 (Example Net)"} {
		if !strings.Contains(formatted, field) {
			t.Errorf("text response omitted %q: %s", field, formatted)
		}
	}
	if formatFloat(1) != "1" || formatFloat(0) != "0" || formatUint(4) != "4" {
		t.Fatal("number formatting mismatch")
	}

	response = httptest.NewRecorder()
	BatchLookupHandler(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("batch method status = %d", response.Code)
	}
	response = httptest.NewRecorder()
	BatchLookupHandler(response, httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{")))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("invalid batch status = %d", response.Code)
	}
	tooMany, _ := json.Marshal(map[string][]string{"ips": make([]string, 101)})
	response = httptest.NewRecorder()
	BatchLookupHandler(response, httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(tooMany))))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("oversized batch status = %d", response.Code)
	}
	response = httptest.NewRecorder()
	BatchLookupHandler(response, httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"ips":["192.0.2.1"]}`)))
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "results") {
		t.Fatalf("batch response = %d %s", response.Code, response.Body.String())
	}
}

func TestDatabaseDownloadHelpers(t *testing.T) {
	dir := t.TempDir()
	paths := GetDatabasePaths(dir)
	if !strings.Contains(paths.CityIPv4DB, filepath.Join("security", "geoip", "dbip-city-ipv4.mmdb")) {
		t.Fatalf("database path = %#v", paths)
	}
	if DatabasesExist(dir) {
		t.Fatal("DatabasesExist true before files were created")
	}
	if err := os.MkdirAll(filepath.Dir(paths.CityIPv4DB), 0755); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{paths.CityIPv4DB, paths.CityIPv6DB, paths.CountryDB, paths.ASNDB} {
		if err := os.WriteFile(path, []byte("db"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if DatabasesExist(dir) {
		t.Fatal("DatabasesExist accepted invalid MMDB files")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte("database"))
	}))
	defer server.Close()
	destination := filepath.Join(t.TempDir(), "db.mmdb")
	if err := downloadFile(server.URL+"/ok", destination); err != nil {
		t.Fatal(err)
	}
	if err := downloadFile(server.URL+"/missing", filepath.Join(t.TempDir(), "missing.mmdb")); err == nil {
		t.Fatal("downloadFile accepted HTTP 404")
	}
	if defaultTimeout <= 0 {
		t.Fatal("invalid timeout")
	}
}

