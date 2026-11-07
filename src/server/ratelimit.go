package server

import (
	"encoding/json"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apimgr/zipcodes/src/config"
	"github.com/apimgr/zipcodes/src/database"
)

func rateLimitMiddleware(db *database.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg := config.Get().Server.RateLimit
			if !cfg.Enabled || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}
			ip := requestIP(r.RemoteAddr)
			keyDigest := sha256.Sum256([]byte(ip))
			clientKey := hex.EncodeToString(keyDigest[:])
			class := "read"
			limit := cfg.Read
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				class = "write"
				limit = cfg.Write
			}
			if isHealthPath(r.URL.Path) {
				class = "health"
				limit = cfg.Health
			}
			if limit.Requests <= 0 || limit.Window <= 0 { limit = config.DefaultConfig().Server.RateLimit.Read }
			wait, err := db.ConsumeRateLimit(clientKey, class, limit.Requests, time.Duration(limit.Window)*time.Second, time.Now())
			if err == nil && wait == 0 && cfg.GlobalBurst > 0 {
				wait, err = db.ConsumeRateLimit(clientKey, "global", cfg.GlobalBurst, time.Minute, time.Now())
			}
			if err != nil {
				writeRequestLimitError(w, http.StatusServiceUnavailable, "RATE_LIMIT_UNAVAILABLE", "Request limiting is temporarily unavailable")
				return
			}
			if wait > 0 {
				_ = db.RecordRateLimitHit(clientKey, class, time.Now())
				w.Header().Set("Retry-After", strconv.Itoa(int(wait.Seconds())+1))
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusTooManyRequests)
				encoder := json.NewEncoder(w)
				encoder.SetIndent("", "  ")
				_ = encoder.Encode(map[string]interface{}{"ok": false, "error": "RATE_LIMITED", "message": "Too many requests", "details": map[string]interface{}{}})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func requestIP(remoteAddress string) string {
	if host, _, err := net.SplitHostPort(remoteAddress); err == nil {
		if ip := net.ParseIP(host); ip != nil { return ip.String() }
	}
	if ip := net.ParseIP(strings.Trim(remoteAddress, "[]")); ip != nil { return ip.String() }
	return "unknown"
}

func isHealthPath(path string) bool {
	return strings.HasSuffix(path, "/healthz") || strings.HasSuffix(path, "/health") || path == "/status"
}

