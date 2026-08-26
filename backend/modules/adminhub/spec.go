package adminhub

import (
	"encoding/json"
	"reflect"

	"github.com/goravel/framework/contracts/database/orm"
)

// Spec declares one admin-manageable resource. Registering a spec is all a
// module needs to do to expose the full CRUD contract consumed by nuxtadmin
// (list/create/read/update/delete/bulk-delete) — see plan §7.3.
type Spec struct {
	// Name is the URL segment, e.g. "articles".
	Name string
	// Model is a zero-value instance of the GORM model.
	Model any
	// Searchable lists the columns the `q` term may hit (whitelist).
	Searchable []string
	// Sortable whitelists sortBy values; falls back to {id, created_at}.
	Sortable []string
	// SortRaw maps extra sort aliases onto fixed SQL fragments (developer
	// constants, never user input), e.g. "hot": "(view_count + reply_count*3)".
	SortRaw map[string]string
	// Fillable restricts which JSON keys create/update accept.
	Fillable []string
	// Hidden removes json keys from every output row (defence in depth on
	// top of model-level json:"-").
	Hidden []string
	// Transform overrides row rendering when provided.
	Transform func(row any) map[string]any
	// SearchOverride replaces the default ILIKE-over-Searchable behaviour for
	// the `q` term (e.g. articles use PG full-text with CJK ILIKE fallback).
	SearchOverride func(query orm.Query, term string) orm.Query
	// BeforeSave may mutate filtered payload data before persistence
	// (e.g. hash "password").
	BeforeSave func(data map[string]any) error
	// AfterSave runs after a successful create/update with the persisted
	// instance and filtered payload (e.g. tag pivot synchronisation).
	AfterSave func(instance any, data map[string]any) error
	// OnShow runs after a successful detail load (e.g. buffered view count).
	OnShow func(instance any)
	// SoftDeletable marks models with a deleted_at column, enabling the
	// trashed=only list filter and the generic restore route.
	SoftDeletable bool
	// PermissionPrefix defaults to Name.
	PermissionPrefix string
}

var registry = map[string]*Spec{}

// Register adds a spec; duplicate names are ignored (first wins) so hot
// re-registration cannot silently swap behaviour.
func Register(spec *Spec) {
	if spec == nil || spec.Name == "" || spec.Model == nil {
		return
	}
	if _, exists := registry[spec.Name]; exists {
		return
	}
	registry[spec.Name] = spec
}

// Get returns the registered spec for a resource name.
func Get(name string) (*Spec, bool) {
	spec, ok := registry[name]
	return spec, ok
}

// tableName resolves the backing table for Table-based recycle-bin queries.
func (s *Spec) tableName() string {
	if t, ok := s.Model.(interface{ TableName() string }); ok {
		return t.TableName()
	}
	return ""
}

/* ---------------- spec helpers ---------------- */

func (s *Spec) prefix() string {
	if s.PermissionPrefix != "" {
		return s.PermissionPrefix
	}
	return s.Name
}

func (s *Spec) sortableColumns() []string {
	if len(s.Sortable) > 0 {
		return s.Sortable
	}
	return []string{"id", "created_at", "updated_at"}
}

func (s *Spec) modelType() reflect.Type {
	return reflect.TypeOf(s.Model)
}

// newInstance allocates a fresh zero-value model instance.
func (s *Spec) newInstance() any {
	return reflect.New(s.modelType()).Interface()
}

// newSlice allocates a *[]*T destination for list queries. The pointer
// wrapper keeps the slice addressable so the ORM can scan rows into it.
func (s *Spec) newSlice() any {
	elem := reflect.PtrTo(s.modelType())
	return reflect.New(reflect.SliceOf(elem)).Interface()
}

func (s *Spec) renderRow(row any) map[string]any {
	if s.Transform != nil {
		out := s.Transform(row)
		for _, h := range s.Hidden {
			delete(out, h)
		}
		return out
	}

	b, err := json.Marshal(row)
	if err != nil {
		return map[string]any{}
	}
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	for _, h := range s.Hidden {
		delete(m, h)
	}
	return m
}

// filterPayload keeps only Fillable keys from an incoming body.
func (s *Spec) filterPayload(body map[string]any) map[string]any {
	allowed := make(map[string]bool, len(s.Fillable))
	for _, f := range s.Fillable {
		allowed[f] = true
	}
	clean := make(map[string]any, len(body))
	for k, v := range body {
		if allowed[k] {
			clean[k] = v
		}
	}
	return clean
}
