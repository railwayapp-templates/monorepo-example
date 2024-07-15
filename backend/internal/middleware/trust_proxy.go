package middleware

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"
)

// PrivateRangesCIDR returns a list of private CIDR range
// strings, which can be used as a TrustProxyConfiguration shortcut.
//
// 192.168.0.0/16, 172.16.0.0/12, "100.0.0.0/12", 127.0.0.1/8, 10.0.0.0/8, fd00::/8, ::1
var privateRanges = []string{
	"192.168.0.0/16",
	"172.16.0.0/12",
	"100.0.0.0/8",
	"127.0.0.1/8",
	"10.0.0.0/8",
	"fd00::/8",
	"::1",
}

var proxyIPHeaders = []string{
	"Fastly-Client-Ip",
	"CF-Connecting-IP",
	"X-Envoy-External-Address",
	"X-Forwarded-For",
	"X-Real-IP",
	"True-Client-IP",
}

var schemeHeaders = []string{
	"X-Forwarded-Proto",
	"X-Forwarded-Scheme",
}

var xForwardedHosts = []string{
	"X-Forwarded-Host",
}

type trustProxyConfig struct {
	TrustIPRanges       []string
	TrustIPHeaders      []string
	TrustSchemeHeaders  []string
	TrustForwardedHosts []string
	ErrorLogger         *slog.Logger
}

func (c *trustProxyConfig) loadDefaults() {
	if len(c.TrustIPRanges) == 0 {
		c.TrustIPRanges = privateRanges
	}

	if len(c.TrustIPHeaders) == 0 {
		c.TrustIPHeaders = proxyIPHeaders
	}

	if len(c.TrustSchemeHeaders) == 0 {
		c.TrustSchemeHeaders = schemeHeaders
	}

	if len(c.TrustForwardedHosts) == 0 {
		c.TrustForwardedHosts = xForwardedHosts
	}

	if c.ErrorLogger == nil {
		c.ErrorLogger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))
	}
}

// TrustProxy checks if the request IP matches one of the provided ranges/IPs
// then inspects common reverse proxy headers and sets the corresponding
// fields in the HTTP request struct for use by middleware or handlers that are next
func TrustProxy(c *trustProxyConfig) func(http.Handler) http.Handler {
	c.loadDefaults()

	// parse passed in trusted IPs into a 'netip.Prefix' slice
	parsedIPs, err := parseIPRanges(c.TrustIPRanges)
	if err != nil {
		c.ErrorLogger.Error(err.Error())
		os.Exit(1)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// check if RemoteAddr is trusted
			trusted, err := isTrustedIP(r.RemoteAddr, parsedIPs)
			if err != nil {
				c.ErrorLogger.Warn(err.Error(), slog.String("ip", r.RemoteAddr))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// if RemoteAddr is not trusted, serve next and return early without trusting any proxy headers
			if !trusted {
				next.ServeHTTP(w, r)
				return
			}

			// RemoteAddr is trusted

			// Set the RemoteAddr with the value passed by the proxy
			if realIP := c.getRealIP(r.Header); realIP != "" {
				r.RemoteAddr = realIP
			}

			// Set the host with the value passed by the proxy
			if realHost := c.getRealHost(r.Header); realHost != "" {
				r.Host = realHost
			}

			// Set the scheme with the value passed by the proxy
			if scheme := c.getScheme(r.Header); scheme != "" {
				r.URL.Scheme = scheme
			}

			next.ServeHTTP(w, r)
		})
	}
}

// parse passed in IP ranges into a 'netip.Prefix' slice
func parseIPRanges(IPRanges []string) ([]netip.Prefix, error) {
	var parsedIPs []netip.Prefix

	for _, ipStr := range IPRanges {
		if strings.Contains(ipStr, "/") {
			ipNet, err := netip.ParsePrefix(ipStr)
			if err != nil {
				return nil, fmt.Errorf("parsing CIDR expression: %w", err)
			}

			parsedIPs = append(parsedIPs, ipNet)
		} else {
			ipAddr, err := netip.ParseAddr(ipStr)
			if err != nil {
				return nil, fmt.Errorf("invalid IP address: '%s': %w", ipStr, err)
			}

			parsedIPs = append(parsedIPs, netip.PrefixFrom(ipAddr, ipAddr.BitLen()))
		}
	}

	return parsedIPs, nil
}

// check if RemoteAddr is trusted
func isTrustedIP(remoteAddr string, trustedIPs []netip.Prefix) (bool, error) {
	ipStr, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		ipStr = remoteAddr
	}

	ipAddr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return false, err
	}

	for _, ipRange := range trustedIPs {
		if ipRange.Contains(ipAddr) {
			return true, nil
		}
	}

	return false, nil
}

// get the real IP from the proxy headers if present
func (c *trustProxyConfig) getRealIP(headers http.Header) string {
	for _, proxyHeader := range c.TrustIPHeaders {
		if value := headers.Get(proxyHeader); value != "" {
			ips := strings.Split(value, ",")
			return strings.TrimSpace(ips[len(ips)-1])
		}
	}

	return ""
}

func (c *trustProxyConfig) getRealHost(headers http.Header) string {
	for _, hostHeader := range c.TrustForwardedHosts {
		if value := headers.Get(hostHeader); value != "" {
			return value
		}
	}

	return ""
}

// get the scheme from the proxy headers if present
func (c *trustProxyConfig) getScheme(headers http.Header) string {
	for _, schemaHeader := range c.TrustSchemeHeaders {
		if value := headers.Get(schemaHeader); value != "" {
			return strings.ToLower(value)
		}
	}

	return ""
}
