package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/apimgr/zipcodes/src/admin"
	"github.com/apimgr/zipcodes/src/api"
	"github.com/apimgr/zipcodes/src/config"
	"github.com/apimgr/zipcodes/src/database"
	"github.com/apimgr/zipcodes/src/geoip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static
var staticFiles embed.FS

//go:embed templates
var templateFiles embed.FS

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
	router       *chi.Mux
	db           *database.DB
	config       *config.Config
	adminHandler *admin.Handler
	address      string
	port         string
	version      string
	buildDate    string
	commit       string
	zipcodesData []byte
}

// New creates a new server instance
func New(db *database.DB, cfg *config.Config, address, port, version, buildDate, commit string, zipcodesData []byte) *Server {
	s := &Server{
		router:       chi.NewRouter(),
		db:           db,
		config:       cfg,
		address:      address,
		port:         port,
		version:      version,
		buildDate:    buildDate,
		commit:       commit,
		zipcodesData: zipcodesData,
	}

	// Initialize admin handler
	s.adminHandler = admin.NewHandler(
		cfg.Server.Admin.Username,
		cfg.Server.Admin.Password,
		cfg.Server.Admin.APIToken,
		cfg.Server.Session.Timeout,
		false, // SSL enabled
		version,
	)

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

	// Real IP
	s.router.Use(middleware.RealIP)

	// Logger
	s.router.Use(middleware.Logger)

	// Compression
	s.router.Use(middleware.Compress(5))

	// Timeout
	s.router.Use(middleware.Timeout(60 * time.Second))

	// Throttle concurrent requests
	s.router.Use(middleware.Throttle(1000))

	// CORS
	s.router.Use(s.corsMiddleware)

	// Security headers
	s.router.Use(s.securityHeadersMiddleware)
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Register admin routes
	s.adminHandler.RegisterRoutes(s.router)

	// Static files
	fileServer := http.FileServer(http.FS(staticFS))
	s.router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Public routes
	s.router.Get("/", s.handleHome)
	s.router.Get("/healthz", s.handleHealth)
	s.router.Get("/health", s.handleHealth)
	s.router.Get("/status", s.handleHealth)

	// PWA support
	s.router.Get("/manifest.json", s.handleManifest)
	s.router.Get("/sw.js", s.handleServiceWorker)
	s.router.Get("/robots.txt", s.handleRobotsTxt)
	s.router.Get("/security.txt", s.handleSecurityTxt)
	s.router.Get("/.well-known/security.txt", s.handleSecurityTxt)

	// Documentation routes
	s.router.Get("/openapi", s.handleSwaggerUI)
	s.router.Get("/graphql", s.handleGraphQLPlayground)

	// API v1 routes (all public, no auth)
	s.router.Route("/api/v1", func(r chi.Router) {
		// Documentation endpoints
		r.Get("/openapi", s.handleSwaggerUI)
		r.Get("/openapi.json", s.handleOpenAPISpec)
		r.Get("/graphql", s.handleGraphQLPlayground)
		r.Post("/graphql", s.handleGraphQL)

		// Health endpoint
		r.Get("/health", s.handleHealth)

		// Raw JSON file endpoint
		r.Get("/zipcodes.json", api.RawJSONHandler)

		// Zipcode endpoints - JSON
		r.Get("/zipcode/search", api.SearchHandler)
		r.Get("/zipcode/autocomplete", api.AutoCompleteHandler)
		r.Get("/zipcode/stats", api.StatsHandler)
		r.Get("/zipcode/{code}", api.GetByZipCodeHandler)
		r.Get("/zipcode/city/{city}", api.GetByCityHandler)
		r.Get("/zipcode/state/{state}", api.GetByStateHandler)
		r.Get("/zipcode/random", api.RandomHandler)

		// Zipcode endpoints - Plain text (.txt)
		r.Get("/zipcode/{code}.txt", api.GetByZipCodeTextHandler)
		r.Get("/zipcode/stats.txt", api.StatsTextHandler)
		r.Get("/zipcode/random.txt", api.RandomTextHandler)

		// GeoIP endpoints
		r.Get("/geoip", geoip.LookupHandler)
		r.Get("/geoip.txt", geoip.LookupTextHandler)
		r.Post("/geoip/batch", geoip.BatchLookupHandler)

		// Stats & count
		r.Get("/stats", api.StatsHandler)
		r.Get("/stats.txt", api.StatsTextHandler)
		r.Get("/count", api.CountHandler)
		r.Get("/count.txt", api.CountTextHandler)
	})

	// Shorthand routes
	s.router.Get("/random", api.RandomHandler)
	s.router.Get("/random.txt", api.RandomTextHandler)
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

// handleHome serves the main page
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	data, err := templateFiles.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// handleHealth returns server health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	stats, err := s.db.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"status":        "healthy",
			"version":       s.version,
			"total_zips":    stats["total_zipcodes"],
			"total_states":  stats["total_states"],
			"geoip_enabled": geoip.GetInstance() != nil,
		},
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
			{"src": "/static/images/icon-192.png", "sizes": "192x192", "type": "image/png"},
			{"src": "/static/images/icon-512.png", "sizes": "512x512", "type": "image/png"},
		},
	}
	respondJSON(w, http.StatusOK, manifest)
}

// handleServiceWorker serves the service worker
func (s *Server) handleServiceWorker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte(`// Service Worker for Zipcodes API
const CACHE_NAME = 'zipcodes-api-v1';
const urlsToCache = [
  '/',
  '/static/css/main.css',
  '/static/js/main.js'
];

self.addEventListener('install', event => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(cache => cache.addAll(urlsToCache))
  );
});

self.addEventListener('fetch', event => {
  event.respondWith(
    caches.match(event.request)
      .then(response => response || fetch(event.request))
  );
});
`))
}

// handleRobotsTxt serves robots.txt
func (s *Server) handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	cfg := config.Get()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "User-agent: *")
	for _, path := range cfg.WebRobots.Allow {
		fmt.Fprintf(w, "Allow: %s\n", path)
	}
	for _, path := range cfg.WebRobots.Deny {
		fmt.Fprintf(w, "Disallow: %s\n", path)
	}
}

// handleSecurityTxt serves security.txt
func (s *Server) handleSecurityTxt(w http.ResponseWriter, r *http.Request) {
	cfg := config.Get()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "# Security Policy")
	if cfg.WebSecurity.Admin != "" {
		fmt.Fprintf(w, "Contact: mailto:%s\n", cfg.WebSecurity.Admin)
	} else {
		fmt.Fprintln(w, "Contact: mailto:security@example.com")
	}
	fmt.Fprintln(w, "Preferred-Languages: en")
}

// APIResponse represents a standardized API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an API error
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// respondJSON writes a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError writes an error JSON response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    status,
			Message: message,
		},
	})
}
