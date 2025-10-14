package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/apimgr/zipcodes/src/admin"
	"github.com/apimgr/zipcodes/src/api"
	"github.com/apimgr/zipcodes/src/database"
	"github.com/apimgr/zipcodes/src/geoip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static
var staticFiles embed.FS

//go:embed templates
var templateFiles embed.FS

// Server represents the HTTP server
type Server struct {
	router *chi.Mux
	db     *database.AppDB
	port   string
}

// New creates a new server instance
func New(db *database.AppDB, port string) *Server {
	if port == "" {
		port = "8080"
	}

	s := &Server{
		router: chi.NewRouter(),
		db:     db,
		port:   port,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Compress(5))
	s.router.Use(middleware.Timeout(60 * time.Second))

	// CORS headers
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Security headers
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			next.ServeHTTP(w, r)
		})
	})
}

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	// Set database for API handlers (use the underlying DB)
	api.SetDatabase(s.db.DB)

	// Initialize admin handlers and middleware
	adminHandler := admin.NewHandler(s.db.GetConn(), templateFiles)
	adminMw := admin.NewMiddleware(s.db.GetConn())

	// Static files
	staticFS, _ := fs.Sub(staticFiles, "static")
	s.router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Health check
	s.router.Get("/healthz", s.healthCheckHandler)

	// Homepage
	s.router.Get("/", s.indexHandler)

	// Admin routes (Basic Auth for web UI)
	s.router.Route("/admin", func(r chi.Router) {
		r.Use(adminMw.RequireBasicAuth)
		r.Get("/", adminHandler.DashboardHandler)
		r.Get("/settings", adminHandler.SettingsHandler)
		r.Post("/settings", adminHandler.SettingsHandler)
		r.Get("/database", adminHandler.DatabaseHandler)
		r.Post("/database/test", adminHandler.DatabaseTestHandler)
		r.Get("/logs", adminHandler.LogsHandler)
		r.Get("/audit", adminHandler.AuditHandler)
	})

	// API routes (public)
	s.router.Route("/api/v1", func(r chi.Router) {
		// Raw JSON file endpoint
		r.Get("/zipcodes.json", api.RawJSONHandler)

		// Zipcode endpoints
		r.Get("/zipcode/search", api.SearchHandler)
		r.Get("/zipcode/autocomplete", api.AutoCompleteHandler)
		r.Get("/zipcode/stats", api.StatsHandler)
		r.Get("/zipcode/{code}", api.GetByZipCodeHandler)
		r.Get("/zipcode/{code}.txt", api.GetByZipCodeTextHandler)
		r.Get("/zipcode/city/{city}", api.GetByCityHandler)
		r.Get("/zipcode/state/{state}", api.GetByStateHandler)

		// GeoIP endpoints
		r.Get("/geoip", geoip.LookupHandler)
		r.Get("/geoip.txt", geoip.LookupTextHandler)
		r.Post("/geoip/batch", geoip.BatchLookupHandler)

		// Admin API routes (Bearer token)
		r.Route("/admin", func(r chi.Router) {
			r.Use(adminMw.RequireBearerToken)
			r.Get("/", adminHandler.AdminInfoHandler)
			r.Get("/settings", adminHandler.SettingsHandler)
			r.Put("/settings", adminHandler.SettingsHandler)
			r.Post("/reload", adminHandler.ReloadHandler)
			r.Get("/stats", adminHandler.AdminStatsHandler)
		})
	})

	// API health endpoint (public)
	s.router.Get("/api/v1/health", s.healthCheckHandler)
}

// indexHandler serves the main page
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	data, err := templateFiles.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// healthCheckHandler provides health status
func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	_, err := s.db.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Simple JSON response
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s","database":{"status":"connected","type":"sqlite"},"features":{"zipcode_lookup":true,"geoip_lookup":%t,"api_enabled":true}}`,
		time.Now().Format(time.RFC3339),
		geoip.GetInstance() != nil,
	)
}

// Start starts the HTTP server
func (s *Server) Start(displayAddr, bindAddr string) error {
	addr := fmt.Sprintf("%s:%s", bindAddr, s.port)

	log.Printf("Listening on %s\n", addr)
	log.Printf("Access at http://%s:%s\n", displayAddr, s.port)

	return http.ListenAndServe(addr, s.router)
}
