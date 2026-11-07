package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apimgr/zipcodes/src/config"
	"github.com/apimgr/zipcodes/src/database"
)

func newTestServer(t *testing.T) (*Server, *database.DB) {
	t.Helper()
	db, err := database.Initialize(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	data := []byte(`[{"state":"CA","city":"Beverly Hills","county":"Los Angeles","zip_code":90210,"latitude":34.0901,"longitude":-118.4065},{"state":"CA","city":"Los Angeles","county":"Los Angeles","zip_code":90001,"latitude":33.9731,"longitude":-118.2479}]`)
	if err := db.LoadFromJSON(data); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(t.TempDir(), "server.yml")
	if _, err := config.Load(configPath); err != nil {
		t.Fatal(err)
	}
	server := New(db, config.Get(), "127.0.0.1", "8080", "test", "today", "abc123", data)
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	return server, db
}

func TestPublicAndAPIEndpoints(t *testing.T) {
	server, _ := newTestServer(t)
	for _, tc := range []struct {
		method, path, body, contains string
		status                       int
	}{
		{http.MethodGet, "/", "", "Zipcodes", http.StatusOK},
		{http.MethodGet, "/healthz", "", "healthy", http.StatusOK},
		{http.MethodGet, "/manifest.json", "", "Zipcodes API", http.StatusOK},
		{http.MethodGet, "/sw.js", "", "CACHE_NAME", http.StatusOK},
		{http.MethodGet, "/robots.txt", "", "User-agent", http.StatusOK},
		{http.MethodGet, "/security.txt", "", "Contact:", http.StatusOK},
		{http.MethodGet, "/.well-known/security.txt", "", "Preferred-Languages", http.StatusOK},
		{http.MethodGet, "/server/swagger", "", "API Documentation", http.StatusOK},
		{http.MethodGet, "/graphql", "", "GraphQL Playground", http.StatusOK},
		{http.MethodGet, "/api/v1/server/openapi.json", "", "openapi", http.StatusOK},
		{http.MethodGet, "/api/v1/zipcodes/search?q=90210", "", "Beverly Hills", http.StatusOK},
		{http.MethodGet, "/api/v1/zipcodes/nearby?lat=34.09&lng=-118.40&radius=10", "", "Beverly Hills", http.StatusOK},
		{http.MethodGet, "/api/v1/zipcodes/autocomplete?q=Bev", "", "Beverly Hills", http.StatusOK},
		{http.MethodGet, "/api/v1/zipcodes/stats", "", "total_zipcodes", http.StatusOK},
		{http.MethodGet, "/api/v1/zipcodes/90210", "", "Beverly Hills", http.StatusOK},
		{http.MethodGet, "/api/v1/zipcodes/90210.txt", "", "Zip Code:", http.StatusOK},
		{http.MethodGet, "/api/v1/zipcodes/city/Beverly%20Hills", "", "Beverly Hills", http.StatusOK},
		{http.MethodGet, "/api/v1/zipcodes/state/CA", "", "90210", http.StatusOK},
		{http.MethodGet, "/api/v1/zipcodes/random", "", "zip_code", http.StatusOK},
		{http.MethodGet, "/api/v1/zipcodes/random.txt", "", "Zip Code:", http.StatusOK},
		{http.MethodGet, "/api/v1/zipcodes/stats.txt", "", "Total Zipcodes", http.StatusOK},
		{http.MethodGet, "/api/v1/zipcodes.json", "", "90210", http.StatusOK},
		{http.MethodGet, "/api/v1/geoip?ip=invalid", "", "GEOIP_LOOKUP_FAILED", http.StatusServiceUnavailable},
		{http.MethodGet, "/api/v1/geoip.txt?ip=invalid", "", "GeoIP lookup is unavailable", http.StatusServiceUnavailable},
		{http.MethodPost, "/api/v1/graphql", `{"query":"{ stats { totalZipcodes } }"}`, "totalZipcodes", http.StatusOK},
	} {
		t.Run(tc.path, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			response := httptest.NewRecorder()
			server.Router().ServeHTTP(response, request)
			if response.Code != tc.status || !strings.Contains(response.Body.String(), tc.contains) {
				t.Fatalf("response = %d %s", response.Code, response.Body.String())
			}
			if response.Header().Get("X-Frame-Options") != "DENY" || response.Header().Get("X-Content-Type-Options") != "nosniff" {
				t.Fatal("security headers missing")
			}
			if response.Header().Get("Access-Control-Allow-Origin") == "" {
				t.Fatal("CORS header missing")
			}
		})
	}
	options := httptest.NewRecorder()
	server.Router().ServeHTTP(options, httptest.NewRequest(http.MethodOptions, "/api/v1/zipcodes/stats", nil))
	if options.Code != http.StatusOK {
		t.Fatalf("OPTIONS status = %d", options.Code)
	}
	batch := httptest.NewRecorder()
	server.Router().ServeHTTP(batch, httptest.NewRequest(http.MethodPost, "/api/v1/geoip/batch", bytes.NewBufferString(`{"ips":["invalid"]}`)))
	if batch.Code != http.StatusOK || !strings.Contains(batch.Body.String(), "results") {
		t.Fatalf("batch response = %d %s", batch.Code, batch.Body.String())
	}
}

