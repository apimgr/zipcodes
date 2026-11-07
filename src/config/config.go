package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"time"
	"sync"

	"gopkg.in/yaml.v3"
)

// OrgName is the organization name for directory paths
const OrgName = "apimgr"

// Config represents the complete server configuration
type Config struct {
	Server      ServerConfig      `yaml:"server"`
	WebUI       WebUIConfig       `yaml:"web-ui"`
	WebRobots   WebRobotsConfig   `yaml:"web-robots"`
	WebSecurity WebSecurityConfig `yaml:"web-security"`
}

// ServerConfig contains server-related settings
type ServerConfig struct {
	Port           string               `yaml:"port"`
	FQDN           string               `yaml:"fqdn"`
	Address        string               `yaml:"address"`
	TrustedProxies TrustedProxiesConfig `yaml:"trusted_proxies"`
	Mode           string               `yaml:"mode"`
	UpdateBranch   string               `yaml:"update_branch"`
	MaxRadiusMiles float64              `yaml:"max_radius_miles"`
	BaseURL        string               `yaml:"baseurl"`
	Limits         RequestLimitsConfig  `yaml:"limits"`
	Compression    CompressionConfig    `yaml:"compression"`
	GeoIP          GeoIPConfig          `yaml:"geoip"`
	RateLimit      RateLimitConfig      `yaml:"rate_limit"`
	Metrics        MetricsConfig        `yaml:"metrics"`
	Logging        LoggingConfig        `yaml:"logging"`
}

// RequestLimitsConfig controls HTTP server timeouts and request-body limits.
type RequestLimitsConfig struct {
	MaxBodySize string `yaml:"max_body_size"`
	ReadTimeout string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	IdleTimeout string `yaml:"idle_timeout"`
}

// CompressionConfig controls response compression.
type CompressionConfig struct {
	Enabled bool `yaml:"enabled"`
	Level int `yaml:"level"`
	Types []string `yaml:"types"`
}

// GeoIPConfig controls optional request location lookups.
type GeoIPConfig struct {
	Enabled bool `yaml:"enabled"`
	Dir string `yaml:"dir"`
	UpdateInterval string `yaml:"update_interval"`
	ASN bool `yaml:"asn"`
	Country bool `yaml:"country"`
	City bool `yaml:"city"`
}

// TrustedProxiesConfig contains operator supplied trusted proxy addresses.
type TrustedProxiesConfig struct {
	Additional []string `yaml:"additional"`
}

// RateLimitConfig contains per-IP request limits.
type RateLimitConfig struct {
	Enabled     bool `yaml:"enabled"`
	Read        RateLimitRule `yaml:"read"`
	Write       RateLimitRule `yaml:"write"`
	Health      RateLimitRule `yaml:"health"`
	GlobalBurst int `yaml:"global_burst"`
}

// RateLimitRule defines a count and sliding-window duration in seconds.
type RateLimitRule struct {
	Requests int `yaml:"requests"`
	Window   int `yaml:"window"`
}

