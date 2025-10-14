package utils

import (
	"net"
	"os"
)

// GetDisplayAddress returns the most appropriate address to display to users
// Priority: specific bind address > external IP > hostname > fallback
func GetDisplayAddress(bindAddr string) string {
	// If binding to a specific address (not 0.0.0.0), use it
	if bindAddr != "" && bindAddr != "0.0.0.0" {
		return bindAddr
	}

	// Try to get external IP
	if externalIP := getExternalIP(); externalIP != "" {
		return externalIP
	}

	// Try to get hostname
	if hostname, err := os.Hostname(); err == nil && hostname != "" {
		return hostname
	}

	// Fallback only if nothing else works
	return "localhost"
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
