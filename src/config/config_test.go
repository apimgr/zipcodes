package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigurationLoadSaveAndMigration(t *testing.T) {
	dir := t.TempDir()
	legacy := filepath.Join(dir, "server.yaml")
	currentPath := filepath.Join(dir, "server.yml")
	if err := os.WriteFile(legacy, []byte("web-ui:\n  theme: light\nweb-security:\n  cors: https://example.test\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(currentPath)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.WebUI.Theme != "light" || cfg.WebSecurity.CORS != "https://example.test" {
		t.Fatalf("loaded config = %#v", cfg)
	}
	if _, err := os.Stat(legacy); !os.IsNotExist(err) {
		t.Fatalf("legacy config remains: %v", err)
	}
	if err := Save(); err != nil {
		t.Fatal(err)
	}
	if got := Get(); got.Server.Address != "0.0.0.0" {
		t.Fatalf("defaults not retained: %#v", got.Server)
	}
	if got := GetCORS(); got != "https://example.test" {
		t.Fatalf("CORS = %q", got)
	}
	if got := GetTheme(); got != "light" {
		t.Fatalf("theme = %q", got)
	}
}

func TestConfigurationDefaultsAndErrors(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Server.Address != "0.0.0.0" || cfg.WebUI.Theme != "dark" || cfg.WebSecurity.CORS != "*" {
		t.Fatalf("defaults = %#v", cfg)
	}
	if got := formatStringSlice(nil); got != "[]" {
		t.Fatalf("empty slice = %q", got)
	}
	if got := formatStringSlice([]string{"/", "/api"}); got != `["/", "/api"]` {
		t.Fatalf("slice = %q", got)
	}
	file := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(file, []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(filepath.Join(file, "server.yml")); err == nil {
		t.Fatal("Load succeeded below a file")
	}
	bad := filepath.Join(t.TempDir(), "bad.yml")
	if err := os.WriteFile(bad, []byte("web-ui: ["), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(bad); err == nil {
		t.Fatal("Load accepted invalid YAML")
	}
	if GetCORS() == "" {
		t.Fatal("CORS default was empty")
	}
}
