package mode

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"
)

const appName = "zipcodes"

// Mode represents the application mode
type Mode string

const (
	// Production mode - optimized for production deployment
	Production Mode = "production"
	// Development mode - optimized for development
	Development Mode = "development"
)

var (
	currentMode Mode = Production // Default mode
	mu          sync.RWMutex
)

// Get returns the current application mode
func Get() Mode {
	mu.RLock()
	defer mu.RUnlock()
	return currentMode
}

// Set sets the application mode
func Set(mode string) error {
	parsed, err := ParseMode(mode)
	if err != nil {
		return err
	}

	mu.Lock()
	currentMode = parsed
	mu.Unlock()

	return nil
}

// ParseMode parses a mode string and returns the corresponding Mode
// Accepts: "dev", "development", "prod", "production"
func ParseMode(s string) (Mode, error) {
	switch s {
	case "dev", "development":
		return Development, nil
	case "prod", "production":
		return Production, nil
	default:
		return "", fmt.Errorf("invalid mode: %s (valid options: dev, development, prod, production)", s)
	}
}

// IsDevelopment returns true if the current mode is development
func IsDevelopment() bool {
	return Get() == Development
}

// IsProduction returns true if the current mode is production
func IsProduction() bool {
	return Get() == Production
}

// Initialize initializes the mode from CLI flag or environment variable
// Priority: CLI flag > MODE env var > default (production)
func Initialize(cliMode string) error {
	var modeStr string

	// Priority 1: CLI flag (if provided)
	if cliMode != "" {
		modeStr = cliMode
	} else {
		// Priority 2: MODE environment variable
		envMode := os.Getenv("MODE")
		if envMode != "" {
			modeStr = envMode
		} else {
			// Priority 3: Default to production
			modeStr = "production"
		}
	}

	return Set(modeStr)
}

// GetErrorDetail returns error details based on the current mode
// In development mode: returns full error details with stack trace
// In production mode: returns generic error message
func GetErrorDetail(err error) string {
	if err == nil {
		return ""
	}

	if IsDevelopment() {
		// Return detailed error with stack trace
		stack := debug.Stack()
		return fmt.Sprintf("Error: %v\n\nStack trace:\n%s", err, string(stack))
	}

	// Return generic error message in production
	return "An internal error occurred. Please try again later."
}

// ShouldShowDebugEndpoints returns true if debug endpoints should be enabled
func ShouldShowDebugEndpoints() bool {
	return IsDevelopment()
}

// GetCacheHeaders returns appropriate cache headers based on the current mode
// Development: no-cache headers to ensure fresh content
// Production: cache headers for optimization
func GetCacheHeaders() map[string]string {
	if IsDevelopment() {
		return map[string]string{
			"Cache-Control": "no-cache, no-store, must-revalidate",
			"Pragma":        "no-cache",
			"Expires":       "0",
		}
	}

	// Production: cache for 1 year for static assets
	return map[string]string{
		"Cache-Control": "public, max-age=31536000, immutable",
	}
}

// GetLogLevel returns the appropriate log level for the current mode
func GetLogLevel() string {
	if IsDevelopment() {
		return "debug"
	}
	return "info"
}

// ShouldCacheTemplates returns true if templates should be cached
func ShouldCacheTemplates() bool {
	return IsProduction()
}

// ShouldCacheStaticFiles returns true if static files should be cached
func ShouldCacheStaticFiles() bool {
	return IsProduction()
}

// ShouldEnableAutoReload returns true if auto-reload should be enabled
func ShouldEnableAutoReload() bool {
	return IsDevelopment()
}

// ShouldEnableProfiling returns true if profiling should be enabled
func ShouldEnableProfiling() bool {
	return IsDevelopment()
}

// GetPanicRecoveryDetail returns panic recovery details based on mode
// Development: full stack trace
// Production: generic error message
func GetPanicRecoveryDetail(recovered interface{}) string {
	if IsDevelopment() {
		stack := debug.Stack()
		return fmt.Sprintf("Panic: %v\n\nStack trace:\n%s", recovered, string(stack))
	}

	return "Internal Server Error"
}

// String returns the string representation of the mode
func (m Mode) String() string {
	return string(m)
}
