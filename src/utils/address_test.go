package utils

import (
	"strings"
	"testing"
)

func TestDisplayAddressNeverUsesUnreachableDefaults(t *testing.T) {
	address := GetDisplayAddress("192.0.2.10")
	if address == "" || address == "localhost" || address == "0.0.0.0" || strings.HasPrefix(address, "127.") {
		t.Fatalf("unsafe display address %q", address)
	}
	if getOutboundIP() == "" {
		t.Log("outbound route is unavailable in this environment")
	}
	if getExternalIP() == "" {
		t.Log("no external-facing interface is available in this environment")
	}
}
