package services

import "testing"

func TestCan(t *testing.T) {
	tests := []struct {
		name        string
		permissions []string
		required    string
		want        bool
	}{
		{"wildcard grants all", []string{"*"}, "anything.at.all", true},
		{"exact match", []string{"posts.view"}, "posts.view", true},
		{"prefix wildcard", []string{"posts.*"}, "posts.delete", true},
		{"prefix wildcard does not leak other prefixes", []string{"posts.*"}, "orders.view", false},
		{"prefix wildcard requires dot boundary", []string{"post.*"}, "posts.view", false},
		{"no permissions denies", nil, "posts.view", false},
		{"empty required denies", []string{"*"}, "", false},
		{"bare star segment only", []string{"*view"}, "users.view", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Can(tt.permissions, tt.required); got != tt.want {
				t.Errorf("Can(%v, %q) = %v, want %v", tt.permissions, tt.required, got, tt.want)
			}
		})
	}
}
