package admin

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/apimgr/zipcodes/src/database"
)

// Middleware handles admin authentication
type Middleware struct {
	db *sql.DB
}

// NewMiddleware creates admin middleware
func NewMiddleware(db *sql.DB) *Middleware {
	return &Middleware{db: db}
}

// RequireBasicAuth requires Basic Auth for web UI
func (m *Middleware) RequireBasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Zipcodes Admin"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !database.VerifyAdminPassword(m.db, username, password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Zipcodes Admin"`)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireBearerToken requires Bearer token for API
func (m *Middleware) RequireBearerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if !database.VerifyAdminToken(m.db, token) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
