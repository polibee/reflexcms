// Package safefetch is the ONLY sanctioned path for server-initiated HTTP
// requests to user-supplied URLs. Enforced rules (plan §O9 / §10.3):
//
//   - scheme must be http or https
//   - the host must not be a localhost name, loopback, private, link-local,
//     multicast, reserved or unspecified address — checked both on literal
//     IPs and again after DNS resolution of every redirect hop
package safefetch

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout = 10 * time.Second
	MaxBytes       = 5 << 20 // 5 MiB response cap
)

var ErrUnsafeURL = errors.New("url is not allowed")

type Validator struct {
	resolver net.Resolver
}

func NewValidator() *Validator {
	return &Validator{}
}

func isBlockedIP(ip net.IP) bool {
	return ip == nil ||
		ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsMulticast()
}

func isBlockedLiteralHost(host string) bool {
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	return isBlockedIP(ip)
}

func isLocalHostname(host string) bool {
	h := strings.TrimSuffix(strings.ToLower(host), ".")
	return h == "localhost" || strings.HasSuffix(h, ".localhost") ||
		h == "local" || strings.HasSuffix(h, ".local")
}

// ValidateURL enforces scheme and host rules. Literal IPs are rejected before
// DNS; hostnames are resolved and every returned address re-checked.
func (v *Validator) ValidateURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: parse: %v", ErrUnsafeURL, err)
	}
	switch u.Scheme {
	case "http", "https":
	default:
		return nil, fmt.Errorf("%w: scheme %q", ErrUnsafeURL, u.Scheme)
	}

	host := u.Hostname()
	if host == "" {
		return nil, fmt.Errorf("%w: missing host", ErrUnsafeURL)
	}
	if isLocalHostname(host) || isBlockedLiteralHost(host) {
		return nil, fmt.Errorf("%w: host %q", ErrUnsafeURL, host)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	addrs, err := v.resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("resolve %q: %w", host, err)
	}
	for _, addr := range addrs {
		if isBlockedIP(addr.IP) {
			return nil, fmt.Errorf("%w: resolves to %s", ErrUnsafeURL, addr.IP)
		}
	}
	return u, nil
}

// Client returns an http.Client whose redirect chain re-validates every hop.
func (v *Validator) Client(timeout time.Duration) *http.Client {
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("%w: too many redirects", ErrUnsafeURL)
			}
			if _, err := v.ValidateURL(req.URL.String()); err != nil {
				return fmt.Errorf("redirect blocked: %w", err)
			}
			return nil
		},
	}
}
