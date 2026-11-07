package ssl

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTLSAndChallengeHandling(t *testing.T) {
	if cfg, err := NewManager(Config{}).GetTLSConfig([]string{"example.test"}); err != nil || cfg != nil {
		t.Fatalf("disabled TLS = %v, %v", cfg, err)
	}
	manager := NewManager(Config{Enabled: true, CertPath: t.TempDir()})
	if _, err := manager.GetTLSConfig([]string{"example.test"}); err == nil {
		t.Fatal("accepted absent certificate")
	}
	if _, err := manager.getLetsEncryptTLSConfig([]string{"example.test"}); err != nil {
		t.Fatal(err)
	}
	if manager.GetHTTPHandler(http.NotFoundHandler()) == nil {
		t.Fatal("fallback handler is nil")
	}

	challenges := NewChallengeServer()
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/outside", nil)
	if challenges.ServeHTTP(response, request) {
		t.Fatal("handled unrelated path")
	}
	request = httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/missing", nil)
	if !challenges.ServeHTTP(response, request) || response.Code != http.StatusNotFound {
		t.Fatal("missing challenge was not rejected")
	}
	challenges.SetToken("token", "response")
	response = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/token", nil)
	if !challenges.ServeHTTP(response, request) || response.Body.String() != "response" {
		t.Fatal("challenge response mismatch")
	}
	challenges.ClearToken("token")
}

func TestParseChallenge(t *testing.T) {
	for input, want := range map[string]string{"http": "http-01", "HTTP01": "http-01", "tls": "tls-alpn-01", "tls-alpn-01": "tls-alpn-01", "dns": "dns-01", "dns01": "dns-01", "other": "http-01"} {
		if got := ParseChallenge(input); got != want {
			t.Errorf("ParseChallenge(%q) = %q, want %q", input, got, want)
		}
	}
}
