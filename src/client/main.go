package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	version    = "dev"
	zipPattern = regexp.MustCompile(`^[0-9]{5}$`)
)

const (
	projectName = "zipcodes"
	internalOrg = "apimgr"
)

type cliConfig struct {
	Server struct {
		Primary string `yaml:"primary"`
	} `yaml:"server"`
}

type lookupResponse struct {
	Success bool `json:"ok"`
	Data    struct {
		Zip       string   `json:"zip"`
		ZipCode   int      `json:"zip_code"`
		City      string   `json:"city"`
		State     string   `json:"state"`
		StateFull string   `json:"state_full"`
		County    string   `json:"county"`
		Latitude  any      `json:"latitude"`
		Longitude any      `json:"longitude"`
		Timezone  string   `json:"timezone"`
		AreaCodes []string `json:"area_codes"`
	} `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	name := filepath.Base(os.Args[0])
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(stderr)
	serverFlag := flags.String("server", "", "Server base URL")
	showVersion := flags.Bool("version", false, "Show version")
	flags.Usage = func() {
		_, _ = fmt.Fprintf(stderr, "Usage: %s [--server URL] ZIPCODE\n", name)
		_, _ = fmt.Fprintln(stderr, "Look up a five-digit US ZIP code using the Zipcodes API.")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if *showVersion {
		_, err := fmt.Fprintf(stdout, "%s %s\n", name, version)
		return err
	}
	if flags.NArg() != 1 || !zipPattern.MatchString(flags.Arg(0)) {
		flags.Usage()
		return errors.New("ZIPCODE must be a five-digit US ZIP code")
	}

	cfg, configPath, err := loadConfig()
	if err != nil {
		return fmt.Errorf("load CLI config: %w", err)
	}
	server := strings.TrimSpace(*serverFlag)
	if server == "" {
		server = strings.TrimSpace(os.Getenv("ZIPCODES_SERVER_PRIMARY"))
	}
	if server == "" {
		server = strings.TrimSpace(cfg.Server.Primary)
	}
	if server == "" {
		return errors.New("server URL is required; set --server or ZIPCODES_SERVER_PRIMARY")
	}
	parsed, err := url.Parse(server)
	if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "http") || parsed.Host == "" || parsed.User != nil {
		return errors.New("server URL must be an absolute HTTP or HTTPS URL without credentials")
	}
	if cfg.Server.Primary == "" && *serverFlag != "" {
		cfg.Server.Primary = server
		if err := saveConfig(cfg, configPath); err != nil {
			return fmt.Errorf("save CLI config: %w", err)
		}
	}

	endpoint := strings.TrimRight(server, "/") + "/api/v1/zipcodes/" + flags.Arg(0)
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", projectName+"-cli/"+version)
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("cannot connect to server: %w", err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		var body lookupResponse
		_ = json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&body)
		if body.Error.Message != "" {
			return fmt.Errorf("server returned %s: %s", response.Status, body.Error.Message)
		}
		return fmt.Errorf("server returned %s", response.Status)
	}
	var body lookupResponse
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&body); err != nil {
		return fmt.Errorf("decode server response: %w", err)
	}
	if !body.Success {
		return errors.New("server returned an unsuccessful lookup response")
	}
	_, err = fmt.Fprintf(stdout, "%s  %s, %s\nCounty: %s\nCoordinates: %v, %v\nTimezone: %s\nArea codes: %s\n", body.Data.Zip, body.Data.City, body.Data.State, body.Data.County, body.Data.Latitude, body.Data.Longitude, body.Data.Timezone, strings.Join(body.Data.AreaCodes, ", "))
	return err
}

func loadConfig() (*cliConfig, string, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return nil, "", err
	}
	path := filepath.Join(root, internalOrg, projectName, "cli.yml")
	cfg := &cliConfig{}
	data, err := os.ReadFile(path)
	if err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, "", err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, "", err
	}
	return cfg, path, nil
}

func saveConfig(cfg *cliConfig, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	if err := file.Chmod(0600); err != nil {
		return errors.Join(err, file.Close())
	}
	if _, err := file.Write(data); err != nil {
		return errors.Join(err, file.Close())
	}
	return file.Close()
}
