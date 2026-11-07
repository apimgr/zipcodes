package server

import (
	"context"
	"net"
	"net/http"
	"strings"
)

type originalPeerKey struct{}

var builtInTrustedProxyCIDRs = []string{
	"127.0.0.0/8", "::1/128", "10.0.0.0/8", "172.16.0.0/12",
	"192.168.0.0/16", "fc00::/7", "169.254.0.0/16", "fe80::/10",
}

func buildTrustedProxies(listenAddress string, additional []string) []*net.IPNet {
	networks := make([]*net.IPNet, 0, len(builtInTrustedProxyCIDRs)+len(additional)+1)
	addCIDR := func(value string) {
		_, network, err := net.ParseCIDR(value)
		if err == nil {
			networks = append(networks, network)
		}
	}
	for _, value := range builtInTrustedProxyCIDRs {
		addCIDR(value)
	}
	if ip := parseListenIP(listenAddress); ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			addCIDR(ipv4.Mask(net.CIDRMask(24, 32)).String() + "/24")
		} else {
			addCIDR(ip.Mask(net.CIDRMask(64, 128)).String() + "/64")
		}
	}
	for _, value := range additional {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if strings.Contains(value, "/") {
			addCIDR(value)
			continue
		}
		if ip := net.ParseIP(value); ip != nil {
			if ipv4 := ip.To4(); ipv4 != nil {
				addCIDR(ipv4.String() + "/32")
			} else {
				addCIDR(ip.String() + "/128")
			}
			continue
		}
		addresses, err := net.LookupIP(value)
		if err != nil {
			continue
		}
		for _, ip := range addresses {
			if ipv4 := ip.To4(); ipv4 != nil {
				addCIDR(ipv4.String() + "/32")
			} else {
				addCIDR(ip.String() + "/128")
			}
		}
	}
	return networks
}

func parseListenIP(address string) net.IP {
	address = strings.TrimSpace(address)
	if host, _, err := net.SplitHostPort(address); err == nil {
		address = host
	}
	address = strings.Trim(address, "[]")
	if address == "" || address == "0.0.0.0" || address == "::" {
		return nil
	}
	return net.ParseIP(address)
}

func (s *Server) isTrustedPeer(r *http.Request) bool {
	peer := r.RemoteAddr
	if original, ok := r.Context().Value(originalPeerKey{}).(string); ok {
		peer = original
	}
	if host, _, err := net.SplitHostPort(peer); err == nil {
		peer = host
	}
	ip := net.ParseIP(peer)
	if ip == nil {
		return false
	}
	for _, network := range s.trustedProxies {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func (s *Server) trustedProxyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalPeer := r.RemoteAddr
		r = r.WithContext(context.WithValue(r.Context(), originalPeerKey{}, originalPeer))
		if s.isTrustedPeer(r) {
			if clientIP := clientIPFromHeaders(r); clientIP != "" {
				r.RemoteAddr = clientIP
			}
		}
		next.ServeHTTP(w, r)
	})
}

func clientIPFromHeaders(r *http.Request) string {
	for _, header := range []string{"CF-Connecting-IP", "True-Client-IP", "X-Real-IP", "X-Forwarded-For", "X-Client-IP"} {
		value := strings.TrimSpace(r.Header.Get(header))
		if header == "X-Forwarded-For" {
			if values := strings.Split(value, ","); len(values) > 0 {
				value = strings.TrimSpace(values[0])
			}
		}
		if host, _, err := net.SplitHostPort(value); err == nil {
			value = host
		}
		if ip := net.ParseIP(value); ip != nil {
			return ip.String()
		}
	}
	return ""
}
