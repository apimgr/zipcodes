package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunLookupAndSaveServer(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configRoot)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/zipcodes/90210" {
			t.Errorf("request path = %q", r.URL.Path)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("Accept = %q", got)
		}
		_, _ = io.WriteString(w, `{"ok":true,"data":{"zip":"90210","city":"Beverly Hills","state":"CA","county":"Los Angeles","latitude":34.09,"longitude":-118.41,"timezone":"America/Los_Angeles","area_codes":["310"]}}`)
	}))
	defer server.Close()

	var stdout, stderr strings.Builder
	if err := run([]string{"--server", server.URL, "90210"}, &stdout, &stderr); err != nil {
		t.Fatalf("run() error = %v; stderr=%s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Beverly Hills, CA") || !strings.Contains(stdout.String(), "310") {
		t.Fatalf("lookup output = %q", stdout.String())
	}
	configPath := filepath.Join(configRoot, internalOrg, projectName, "cli.yml")
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("saved config: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("config permissions = %o, want 600", info.Mode().Perm())
	}
}

func TestRunRejectsArgumentsAndServerURLs(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{name: "missing zip", args: []string{"--server", "https://example.test"}},
		{name: "malformed zip", args: []string{"--server", "https://example.test", "12x45"}},
		{name: "credential URL", args: []string{"--server", "https://user" + ":pass@example.test", "90210"}},
		{name: "relative URL", args: []string{"--server", "/api", "90210"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("XDG_CONFIG_HOME", t.TempDir())
			var stdout, stderr strings.Builder
			if err := run(tc.args, &stdout, &stderr); err == nil {
				t.Fatal("run() error = nil")
			}
		})
	}
}

func TestRunHandlesServerErrorsAndInvalidResponses(t *testing.T) {
	for _, tc := range []struct {
		name string
		code int
		body string
		want string
	}{
		{name: "API error", code: http.StatusNotFound, body: `{"error":{"message":"ZIP code not found"}}`, want: "ZIP code not found"},
		{name: "unsuccessful envelope", code: http.StatusOK, body: `{"ok":false}`, want: "unsuccessful"},
		{name: "invalid JSON", code: http.StatusOK, body: `{`, want: "decode server response"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("XDG_CONFIG_HOME", t.TempDir())
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.code)
				_, _ = io.WriteString(w, tc.body)
			}))
			defer server.Close()
			var stdout, stderr strings.Builder
			err := run([]string{"--server", server.URL, "90210"}, &stdout, &stderr)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("run() error = %v, want substring %q", err, tc.want)
			}
		})
	}
}

func TestRunHelpAndVersion(t *testing.T) {
	var stdout, stderr strings.Builder
	if err := run([]string{"--help"}, &stdout, &stderr); err != nil {
		t.Fatalf("help error = %v", err)
	}
	if !strings.Contains(stderr.String(), "five-digit US ZIP code") {
		t.Errorf("help = %q", stderr.String())
	}
	stdout.Reset()
	if err := run([]string{"--version"}, &stdout, &stderr); err != nil {
		t.Fatalf("version error = %v", err)
	}
	if !strings.Contains(stdout.String(), "dev") {
		t.Errorf("version = %q", stdout.String())
	}
}
