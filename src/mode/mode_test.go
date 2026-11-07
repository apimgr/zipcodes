package mode

import (
	"errors"
	"strings"
	"testing"
)

func TestModesAndPresentation(t *testing.T) {
	if err := Initialize("dev"); err != nil {
		t.Fatal(err)
	}
	if !IsDevelopment() || IsProduction() || !ShouldShowDebugEndpoints() || ShouldCacheTemplates() || ShouldCacheStaticFiles() || !ShouldEnableAutoReload() || !ShouldEnableProfiling() {
		t.Fatal("development mode flags are inconsistent")
	}
	if got := GetCacheHeaders(); got["Cache-Control"] != "no-cache, no-store, must-revalidate" {
		t.Fatalf("development cache headers = %#v", got)
	}
	if detail := GetErrorDetail(errors.New("detail")); !strings.Contains(detail, "Stack trace:") {
		t.Fatalf("development detail = %q", detail)
	}
	if detail := GetPanicRecoveryDetail("panic"); !strings.Contains(detail, "Panic: panic") {
		t.Fatalf("panic detail = %q", detail)
	}
	if (Development).String() != "development" {
		t.Fatal("mode string mismatch")
	}
	if err := Initialize("prod"); err != nil {
		t.Fatal(err)
	}
	if !IsProduction() || GetLogLevel() != "info" || GetErrorDetail(errors.New("secret")) != "An internal error occurred. Please try again later." || GetPanicRecoveryDetail("secret") != "Internal Server Error" {
		t.Fatal("production mode behavior mismatch")
	}
	if got := GetCacheHeaders(); got["Cache-Control"] != "public, max-age=31536000, immutable" {
		t.Fatalf("production cache headers = %#v", got)
	}
	if _, err := ParseMode("invalid"); err == nil {
		t.Fatal("ParseMode accepted an invalid value")
	}
	if err := Initialize("invalid"); err == nil {
		t.Fatal("Initialize accepted an invalid mode")
	}
}

func TestModeEnvironmentInitialization(t *testing.T) {
	t.Setenv("MODE", "development")
	if err := Initialize(""); err != nil {
		t.Fatal(err)
	}
	if !IsDevelopment() {
		t.Fatal("environment mode was ignored")
	}
	t.Setenv("MODE", "production")
	if err := Initialize("dev"); err != nil {
		t.Fatal(err)
	}
	if !IsDevelopment() {
		t.Fatal("CLI mode did not take priority")
	}
}
