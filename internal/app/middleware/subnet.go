package middleware

import (
	"fmt"
	"net"
	"net/http"

	"github.com/Te8va/shortURL/internal/app/config"
)

// TrustedSubnetMiddleware creates middleware that checks if client IP is in trusted subnet
func TrustedSubnetMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.TrustedSubnet == "" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			clientIP := getClientIP(r)
			isAllowed, err := isIPInSubnet(clientIP, cfg.TrustedSubnet)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !isAllowed {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isIPInSubnet(ipStr, subnetCIDR string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	_, subnet, err := net.ParseCIDR(subnetCIDR)
	if err != nil {
		return false, fmt.Errorf("invalid trusted subnet: %w", err)
	}

	return subnet.Contains(ip), nil
}

func getClientIP(r *http.Request) string {
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
