// Package slug derives URL-safe slugs and guarantees uniqueness against a
// caller-supplied lookup, keeping the logic pure and unit-testable.
package slug

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"unicode"
)

// Slugify converts a title into a URL-safe slug. Unicode letters (incl. CJK)
// are preserved so Chinese titles keep readable slugs; everything else
// collapses to '-'. An empty result falls back to "item".
func Slugify(title string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.TrimSpace(strings.ToLower(title)) {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if dash && b.Len() > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(r)
			dash = false
		default:
			dash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "item"
	}
	if len(out) > 180 {
		out = out[:180]
	}
	return out
}

// EnsureUnique appends a short random suffix on slug collisions. The exists
// callback supplies the database lookup so tests can run offline.
func EnsureUnique(baseSlug string, maxLen int, exists func(slug string) (bool, error)) (string, error) {
	slug := baseSlug
	for attempt := 0; attempt < 6; attempt++ {
		taken, err := exists(slug)
		if err != nil {
			return "", fmt.Errorf("check slug %q: %w", slug, err)
		}
		if !taken {
			return slug, nil
		}
		slug = truncateWithSuffix(baseSlug, randomSuffix(), maxLen)
	}
	return "", fmt.Errorf("could not derive unique slug from %q", baseSlug)
}

func truncateWithSuffix(base, suffix string, maxLen int) string {
	limit := maxLen - len(suffix) - 1
	if limit < 1 {
		return suffix
	}
	base = strings.TrimRight(base, "-")
	if len(base) > limit {
		base = base[:limit]
	}
	return base + "-" + suffix
}

func randomSuffix() string {
	raw := make([]byte, 3)
	if _, err := rand.Read(raw); err != nil {
		return fmt.Sprintf("%x", time.Now().UnixNano()%1e9)
	}
	return hex.EncodeToString(raw)
}
