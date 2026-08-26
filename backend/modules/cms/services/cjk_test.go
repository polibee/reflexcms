package services

import "testing"

func TestContainsCJK(t *testing.T) {
	positive := []string{"你好", "Nuxt 3 中文指南", "テスト", "한국어", "mixed English 中文"}
	for _, s := range positive {
		if !ContainsCJK(s) {
			t.Errorf("ContainsCJK(%q) = false, want true", s)
		}
	}
	negative := []string{"", "plain english", "nuxt-3-release", "123 !@#"}
	for _, s := range negative {
		if ContainsCJK(s) {
			t.Errorf("ContainsCJK(%q) = true, want false", s)
		}
	}
}
