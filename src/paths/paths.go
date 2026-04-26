package paths

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	// OrgName is the organization name used for directory structure
	OrgName = "apimgr"
	// ProjectName is the name of this project
	ProjectName = "zipcodes"
)

// GetConfigDir returns the OS-specific configuration directory
func GetConfigDir(appName string) string {
	// Check environment variable first
	if configDir := os.Getenv("CONFIG_DIR"); configDir != "" {
		return configDir
	}

	var baseDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%\OrgName\AppName
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			baseDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
		return filepath.Join(baseDir, capitalize(OrgName), capitalize(appName))

	case "darwin":
		// macOS: ~/Library/Application Support/OrgName/AppName
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, "Library", "Application Support", capitalize(OrgName), capitalize(appName))

	default:
		// Linux/Unix: Check if running as root
		if os.Geteuid() == 0 {
			// Root user: /etc/orgname/appname
			return filepath.Join("/etc", OrgName, appName)
		}
		// Regular user: ~/.config/orgname/appname
		if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
			return filepath.Join(xdgConfig, OrgName, appName)
		}
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".config", OrgName, appName)
	}
}

// GetDataDir returns the OS-specific data directory
func GetDataDir(appName string) string {
	// Check environment variable first
	if dataDir := os.Getenv("DATA_DIR"); dataDir != "" {
		return dataDir
	}

	var baseDir string

	switch runtime.GOOS {
	case "windows":
		// Windows: %LOCALAPPDATA%\OrgName\AppName
		baseDir = os.Getenv("LOCALAPPDATA")
		if baseDir == "" {
			baseDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
		return filepath.Join(baseDir, capitalize(OrgName), capitalize(appName))

	case "darwin":
		// macOS: ~/Library/Application Support/OrgName/AppName
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, "Library", "Application Support", capitalize(OrgName), capitalize(appName))

	default:
		// Linux/Unix: Check if running as root
		if os.Geteuid() == 0 {
			// Root user: /var/lib/orgname/appname
			return filepath.Join("/var/lib", OrgName, appName)
		}
		// Regular user: ~/.local/share/orgname/appname
		if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
			return filepath.Join(xdgData, OrgName, appName)
		}
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, ".local", "share", OrgName, appName)
	}
}

// GetLogsDir returns the OS-specific logs directory
func GetLogsDir(appName string) string {
	// Check environment variable first
	if logsDir := os.Getenv("LOGS_DIR"); logsDir != "" {
		return logsDir
	}

	switch runtime.GOOS {
	case "windows":
		// Windows: %LOCALAPPDATA%\OrgName\AppName\logs
		return filepath.Join(GetDataDir(appName), "logs")

	case "darwin":
		// macOS: ~/Library/Logs/OrgName/AppName
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, "Library", "Logs", capitalize(OrgName), capitalize(appName))

	default:
		// Linux/Unix: Check if running as root
		if os.Geteuid() == 0 {
			// Root user: /var/log/orgname/appname
			return filepath.Join("/var/log", OrgName, appName)
		}
		// Regular user: ~/.local/share/orgname/appname/logs
		return filepath.Join(GetDataDir(appName), "logs")
	}
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// capitalize returns the string with first letter capitalized
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
