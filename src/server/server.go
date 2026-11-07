package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/apimgr/zipcodes/src/api"
	"github.com/apimgr/zipcodes/src/config"
	"github.com/apimgr/zipcodes/src/database"
	"github.com/apimgr/zipcodes/src/geoip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static
var staticFiles embed.FS

// For use in handlers
var staticFS fs.FS

func init() {
	var err error
	staticFS, err = fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
}

// Server represents the HTTP server
type Server struct {
	router         *chi.Mux
	db             *database.DB
	config         *config.Config
	address        string
	port           string
	version        string
	buildDate      string
	commit         string
	zipcodesData   []byte
	trustedProxies []*net.IPNet
}

// New creates a new server instance
func New(db *database.DB, cfg *config.Config, address, port, version, buildDate, commit string, zipcodesData []byte) *Server {
	s := &Server{
		router:         chi.NewRouter(),
		db:             db,
		config:         cfg,
		address:        address,
		port:           port,
		version:        version,
		buildDate:      buildDate,
		commit:         commit,
		zipcodesData:   zipcodesData,
		trustedProxies: buildTrustedProxies(address, cfg.Server.TrustedProxies.Additional),
	}

	// Set embedded JSON data for API handlers
	api.SetZipcodesJSON(zipcodesData)
	api.SetDatabase(db)

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

// Router returns the chi router
func (s *Server) Router() *chi.Mux {
	return s.router
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	// Recovery from panics
	s.router.Use(middleware.Recoverer)

	// Request ID
	s.router.Use(middleware.RequestID)

	// Resolve forwarding headers only when the immediate peer is trusted.
	s.router.Use(s.trustedProxyMiddleware)

	// Enforce public per-IP abuse limits after trusted proxy resolution.
	s.router.Use(rateLimitMiddleware(s.db))

	// Logger
	s.router.Use(middleware.Logger)

	// Compression and request deadlines follow validated server configuration.
	if s.config.Server.Compression.Enabled {
		s.router.Use(middleware.Compress(s.config.Server.Compression.Level, s.config.Server.Compression.Types...))
	}
	writeTimeout, err := time.ParseDuration(s.config.Server.Limits.WriteTimeout)
	if err != nil || writeTimeout <= 0 {
		writeTimeout = 30 * time.Second
	}
	s.router.Use(middleware.Timeout(writeTimeout))

	// Bound request-line and body sizes before handlers parse untrusted input.
	s.router.Use(requestSizeMiddleware)

	// Throttle concurrent requests
	s.router.Use(middleware.Throttle(1000))

	// CORS
	s.router.Use(s.corsMiddleware)

	// Security headers
	s.router.Use(s.securityHeadersMiddleware)
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Static files
	fileServer := http.FileServer(http.FS(staticFS))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Public routes
	s.router.Get("/", s.handleSearchPage)
	s.router.Get("/search", s.handleSearchPage)
	s.router.Get("/healthz", s.handleHealth)
	s.router.Get("/server/about", s.handleAbout)
	s.router.Get("/server/help", s.handleHelp)
	s.router.Get("/server/healthz", s.handleHealthPage)
	s.router.Get("/server/privacy", s.handlePrivacy)
	s.router.Get("/server/terms", s.handleTerms)

	// PWA support
	s.router.Get("/manifest.json", s.handleManifest)
	s.router.Get("/sw.js", s.handleServiceWorker)
	s.router.Get("/robots.txt", s.handleRobotsTxt)
	s.router.Get("/security.txt", s.handleSecurityTxt)
	s.router.Get("/.well-known/security.txt", s.handleSecurityTxt)

	// Documentation routes
	s.router.Get("/server/swagger", s.handleSwaggerUI)
	s.router.Get("/graphql", s.handleGraphQLPlayground)

	// API v1 routes (all public, no auth)
	s.router.Route("/api/v1", func(r chi.Router) {
		// Documentation endpoints
		r.Get("/server/swagger", s.handleSwaggerUI)
		r.Get("/server/openapi.json", s.handleOpenAPISpec)
		r.Get("/graphql", s.handleGraphQLPlayground)
		r.Post("/graphql", s.handleGraphQL)

		// Server information pages and API equivalents
		r.Get("/server/about", s.handleAbout)
		r.Get("/server/help", s.handleHelp)
		r.Get("/server/healthz", s.handleHealthPage)
		r.Get("/server/privacy", s.handlePrivacy)
		r.Get("/server/terms", s.handleTerms)

		// Raw JSON file endpoint
		r.Get("/zipcodes.json", api.RawJSONHandler)

		// Zipcode endpoints - JSON
		r.Get("/zipcodes/search", api.SearchHandler)
		r.Get("/zipcodes/autocomplete", api.AutoCompleteHandler)
		r.Get("/zipcodes/nearby", api.NearbyHandler)
		r.Get("/zipcodes/stats", api.StatsHandler)
		r.Get("/zipcodes/{code}", api.GetByZipCodeHandler)
		r.Get("/zipcodes/city/{city}", api.GetByCityHandler)
		r.Get("/zipcodes/state/{state}", api.GetByStateHandler)
		r.Get("/zipcodes/random", api.RandomHandler)

		// Zipcode endpoints - Plain text (.txt)
		r.Get("/zipcodes/{code}.txt", api.GetByZipCodeTextHandler)
		r.Get("/zipcodes/stats.txt", api.StatsTextHandler)
		r.Get("/zipcodes/random.txt", api.RandomTextHandler)

		// GeoIP endpoints
		r.Get("/me", geoip.MeHandler)
		r.Get("/geoip", geoip.LookupHandler)
		r.Get("/geoip.txt", geoip.LookupTextHandler)
		r.Post("/geoip/batch", geoip.BatchLookupHandler)

		// Stats & count
	})

}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors := config.GetCORS()
		w.Header().Set("Access-Control-Allow-Origin", cors)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// securityHeadersMiddleware adds security headers
func (s *Server) securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

// handleHealth returns server health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	stats, err := s.db.GetStats()
	if err != nil {
		if strings.Contains(strings.ToLower(r.Header.Get("Accept")), "text/plain") {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = fmt.Fprintln(w, "ERROR: The service is unavailable.")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "unavailable"})
		return
	}

	if strings.Contains(strings.ToLower(r.Header.Get("Accept")), "text/plain") {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintln(w, "OK: Zip code lookup is available.")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":        "healthy",
		"version":       s.version,
		"total_zips":    stats["total_zipcodes"],
		"total_states":  stats["total_states"],
		"geoip_enabled": geoip.GetInstance() != nil,
	})
}