// MetricsConfig contains metrics settings
type MetricsConfig struct {
	Enabled       bool   `yaml:"enabled"`
	Endpoint      string `yaml:"endpoint"`
	IncludeSystem bool   `yaml:"include_system"`
	IncludeApp    bool   `yaml:"include_app"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	AccessFormat string `yaml:"access_format"`
	Level        string `yaml:"level"`
}

// WebUIConfig contains web UI settings
type WebUIConfig struct {
	Theme         string              `yaml:"theme"`
	Logo          string              `yaml:"logo"`
	Favicon       string              `yaml:"favicon"`
	Notifications NotificationsConfig `yaml:"notifications"`
}

// NotificationsConfig contains notification settings
type NotificationsConfig struct {
	Enabled       bool     `yaml:"enabled"`
	Announcements []string `yaml:"announcements"`
}

// WebRobotsConfig contains robots.txt settings
type WebRobotsConfig struct {
	Allow []string `yaml:"allow"`
	Deny  []string `yaml:"deny"`
}

// WebSecurityConfig contains security settings
type WebSecurityConfig struct {
	Admin string `yaml:"admin"`
	CORS  string `yaml:"cors"`
}

var (
	current    *Config
	mu         sync.RWMutex
	configPath string
)

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:           "",
			FQDN:           "",
			Address:        "0.0.0.0",
			Mode:           "production",
			UpdateBranch:   "stable",
			MaxRadiusMiles: 100,
			BaseURL: "/",
			Limits: RequestLimitsConfig{MaxBodySize: "10MB", ReadTimeout: "30s", WriteTimeout: "30s", IdleTimeout: "120s"},
			Compression: CompressionConfig{Enabled: true, Level: 5, Types: []string{"text/html", "text/css", "text/javascript", "application/json", "application/xml"}},
			GeoIP: GeoIPConfig{Enabled: true, Dir: "security/geoip", UpdateInterval: "weekly", ASN: true, Country: true, City: true},
			RateLimit: RateLimitConfig{
				Enabled: true,
				Read: RateLimitRule{Requests: 120, Window: 60},
				Write: RateLimitRule{Requests: 10, Window: 60},
				Health: RateLimitRule{Requests: 120, Window: 60},
				GlobalBurst: 240,
			},
			TrustedProxies: TrustedProxiesConfig{Additional: []string{}},
			Metrics: MetricsConfig{
				Enabled:       false,
				Endpoint:      "/metrics",
				IncludeSystem: true,
				IncludeApp:    true,
			},
			Logging: LoggingConfig{
				AccessFormat: "apache",
				Level:        "info",
			},
		},
		WebUI: WebUIConfig{
			Theme:   "dark",
			Logo:    "",
			Favicon: "",
			Notifications: NotificationsConfig{
				Enabled:       true,
				Announcements: []string{},
			},
		},
		WebRobots: WebRobotsConfig{
			Allow: []string{"/", "/api"},
			Deny:  []string{"/debug"},
		},
		WebSecurity: WebSecurityConfig{
			Admin: "",
			CORS:  "*",
		},
	}
}

// migrateYamlToYml migrates old .yaml files to .yml
func migrateYamlToYml(path string) {
	if !strings.HasSuffix(path, ".yml") {
		return
	}

	oldPath := strings.TrimSuffix(path, ".yml") + ".yaml"

	// Check if old .yaml file exists and new .yml doesn't
	if _, err := os.Stat(oldPath); err == nil {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Rename .yaml to .yml
			if err := os.Rename(oldPath, path); err == nil {
				fmt.Printf("Migrated configuration: %s -> %s\n", oldPath, path)
			}
		}
	}
}

// Load loads configuration from a YAML file
func Load(path string) (*Config, error) {
	mu.Lock()
	defer mu.Unlock()

	// Migrate old .yaml to .yml if needed
	migrateYamlToYml(path)

	configPath = path

	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := saveConfig(cfg, path); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		current = cfg
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	validateConfig(cfg)
	current = cfg
	return cfg, nil
}

// Get returns the current configuration
func Get() *Config {
	mu.RLock()
	defer mu.RUnlock()
	if current == nil {
		return DefaultConfig()
	}
	return current
}

// Save saves the current configuration to file
func Save() error {
	mu.Lock()
	defer mu.Unlock()
	if current == nil || configPath == "" {
		return fmt.Errorf("no configuration loaded")
	}
	return saveConfig(current, configPath)
}

// saveConfig writes configuration to a YAML file
func saveConfig(cfg *Config, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	content := generateConfigYAML(cfg)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// generateConfigYAML generates YAML content with comments
func generateConfigYAML(cfg *Config) string {
	return fmt.Sprintf(`# Zipcodes API Server Configuration
# Documentation: https://github.com/apimgr/zipcodes

server:
  port: "%s"
  fqdn: "%s"
  address: "%s"
  baseurl: "%s"
  max_radius_miles: %.0f
  limits:
    max_body_size: "%s"
    read_timeout: "%s"
    write_timeout: "%s"
    idle_timeout: "%s"
  compression:
    enabled: %t
    level: %d
    types: %s
  geoip:
    enabled: %t
    dir: "%s"
    update_interval: "%s"
    asn: %t
    country: %t
    city: %t

  trusted_proxies:
    additional: %s

  rate_limit:
    enabled: %t
    read:
      requests: %d
      window: %d
    write:
      requests: %d
      window: %d
    health:
      requests: %d
      window: %d
    global_burst: %d

  metrics:
    enabled: %t
    endpoint: "%s"
    include_system: %t
    include_app: %t

  logging:
    access_format: "%s"
    level: "%s"

web-ui:
  theme: "%s"
  logo: "%s"
  favicon: "%s"
  notifications:
    enabled: %t
    announcements: %s

web-robots:
  allow: %s
  deny: %s

web-security:
  admin: "%s"
  cors: "%s"
`,
		cfg.Server.Port,
		cfg.Server.FQDN,
		cfg.Server.Address,
		cfg.Server.BaseURL,
		cfg.Server.MaxRadiusMiles,
		cfg.Server.Limits.MaxBodySize,
		cfg.Server.Limits.ReadTimeout,
		cfg.Server.Limits.WriteTimeout,
		cfg.Server.Limits.IdleTimeout,
		cfg.Server.Compression.Enabled,
		cfg.Server.Compression.Level,
		formatStringSlice(cfg.Server.Compression.Types),
		cfg.Server.GeoIP.Enabled,
		cfg.Server.GeoIP.Dir,
		cfg.Server.GeoIP.UpdateInterval,
		cfg.Server.GeoIP.ASN,
		cfg.Server.GeoIP.Country,
		cfg.Server.GeoIP.City,
		formatStringSlice(cfg.Server.TrustedProxies.Additional),
		cfg.Server.RateLimit.Enabled,
		cfg.Server.RateLimit.Read.Requests,
		cfg.Server.RateLimit.Read.Window,
		cfg.Server.RateLimit.Write.Requests,
		cfg.Server.RateLimit.Write.Window,
		cfg.Server.RateLimit.Health.Requests,
		cfg.Server.RateLimit.Health.Window,
		cfg.Server.RateLimit.GlobalBurst,
		cfg.Server.Metrics.Enabled,
		cfg.Server.Metrics.Endpoint,
		cfg.Server.Metrics.IncludeSystem,
		cfg.Server.Metrics.IncludeApp,
		cfg.Server.Logging.AccessFormat,
		cfg.Server.Logging.Level,
		cfg.WebUI.Theme,
		cfg.WebUI.Logo,
		cfg.WebUI.Favicon,
		cfg.WebUI.Notifications.Enabled,
		formatStringSlice(cfg.WebUI.Notifications.Announcements),
		formatStringSlice(cfg.WebRobots.Allow),
		formatStringSlice(cfg.WebRobots.Deny),
		cfg.WebSecurity.Admin,
		cfg.WebSecurity.CORS,
	)
}

func formatStringSlice(s []string) string {
	if len(s) == 0 {
		return "[]"
	}
	result := "["
	for i, v := range s {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("\"%s\"", v)
	}
	result += "]"
	return result
}

// GetTheme returns the current theme
func GetTheme() string {
	cfg := Get()
	return cfg.WebUI.Theme
}

// GetCORS returns the CORS setting
func GetCORS() string {
	cfg := Get()
	if cfg.WebSecurity.CORS == "" {
		return "*"
	}
	return cfg.WebSecurity.CORS
}


func validateConfig(cfg *Config) {
	defaults := DefaultConfig()
	if cfg.Server.Port != "" {
		port, err := strconv.Atoi(cfg.Server.Port)
		if err != nil || port < 1 || port > 65535 {
			log.Printf("Warning: invalid server.port; using a random port")
			cfg.Server.Port = ""
		}
	}
	if cfg.Server.MaxRadiusMiles <= 0 || cfg.Server.MaxRadiusMiles > 100 {
		log.Printf("Warning: invalid server.max_radius_miles; using %.0f", defaults.Server.MaxRadiusMiles)
		cfg.Server.MaxRadiusMiles = defaults.Server.MaxRadiusMiles
	}
	if cfg.Server.BaseURL == "" { cfg.Server.BaseURL = defaults.Server.BaseURL }
	if !strings.HasPrefix(cfg.Server.BaseURL, "/") || strings.Contains(cfg.Server.BaseURL, "..") || strings.ContainsAny(cfg.Server.BaseURL, "?#") {
		log.Printf("Warning: invalid server.baseurl; using %s", defaults.Server.BaseURL)
		cfg.Server.BaseURL = defaults.Server.BaseURL
	}
	if cfg.Server.GeoIP.Dir == "" { cfg.Server.GeoIP.Dir = defaults.Server.GeoIP.Dir }
	if cfg.Server.GeoIP.UpdateInterval == "" { cfg.Server.GeoIP.UpdateInterval = defaults.Server.GeoIP.UpdateInterval }
	if cfg.Server.GeoIP.UpdateInterval != "minutely" && cfg.Server.GeoIP.UpdateInterval != "hourly" && cfg.Server.GeoIP.UpdateInterval != "daily" && cfg.Server.GeoIP.UpdateInterval != "weekly" && cfg.Server.GeoIP.UpdateInterval != "monthly" {
		duration, err := time.ParseDuration(cfg.Server.GeoIP.UpdateInterval)
		if err != nil || duration < time.Hour || duration > 365*24*time.Hour {
			log.Printf("Warning: invalid server.geoip.update_interval; using %s", defaults.Server.GeoIP.UpdateInterval)
			cfg.Server.GeoIP.UpdateInterval = defaults.Server.GeoIP.UpdateInterval
		}
	}
	if !cfg.Server.GeoIP.City && !cfg.Server.GeoIP.Country && !cfg.Server.GeoIP.ASN {
		log.Printf("Warning: GeoIP is enabled but no database categories are selected; disabling GeoIP")
		cfg.Server.GeoIP.Enabled = false
	}
	if cfg.Server.Limits.MaxBodySize == "" { cfg.Server.Limits.MaxBodySize = defaults.Server.Limits.MaxBodySize }
	if cfg.Server.Limits.ReadTimeout == "" { cfg.Server.Limits.ReadTimeout = defaults.Server.Limits.ReadTimeout }
	if cfg.Server.Limits.WriteTimeout == "" { cfg.Server.Limits.WriteTimeout = defaults.Server.Limits.WriteTimeout }
	if cfg.Server.Limits.IdleTimeout == "" { cfg.Server.Limits.IdleTimeout = defaults.Server.Limits.IdleTimeout }
	if _, err := ParseSize(cfg.Server.Limits.MaxBodySize); err != nil { log.Printf("Warning: invalid server.limits.max_body_size; using default"); cfg.Server.Limits.MaxBodySize = defaults.Server.Limits.MaxBodySize }
	for value, fallback := range map[*string]string{&cfg.Server.Limits.ReadTimeout: defaults.Server.Limits.ReadTimeout, &cfg.Server.Limits.WriteTimeout: defaults.Server.Limits.WriteTimeout, &cfg.Server.Limits.IdleTimeout: defaults.Server.Limits.IdleTimeout} {
		duration, err := time.ParseDuration(*value)
		if err != nil || duration < time.Second || duration > 10*time.Minute { log.Printf("Warning: invalid server timeout; using %s", fallback); *value = fallback }
	}
	if cfg.Server.Compression.Level < 1 || cfg.Server.Compression.Level > 9 {
		log.Printf("Warning: invalid server.compression.level; using %d", defaults.Server.Compression.Level)
		cfg.Server.Compression.Level = defaults.Server.Compression.Level
	}
	if len(cfg.Server.Compression.Types) == 0 { cfg.Server.Compression.Types = defaults.Server.Compression.Types }
	if cfg.Server.RateLimit.Read.Requests <= 0 || cfg.Server.RateLimit.Read.Window <= 0 { cfg.Server.RateLimit.Read = defaults.Server.RateLimit.Read }
	if cfg.Server.RateLimit.Write.Requests <= 0 || cfg.Server.RateLimit.Write.Window <= 0 { cfg.Server.RateLimit.Write = defaults.Server.RateLimit.Write }
	if cfg.Server.RateLimit.Health.Requests <= 0 || cfg.Server.RateLimit.Health.Window <= 0 { cfg.Server.RateLimit.Health = defaults.Server.RateLimit.Health }
	if cfg.Server.RateLimit.GlobalBurst <= 0 { cfg.Server.RateLimit.GlobalBurst = defaults.Server.RateLimit.GlobalBurst }
	for _, rule := range []*RateLimitRule{&cfg.Server.RateLimit.Read, &cfg.Server.RateLimit.Write, &cfg.Server.RateLimit.Health} {
		if rule.Requests > 1000000 { rule.Requests = 1000000 }
		if rule.Window > 86400 { rule.Window = 86400 }
	}
	if cfg.Server.RateLimit.GlobalBurst > 1000000 { cfg.Server.RateLimit.GlobalBurst = 1000000 }
	if cfg.WebSecurity.CORS != "" && cfg.WebSecurity.CORS != "*" {
		origin, err := url.Parse(cfg.WebSecurity.CORS)
		if err != nil || (origin.Scheme != "http" && origin.Scheme != "https") || origin.Host == "" || origin.Path != "" || origin.RawQuery != "" || origin.Fragment != "" {
			log.Printf("Warning: invalid web-security.cors; allowing all origins for this public API")
			cfg.WebSecurity.CORS = "*"
		}
	}
	if cfg.WebUI.Theme != "dark" && cfg.WebUI.Theme != "light" && cfg.WebUI.Theme != "auto" {
		log.Printf("Warning: invalid web-ui.theme; using %s", defaults.WebUI.Theme)
		cfg.WebUI.Theme = defaults.WebUI.Theme
	}
}


// ParseSize converts a byte count with an optional binary unit suffix.
func ParseSize(value string) (int64, error) {
	value = strings.TrimSpace(strings.ToUpper(value))
	unit := int64(1)
	for _, suffix := range []struct { name string; size int64 }{{"GIB", 1 << 30}, {"GB", 1 << 30}, {"G", 1 << 30}, {"MIB", 1 << 20}, {"MB", 1 << 20}, {"M", 1 << 20}, {"KIB", 1 << 10}, {"KB", 1 << 10}, {"K", 1 << 10}, {"B", 1}} {
		if strings.HasSuffix(value, suffix.name) { unit = suffix.size; value = strings.TrimSpace(strings.TrimSuffix(value, suffix.name)); break }
	}
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil || number <= 0 || number > (1<<40)/unit { return 0, fmt.Errorf("invalid size") }
	return number * unit, nil
}
