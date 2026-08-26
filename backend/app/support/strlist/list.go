package strlist

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// List persists a []string as a JSONB column (role permissions, article AI
// keywords, ...). Shared by all modules so model packages never import each
// other's internals (plan §7.2).
type List []string

func (l List) Value() (driver.Value, error) {
	if l == nil {
		return "[]", nil
	}
	b, err := json.Marshal(l)
	if err != nil {
		return nil, fmt.Errorf("marshal strlist: %w", err)
	}
	return string(b), nil
}

func (l *List) Scan(value any) error {
	switch v := value.(type) {
	case nil:
		*l = List{}
		return nil
	case []byte:
		return json.Unmarshal(v, l)
	case string:
		return json.Unmarshal([]byte(v), l)
	default:
		return fmt.Errorf("unsupported strlist source type %T", value)
	}
}

// Has reports whether the list grants the given permission, honouring the
// wildcard semantics shared with the admin frontend can() matcher:
// "*" grants everything, "prefix.*" grants every permission under prefix.
func (l List) Has(permission string) bool {
	for _, p := range l {
		if p == "*" || p == permission {
			return true
		}
		if len(p) > 1 && p[len(p)-1] == '*' && len(permission) >= len(p)-1 &&
			permission[:len(p)-1] == p[:len(p)-1] {
			return true
		}
	}
	return false
}
