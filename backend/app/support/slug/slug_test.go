package slug

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Hello World", "hello-world"},
		{"  Getting -- Started with Nuxt! ", "getting-started-with-nuxt"},
		{"你好 世界", "你好-世界"},
		{"Go 1.23 release", "go-1-23-release"},
		{"!!!", "item"},
		{"", "item"},
	}
	for _, tt := range tests {
		if got := Slugify(tt.in); got != tt.want {
			t.Errorf("Slugify(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestEnsureUnique(t *testing.T) {
	t.Run("first try free", func(t *testing.T) {
		got, err := EnsureUnique("hello", 191, func(string) (bool, error) { return false, nil })
		if err != nil || got != "hello" {
			t.Fatalf("got %q err %v", got, err)
		}
	})

	t.Run("collision appends random suffix keeping prefix", func(t *testing.T) {
		taken := map[string]bool{"hello": true}
		got, err := EnsureUnique("hello", 191, func(s string) (bool, error) { return taken[s], nil })
		if err != nil {
			t.Fatal(err)
		}
		if len(got) <= len("hello-") || got[:len("hello")] != "hello" || got == "hello" {
			t.Fatalf("expected suffixed variant of hello, got %q", got)
		}
	})

	t.Run("long base is trimmed before suffixing", func(t *testing.T) {
		base := string(make([]byte, 200))
		for i := range base {
			base = base[:i] + "a" + base[i+1:]
		}
		taken := map[string]bool{base: true}
		got, err := EnsureUnique(base, 191, func(s string) (bool, error) { return taken[s], nil })
		if err != nil {
			t.Fatal(err)
		}
		if len(got) > 191 {
			t.Fatalf("slug exceeds max length: %d", len(got))
		}
	})
}

type testErr struct{}

func (*testErr) Error() string { return "boom" }

func TestExistsErrorPropagates(t *testing.T) {
	var boom error = &testErr{}
	if _, err := EnsureUnique("x", 10, func(string) (bool, error) { return false, boom }); err == nil {
		t.Fatal("expected error to propagate")
	}
}
