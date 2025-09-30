package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Te8va/shortURL/internal/app/config"
)

// TestTrustedSubnetMiddleware tests the TrustedSubnetMiddleware with various scenarios
func TestTrustedSubnetMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		trustedSubnet  string
		clientIP       string
		xRealIP        string
		remoteAddr     string
		expectedStatus int
		expectedBody   string
		shouldCallNext bool
	}{
		{
			name:           "empty trusted subnet should forbid access",
			trustedSubnet:  "",
			clientIP:       "192.168.1.100",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Forbidden\n",
			shouldCallNext: false,
		},

		{
			name:           "valid IP in subnet via X-Real-IP should allow access",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "192.168.1.100",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			shouldCallNext: true,
		},
		{
			name:           "valid IP in subnet via RemoteAddr should allow access",
			trustedSubnet:  "192.168.1.0/24",
			remoteAddr:     "192.168.1.100:8080",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			shouldCallNext: true,
		},
		{
			name:           "valid IP at subnet boundary should allow access",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "192.168.1.254",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			shouldCallNext: true,
		},

		{
			name:           "IP outside subnet should forbid access",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "10.0.0.100",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Forbidden\n",
			shouldCallNext: false,
		},
		{
			name:           "IP from different subnet should forbid access",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "192.168.2.100",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Forbidden\n",
			shouldCallNext: false,
		},
		{
			name:           "X-Real-IP should take precedence over RemoteAddr",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "192.168.1.100",
			remoteAddr:     "10.0.0.100:8080",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			shouldCallNext: true,
		},
		{
			name:           "larger subnet /16 should allow access",
			trustedSubnet:  "192.168.0.0/16",
			xRealIP:        "192.168.100.100",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			shouldCallNext: true,
		},
		{
			name:           "smaller subnet /30 should allow access to specific IPs",
			trustedSubnet:  "192.168.1.0/30",
			xRealIP:        "192.168.1.1",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			shouldCallNext: true,
		},
		{
			name:           "smaller subnet /30 should forbid access to IPs outside range",
			trustedSubnet:  "192.168.1.0/30",
			xRealIP:        "192.168.1.5",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Forbidden\n",
			shouldCallNext: false,
		},
		{
			name:           "IPv6 in subnet should allow access",
			trustedSubnet:  "2001:db8::/32",
			xRealIP:        "2001:db8::1",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			shouldCallNext: true,
		},
		{
			name:           "IPv6 outside subnet should forbid access",
			trustedSubnet:  "2001:db8::/32",
			xRealIP:        "2001:db9::1",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Forbidden\n",
			shouldCallNext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			cfg := &config.Config{
				TrustedSubnet: tt.trustedSubnet,
			}

			middleware := TrustedSubnetMiddleware(cfg)(nextHandler)

			req := httptest.NewRequest("GET", "/api/internal/stats", nil)

			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}
			if tt.remoteAddr != "" {
				req.RemoteAddr = tt.remoteAddr
			}

			rr := httptest.NewRecorder()

			middleware.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}

			if nextCalled != tt.shouldCallNext {
				t.Errorf("expected next called %v, got %v", tt.shouldCallNext, nextCalled)
			}
		})
	}
}

// TestIsIPInSubnet tests the isIPInSubnet function with various IP and subnet combinations
func TestIsIPInSubnet(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		subnet   string
		expected bool
		wantErr  bool
	}{
		{
			name:     "valid IPv4 in subnet",
			ip:       "192.168.1.100",
			subnet:   "192.168.1.0/24",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "valid IPv4 outside subnet",
			ip:       "10.0.0.100",
			subnet:   "192.168.1.0/24",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "valid IPv6 in subnet",
			ip:       "2001:db8::1",
			subnet:   "2001:db8::/32",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "valid IPv6 outside subnet",
			ip:       "2001:db9::1",
			subnet:   "2001:db8::/32",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "invalid IP address",
			ip:       "invalid-ip",
			subnet:   "192.168.1.0/24",
			expected: false,
			wantErr:  true,
		},
		{
			name:     "invalid subnet CIDR",
			ip:       "192.168.1.100",
			subnet:   "invalid-subnet",
			expected: false,
			wantErr:  true,
		},
		{
			name:     "empty IP address",
			ip:       "",
			subnet:   "192.168.1.0/24",
			expected: false,
			wantErr:  true,
		},
		{
			name:     "IP at subnet network address",
			ip:       "192.168.1.0",
			subnet:   "192.168.1.0/24",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "IP at subnet broadcast address",
			ip:       "192.168.1.255",
			subnet:   "192.168.1.0/24",
			expected: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := isIPInSubnet(tt.ip, tt.subnet)

			if (err != nil) != tt.wantErr {
				t.Errorf("isIPInSubnet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result != tt.expected {
				t.Errorf("isIPInSubnet() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestGetClientIP tests the getClientIP function with different header scenarios
func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name       string
		xRealIP    string
		remoteAddr string
		expected   string
	}{
		{
			name:       "X-Real-IP should be preferred",
			xRealIP:    "192.168.1.100",
			remoteAddr: "10.0.0.100:8080",
			expected:   "192.168.1.100",
		},
		{
			name:       "should use RemoteAddr when X-Real-IP is empty",
			xRealIP:    "",
			remoteAddr: "10.0.0.100:8080",
			expected:   "10.0.0.100",
		},
		{
			name:       "should handle IPv6 RemoteAddr",
			xRealIP:    "",
			remoteAddr: "[2001:db8::1]:8080",
			expected:   "2001:db8::1",
		},
		{
			name:       "should handle malformed RemoteAddr",
			xRealIP:    "",
			remoteAddr: "malformed-address",
			expected:   "malformed-address",
		},
		{
			name:       "should handle empty RemoteAddr port",
			xRealIP:    "",
			remoteAddr: "192.168.1.100:",
			expected:   "192.168.1.100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}
			if tt.remoteAddr != "" {
				req.RemoteAddr = tt.remoteAddr
			}

			result := getClientIP(req)
			if result != tt.expected {
				t.Errorf("getClientIP() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestTrustedSubnetMiddlewareErrorHandling tests error scenarios in the middleware
func TestTrustedSubnetMiddlewareErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		trustedSubnet string
		xRealIP       string
		expectedBody  string
	}{
		{
			name:          "invalid IP should return internal server error",
			trustedSubnet: "192.168.1.0/24",
			xRealIP:       "invalid-ip",
			expectedBody:  "Internal server error\n",
		},
		{
			name:          "invalid subnet should return internal server error",
			trustedSubnet: "invalid-subnet",
			xRealIP:       "192.168.1.100",
			expectedBody:  "Internal server error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("next handler should not be called on error")
			})

			cfg := &config.Config{
				TrustedSubnet: tt.trustedSubnet,
			}

			middleware := TrustedSubnetMiddleware(cfg)(nextHandler)

			req := httptest.NewRequest("GET", "/api/internal/stats", nil)
			req.Header.Set("X-Real-IP", tt.xRealIP)

			rr := httptest.NewRecorder()
			middleware.ServeHTTP(rr, req)

			if rr.Code != http.StatusInternalServerError {
				t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
