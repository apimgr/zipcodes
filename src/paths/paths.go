package paths

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// GetDefaultDirs returns OS-specific default directories
func GetDefaultDirs(projectName string) (configDir, dataDir, logsDir string) {
	// Check if running with root/admin privileges
	isRoot := false
	if runtime.GOOS == "windows" {
		isRoot = os.Getenv("USERDOMAIN") == os.Getenv("COMPUTERNAME")
	} else {
		isRoot = os.Geteuid() == 0
	}

	if isRoot {
		// System-wide installation
		switch runtime.GOOS {
		case "windows":
			programData := os.Getenv("ProgramData")
			if programData == "" {
				programData = "C:\\ProgramData"
			}
			baseDir := filepath.Join(programData, capitalizeFirst(projectName))
			configDir = filepath.Join(baseDir, "config")
			dataDir = filepath.Join(baseDir, "data")
			logsDir = filepath.Join(baseDir, "logs")

		case "darwin":
			// macOS system-wide
			baseDir := filepath.Join("/Library/Application Support", capitalizeFirst(projectName))
			configDir = baseDir
			dataDir = filepath.Join(baseDir, "data")
			logsDir = filepath.Join("/Library/Logs", capitalizeFirst(projectName))

		default:
			// Linux/BSD system-wide
			configDir = filepath.Join("/etc", projectName)
			dataDir = filepath.Join("/var/lib", projectName)
			logsDir = filepath.Join("/var/log", projectName)
		}
	} else {
		// User-specific installation
		homeDir := ""
		currentUser, err := user.Current()
		if err == nil {
			homeDir = currentUser.HomeDir
		}
		if homeDir == "" {
			homeDir = os.Getenv("HOME")
		}
		if homeDir == "" && runtime.GOOS == "windows" {
			homeDir = os.Getenv("USERPROFILE")
		}

		switch runtime.GOOS {
		case "windows":
			// Windows user paths
			appData := os.Getenv("APPDATA")
			if appData == "" {
				appData = filepath.Join(homeDir, "AppData", "Roaming")
			}
			localAppData := os.Getenv("LOCALAPPDATA")
			if localAppData == "" {
				localAppData = filepath.Join(homeDir, "AppData", "Local")
			}

			configDir = filepath.Join(appData, capitalizeFirst(projectName))
			dataDir = filepath.Join(localAppData, capitalizeFirst(projectName))
			logsDir = filepath.Join(localAppData, capitalizeFirst(projectName), "logs")

		case "darwin":
			// macOS user paths
			baseDir := filepath.Join(homeDir, "Library", "Application Support", capitalizeFirst(projectName))
			configDir = baseDir
			dataDir = filepath.Join(baseDir, "data")
			logsDir = filepath.Join(homeDir, "Library", "Logs", capitalizeFirst(projectName))

		default:
			// Linux/BSD user paths (XDG spec)
			xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
			if xdgConfigHome == "" {
				xdgConfigHome = filepath.Join(homeDir, ".config")
			}

			xdgDataHome := os.Getenv("XDG_DATA_HOME")
			if xdgDataHome == "" {
				xdgDataHome = filepath.Join(homeDir, ".local", "share")
			}

			xdgStateHome := os.Getenv("XDG_STATE_HOME")
			if xdgStateHome == "" {
				xdgStateHome = filepath.Join(homeDir, ".local", "state")
			}

			configDir = filepath.Join(xdgConfigHome, projectName)
			dataDir = filepath.Join(xdgDataHome, projectName)
			logsDir = filepath.Join(xdgStateHome, projectName)
		}
	}

	return configDir, dataDir, logsDir
}

// GetDirs returns directories with environment variable and flag overrides
func GetDirs(projectName, configFlag, dataFlag, logsFlag string) (configDir, dataDir, logsDir string) {
	// Priority order:
	// 1. Command-line flags (highest)
	// 2. Environment variables
	// 3. OS-specific defaults (lowest)

	configDir, dataDir, logsDir = GetDefaultDirs(projectName)

	// Override with environment variables
	if envConfig := os.Getenv("CONFIG_DIR"); envConfig != "" {
		configDir = envConfig
	}
	if envData := os.Getenv("DATA_DIR"); envData != "" {
		dataDir = envData
	}
	if envLogs := os.Getenv("LOGS_DIR"); envLogs != "" {
		logsDir = envLogs
	}

	// Override with command-line flags (highest priority)
	if configFlag != "" {
		configDir = configFlag
	}
	if dataFlag != "" {
		dataDir = dataFlag
	}
	if logsFlag != "" {
		logsDir = logsFlag
	}

	return configDir, dataDir, logsDir
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	// Convert first character to uppercase if it's lowercase
	first := s[0]
	if first >= 'a' && first <= 'z' {
		first = first - 32
	}
	return string(first) + s[1:]
}
