package server

import (
	"encoding/json"
	"net/http"

	"github.com/apimgr/zipcodes/src/config"
)

func requestSizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.RequestURI) > 8192 {
			writeRequestLimitError(w, http.StatusRequestURITooLong, "REQUEST_URI_TOO_LONG", "Request URI is too long")
			return
		}
		if r.Body != nil {
			maxBodySize, err := config.ParseSize(config.Get().Server.Limits.MaxBodySize)
			if err != nil { maxBodySize = 10 << 20 }
			r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		}
		next.ServeHTTP(w, r)
	})
}

func writeRequestLimitError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(map[string]interface{}{"ok": false, "error": code, "message": message, "details": map[string]interface{}{}})
}
