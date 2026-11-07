package main

import (
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, run func()) string {
	t.Helper()
	previous := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	run()
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = previous
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	return string(output)
}

func TestCLIHelpAndVersion(t *testing.T) {
	oldArgs, oldFlags := os.Args, flag.CommandLine
	t.Cleanup(func() { os.Args, flag.CommandLine = oldArgs, oldFlags })
	flag.CommandLine = flag.NewFlagSet("zipcodes", flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
	os.Args = []string{"zipcodes", "--help"}
	if output := captureStdout(t, main); !strings.Contains(output, "Usage: zipcodes") {
		t.Fatalf("help output = %q", output)
	}
	flag.CommandLine = flag.NewFlagSet("zipcodes", flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
	os.Args = []string{"zipcodes", "--version"}
	if output := captureStdout(t, main); !strings.Contains(output, Version) {
		t.Fatalf("version output = %q", output)
	}
}

func TestCommandHelpers(t *testing.T) {
	for _, command := range []string{"install", "uninstall", "start", "stop", "restart", "status", "unknown"} {
		if output := captureStdout(t, func() { handleServiceCommand(command) }); output == "" {
			t.Errorf("service command %q returned no help", command)
		}
	}
	for _, mode := range []string{"on", "off", "unknown"} {
		if output := captureStdout(t, func() { handleMaintenanceMode(mode) }); output == "" {
			t.Errorf("maintenance mode %q returned no output", mode)
		}
	}
	if output := captureStdout(t, printHelp); !strings.Contains(output, "Documentation:") {
		t.Fatalf("help output = %q", output)
	}
}

func TestServerStatus(t *testing.T) {
	oldPort := os.Getenv("PORT")
	t.Cleanup(func() { _ = os.Setenv("PORT", oldPort) })
	if err := os.Unsetenv("PORT"); err != nil {
		t.Fatal(err)
	}
	if status := checkServerStatus(); status != 1 {
		t.Fatalf("missing PORT status = %d", status)
	}
	listener := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}))
	defer listener.Close()
	port := strings.TrimPrefix(listener.URL, "http://127.0.0.1:")
	if port == listener.URL {
		port = strings.TrimPrefix(listener.URL, "http://")
	}
	if err := os.Setenv("PORT", port); err != nil {
		t.Fatal(err)
	}
	if status := checkServerStatus(); status != 0 {
		t.Fatalf("healthy server status = %d", status)
	}
	if status := checkServerStatus(port); status != 0 {
		t.Fatalf("healthy server status from port override = %d", status)
	}
	listener.Close()
	if status := checkServerStatus(); status != 1 {
		t.Fatalf("stopped server status = %d", status)
	}
}
