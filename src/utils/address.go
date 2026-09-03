package utils

import (
	"net"
	"os"
	"strings"
)

// GetDisplayAddress returns the most appropriate address to display to users
// Priority: FQDN > specific bind address > external IP > hostname > fallback
// NEVER returns localhost, 127.0.0.1, or 0.0.0.0 per SPEC.md
func GetDisplayAddress(bindAddr string) string {
	// Try to get hostname first and check if it's a valid FQDN
	if hostname, err := os.Hostname(); err == nil && hostname != "" && hostname != "localhost" {
		// Try to resolve hostname to see if it's a valid FQDN
		if addrs, err := net.LookupHost(hostname); err == nil && len(addrs) > 0 {
			return hostname
		}
	}

	// If binding to a specific address (not 0.0.0.0), use it
	if bindAddr != "" && bindAddr != "0.0.0.0" && !strings.HasPrefix(bindAddr, "127.") {
		return bindAddr
	}

	// Try to get outbound IP (most likely accessible IP)
	if externalIP := getOutboundIP(); externalIP != "" {
		return externalIP
	}

	// Try to get external IP from interfaces
	if externalIP := getExternalIP(); externalIP != "" {
		return externalIP
	}

	// Try to get hostname (even if not FQDN)
	if hostname, err := os.Hostname(); err == nil && hostname != "" && hostname != "localhost" {
		return hostname
	}

	// Last resort: generic message (NEVER "localhost")
	return "<your-host>"
}

// getOutboundIP gets the preferred outbound IP of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// getExternalIP attempts to get the external-facing IP address
func getExternalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {
		// Skip down interfaces
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// Skip loopback
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			// Prefer IPv4
			ip = ip.To4()
			if ip != nil {
				return ip.String()
			}
		}
	}

	return ""
}
