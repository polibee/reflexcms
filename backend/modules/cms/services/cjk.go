package services

import "unicode"

// ContainsCJK reports whether the term includes CJK ideographs, kana or
// hangul. PG 'simple' full-text search cannot segment CJK text into lexemes,
// so the article search degrades to ILIKE for those terms (plan §M2).
func ContainsCJK(s string) bool {
	for _, r := range s {
		switch {
		case unicode.Is(unicode.Han, r): // Chinese
			return true
		case unicode.Is(unicode.Hiragana, r), unicode.Is(unicode.Katakana, r):
			return true
		case unicode.Is(unicode.Hangul, r):
			return true
		}
	}
	return false
}
