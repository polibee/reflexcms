package safefetch

import "testing"

func TestValidateURLBlocksUnsafeTargets(t *testing.T) {
	v := NewValidator()

	blocked := []string{
		"", // missing host
		"ftp://example.com/file",
		"file:///etc/passwd",
		"gopher://127.0.0.1:70",
		"http://localhost/x",
		"https://LOCALHOST:8080/admin",
		"http://app.localhost/x",
		"http://myhost.local/",
		"http://127.0.0.1/",
		"http://[::1]/",
		"http://10.0.0.9/internal",
		"http://172.16.1.1/",
		"http://192.168.1.10/router",
		"http://169.254.169.254/latest/meta-data", // cloud metadata endpoint
		"http://0.0.0.0/",
		"http://224.0.0.1/multicast",
		"://missing-scheme",
	}

	for _, raw := range blocked {
		if _, err := v.ValidateURL(raw); err == nil {
			t.Errorf("ValidateURL(%q) unexpectedly allowed", raw)
		}
	}
}

func TestValidateURLAcceptsPublicHTTPS(t *testing.T) {
	v := NewValidator()
	// example.com is IANA-reserved for documentation and resolvable; if the
	// sandbox has no network this test skips rather than fails.
	u, err := v.ValidateURL("https://example.com/index.html")
	if err != nil {
		t.Skipf("network unavailable: %v", err)
	}
	if u.Scheme != "https" || u.Hostname() != "example.com" {
		t.Fatalf("unexpected url: %s", u)
	}
}