// handleManifest serves the PWA manifest
func (s *Server) handleManifest(w http.ResponseWriter, r *http.Request) {
	manifest := map[string]interface{}{
		"name":             "Zipcodes API",
		"short_name":       "Zipcodes",
		"description":      "US ZIP code lookup API with GeoIP support",
		"start_url":        "/",
		"display":          "standalone",
		"background_color": "#1a1a1a",
		"theme_color":      "#0066cc",
		"icons": []map[string]string{
			{"src": "/static/images/icon.svg", "sizes": "any", "type": "image/svg+xml", "purpose": "any maskable"},
		},
	}
	respondJSON(w, http.StatusOK, manifest)
}

// handleServiceWorker serves the service worker
func (s *Server) handleServiceWorker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	if _, err := io.WriteString(w, `const CACHE_NAME = "zipcodes-v2";
const APP_SHELL = ["/", "/static/css/main.css", "/static/js/theme.js", "/static/images/icon.svg", "/manifest.json"];

self.addEventListener("install", (event) => {
  event.waitUntil(caches.open(CACHE_NAME).then((cache) => cache.addAll(APP_SHELL)));
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  event.waitUntil(caches.keys().then((keys) => Promise.all(keys.filter((key) => key.startsWith("zipcodes-") && key !== CACHE_NAME).map((key) => caches.delete(key)))));
  self.clients.claim();
});

self.addEventListener("fetch", (event) => {
  const request = event.request;
  const url = new URL(request.url);
  if (request.method !== "GET" || url.origin !== self.location.origin || url.pathname.startsWith("/api/")) return;
  if (request.mode === "navigate") {
    event.respondWith(fetch(request).then((response) => {
      if (response.ok) caches.open(CACHE_NAME).then((cache) => cache.put("/", response.clone()));
      return response;
    }).catch(async () => (await caches.match("/")) || Response.error()));
    return;
  }
  event.respondWith(caches.match(request).then((cached) => {
    const refresh = fetch(request).then((response) => {
      if (response.ok) caches.open(CACHE_NAME).then((cache) => cache.put(request, response.clone()));
      return response;
    });
    return cached || refresh;
  }));
});
`); err != nil {
		return
	}
}

// handleRobotsTxt serves robots.txt
func (s *Server) handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	cfg := config.Get()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := fmt.Fprintln(w, "User-agent: *"); err != nil {
		return
	}
	for _, path := range cfg.WebRobots.Allow {
		if _, err := fmt.Fprintf(w, "Allow: %s\n", path); err != nil {
			return
		}
	}
	for _, path := range cfg.WebRobots.Deny {
		if _, err := fmt.Fprintf(w, "Disallow: %s\n", path); err != nil {
			return
		}
	}
}

// handleSecurityTxt serves security.txt
func (s *Server) handleSecurityTxt(w http.ResponseWriter, r *http.Request) {
	cfg := config.Get()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := fmt.Fprintln(w, "# Security Policy"); err != nil {
		return
	}
	if cfg.WebSecurity.Admin != "" {
		if _, err := fmt.Fprintf(w, "Contact: mailto:%s\n", cfg.WebSecurity.Admin); err != nil {
			return
		}
	} else {
		if _, err := fmt.Fprintln(w, "Contact: https://github.com/apimgr/zipcodes/security/advisories/new"); err != nil {
			return
		}
	}
	if _, err := fmt.Fprintf(w, "Expires: %s\n", time.Now().UTC().AddDate(1, 0, 0).Format(time.RFC3339)); err != nil {
		return
	}
	if _, err := fmt.Fprintln(w, "Preferred-Languages: en"); err != nil {
		return
	}
	if _, err := fmt.Fprintln(w, "Canonical: https://github.com/apimgr/zipcodes/.well-known/security.txt"); err != nil {
		return
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]interface{}{
		"ok":      false,
		"error":   "REQUEST_FAILED",
		"message": message,
		"details": map[string]interface{}{},
	})
}
