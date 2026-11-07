package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestServiceIdentityAndCopy(t *testing.T) {
	if GetBinaryPath() == "" {
		t.Fatal("binary path is empty")
	}
	if got := DetectServiceManager(); got < ServiceUnknown || got > ServiceBSDRC {
		t.Fatalf("unknown manager enum %d", got)
	}
	dir := t.TempDir()
	source := filepath.Join(dir, "source")
	destination := filepath.Join(dir, "nested", "binary")
	if err := os.WriteFile(source, []byte("payload"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := copyBinary(source, destination); err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(destination); err != nil || string(got) != "payload" {
		t.Fatalf("copied content = %q, %v", got, err)
	}
	if err := copyBinary(filepath.Join(dir, "missing"), filepath.Join(dir, "out")); err == nil {
		t.Fatal("copyBinary accepted missing source")
	}
}

func TestInstallersWriteOnlyUnderConfiguredRoot(t *testing.T) {
	root := t.TempDir()
	oldRoot := installRoot
	installRoot = root
	t.Cleanup(func() { installRoot = oldRoot })

	binDir := filepath.Join(root, "commands")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatal(err)
	}
	commandLog := filepath.Join(root, "commands.log")
	command := "#!/bin/sh\nprintf '%s\\n' \"$0 $*\" >> \"$COMMAND_LOG\"\n"
	for _, name := range []string{"systemctl", "sv", "launchctl", "sc.exe", "service"} {
		path := filepath.Join(binDir, name)
		if err := os.WriteFile(path, []byte(command), 0755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("COMMAND_LOG", commandLog)
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	cases := []struct {
		name    string
		manager ServiceType
		path    string
	}{
		{"systemd", ServiceSystemd, "/etc/systemd/system/zipcodes.service"},
		{"runit", ServiceRunit, "/etc/sv/zipcodes/run"},
		{"launchd", ServiceLaunchd, "/Library/LaunchDaemons/com.apimgr.zipcodes.plist"},
		{"windows", ServiceWindows, "/usr/local/bin/zipcodes"},
		{"bsd", ServiceBSDRC, "/usr/local/etc/rc.d/zipcodes"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			oldRoot := installRoot
			oldManager := detectServiceManager
			installRoot = filepath.Join(root, tc.name)
			detectServiceManager = func() ServiceType { return tc.manager }
			t.Cleanup(func() {
				installRoot = oldRoot
				detectServiceManager = oldManager
			})
			if err := Install(); err != nil {
				t.Fatalf("install: %v", err)
			}
			path := rootedPath(tc.path)
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("installed artifact %q: %v", path, err)
			}
			if err := Uninstall(); err != nil {
				t.Fatalf("uninstall: %v", err)
			}
		})
	}
	if _, err := os.Stat(commandLog); err != nil {
		t.Fatalf("service commands were not invoked: %v", err)
	}
}

func TestRootedPathKeepsAbsolutePathsWithinTestRoot(t *testing.T) {
	oldRoot := installRoot
	installRoot = t.TempDir()
	t.Cleanup(func() { installRoot = oldRoot })
	path := rootedPath("/etc/apimgr/zipcodes/server.yml")
	if !filepath.IsAbs(path) || !strings.HasPrefix(path, installRoot+string(filepath.Separator)) {
		t.Fatalf("rooted path = %q, root = %q", path, installRoot)
	}
}

func TestServiceManagerControlsInvokeCommands(t *testing.T) {
	root := t.TempDir()
	binDir := filepath.Join(root, "commands")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatal(err)
	}
	commandLog := filepath.Join(root, "commands.log")
	command := "#!/bin/sh\nprintf '%s\\n' \"$0 $*\" >> \"$COMMAND_LOG\"\n"
	for _, name := range []string{"systemctl", "sv", "launchctl", "sc.exe", "service"} {
		if err := os.WriteFile(filepath.Join(binDir, name), []byte(command), 0755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("COMMAND_LOG", commandLog)
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	oldDetector := detectServiceManager
	defer func() { detectServiceManager = oldDetector }()
	for _, manager := range []ServiceType{ServiceSystemd, ServiceRunit, ServiceLaunchd, ServiceWindows, ServiceBSDRC} {
		detectServiceManager = func() ServiceType { return manager }
		for name, action := range map[string]func() error{
			"start": Start, "stop": Stop, "restart": Restart, "reload": Reload,
		} {
			if err := action(); err != nil {
				t.Errorf("manager %d %s: %v", manager, name, err)
			}
		}
	}
	if _, err := os.Stat(commandLog); err != nil {
		t.Fatalf("service commands were not invoked: %v", err)
	}
}
