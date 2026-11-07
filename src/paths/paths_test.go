package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirectoryOverridesAndDefaults(t *testing.T) {
	t.Setenv("CONFIG_DIR", "/tmp/config-override")
	t.Setenv("DATA_DIR", "/tmp/data-override")
	t.Setenv("LOGS_DIR", "/tmp/logs-override")
	if GetConfigDir("zipcodes") != "/tmp/config-override" || GetDataDir("zipcodes") != "/tmp/data-override" || GetLogsDir("zipcodes") != "/tmp/logs-override" {
		t.Fatal("directory overrides were ignored")
	}
	t.Setenv("CONFIG_DIR", "")
	t.Setenv("DATA_DIR", "")
	t.Setenv("LOGS_DIR", "")
	if GetConfigDir("zipcodes") == "" || GetDataDir("zipcodes") == "" || GetLogsDir("zipcodes") == "" {
		t.Fatal("default directories were empty")
	}
	if got := capitalize(""); got != "" {
		t.Fatalf("capitalize(empty) = %q", got)
	}
	if got := capitalize("zipcodes"); got != "Zipcodes" {
		t.Fatalf("capitalize = %q", got)
	}
	path := filepath.Join(t.TempDir(), "nested", "directory")
	if err := EnsureDir(path); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(path); err != nil || !info.IsDir() {
		t.Fatalf("EnsureDir result: %v", err)
	}
}
