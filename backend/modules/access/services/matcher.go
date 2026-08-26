package services

// Can mirrors the nuxtadmin frontend can() semantics exactly:
//   - "*"            grants every permission
//   - "prefix.*"     grants everything under prefix (e.g. "articles.*")
//   - exact match    otherwise
//
// Both sides must stay in sync; see docs/工程化开发计划.md §5.2.
func Can(permissions []string, required string) bool {
	if required == "" {
		return false
	}
	for _, p := range permissions {
		if p == "*" || p == required {
			return true
		}
		if ok, _ := matchPrefix(p, required); ok {
			return true
		}
	}
	return false
}

func matchPrefix(pattern, required string) (bool, bool) {
	if len(pattern) < 2 || pattern[len(pattern)-1] != '*' {
		return false, false
	}
	prefix := pattern[:len(pattern)-1]
	if prefix == "" {
		return false, false
	}
	if len(required) >= len(prefix) && required[:len(prefix)] == prefix {
		return true, true
	}
	return false, true
}