func TestRespondHelpersAndStaticAsset(t *testing.T) {
	server, _ := newTestServer(t)
	response := httptest.NewRecorder()
	server.Router().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/static/css/main.css", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("static status = %d", response.Code)
	}
	response = httptest.NewRecorder()
	respondError(response, http.StatusBadRequest, "bad request")
	if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), "bad request") {
		t.Fatalf("error response = %d %s", response.Code, response.Body.String())
	}
	response = httptest.NewRecorder()
	respondJSON(response, http.StatusCreated, map[string]string{"ok": "yes"})
	if response.Code != http.StatusCreated || !strings.Contains(response.Body.String(), "yes") {
		t.Fatalf("JSON response = %d %s", response.Code, response.Body.String())
	}
}

func TestHealthDatabaseFailure(t *testing.T) {
	server, db := newTestServer(t)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	server.Router().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if response.Code != http.StatusServiceUnavailable || strings.Contains(response.Body.String(), "database is closed") {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
}

func TestTrustedProxyMiddleware(t *testing.T) {
	server, _ := newTestServer(t)
	var seenRemote, originalPeer string
	var trusted bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenRemote = r.RemoteAddr
		if peer, ok := r.Context().Value(originalPeerKey{}).(string); ok {
			originalPeer = peer
		}
		trusted = server.isTrustedPeer(r)
	})
	handler := server.trustedProxyMiddleware(next)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = "10.2.3.4:8080"
	request.Header.Set("X-Forwarded-For", "198.51.100.77, 10.2.3.4")
	handler.ServeHTTP(httptest.NewRecorder(), request)
	if seenRemote != "198.51.100.77" || originalPeer != "10.2.3.4:8080" {
		t.Fatalf("remote=%q peer=%q", seenRemote, originalPeer)
	}
	if !trusted {
		t.Fatal("rewritten request lost trusted original peer")
	}

	request = httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = "198.51.100.20:8080"
	request.Header.Set("X-Forwarded-For", "203.0.113.8")
	seenRemote = ""
	handler.ServeHTTP(httptest.NewRecorder(), request)
	if seenRemote != "198.51.100.20:8080" {
		t.Fatalf("untrusted peer rewrote RemoteAddr to %q", seenRemote)
	}
	if trusted {
		t.Fatal("untrusted peer was classified as trusted")
	}
}

func TestTrustedProxySet(t *testing.T) {
	server := &Server{trustedProxies: buildTrustedProxies("127.0.0.1:8080", []string{"203.0.113.10", "198.51.100.0/24", "not-a-proxy"})}
	for peer, want := range map[string]bool{
		"127.0.0.1:1": true, "127.0.0.2:2": true, "203.0.113.10:443": true,
		"198.51.100.23:80": true, "8.8.8.8:53": false, "not-an-ip": false,
	} {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.RemoteAddr = peer
		if got := server.isTrustedPeer(request); got != want {
			t.Errorf("isTrustedPeer(%q) = %t, want %t", peer, got, want)
		}
	}
	if parseListenIP("0.0.0.0:8080") != nil || parseListenIP(":8080") != nil || parseListenIP("127.0.0.1:8080") == nil {
		t.Fatal("listen address parsing mismatch")
	}
}
