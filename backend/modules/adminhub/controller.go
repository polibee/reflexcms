package adminhub

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/http"

	httpx "reflexcms/backend/app/support/httpx"
	accessmatcher "reflexcms/backend/modules/access/services"
	auditservices "reflexcms/backend/modules/audit/services"
	authcontrollers "reflexcms/backend/modules/auth/http/controllers"
)

type GatewayController struct{}

func NewGatewayController() *GatewayController { return &GatewayController{} }

/* ---------------- GET /api/admin/{resource} ---------------- */

func (r *GatewayController) Index(ctx http.Context) http.Response {
	spec, ok := currentSpec(ctx)
	if !ok {
		return httpx.Error(ctx, 404, unknownResource(ctx))
	}
	if resp := requirePermission(ctx, spec.prefix()+".view"); resp != nil {
		return *resp
	}

	page := clampInt(ctx.Request().Query("page", "1"), 1, 1<<30)
	perPage := clampInt(ctx.Request().Query("perPage", "10"), 1, 200)

	// Soft-deletable specs go through a Table-based query so the recycle bin
	// can opt into trashed rows (gorm's implicit "deleted_at IS NULL" scope
	// does not apply to Table queries).
	trashedMode := ctx.Request().Query("trashed", "")
	var query orm.Query
	if spec.SoftDeletable && spec.tableName() != "" {
		query = facadesQuery().Table(spec.tableName())
		switch trashedMode {
		case "only":
			query = query.Where("deleted_at IS NOT NULL")
		case "all":
			// no filter
		default:
			query = query.Where("deleted_at IS NULL")
		}
	} else {
		query = facadesQuery().Model(spec.Model)
	}

	if term := strings.TrimSpace(ctx.Request().Query("q", "")); term != "" {
		if spec.SearchOverride != nil {
			query = spec.SearchOverride(query, term)
		} else if len(spec.Searchable) > 0 {
			// Column names come from the Searchable whitelist (never raw
			// input); the search term itself is always a bound parameter.
			conds := make([]string, 0, len(spec.Searchable))
			args := make([]any, 0, len(spec.Searchable))
			for _, col := range spec.Searchable {
				conds = append(conds, "lower("+col+") LIKE ?")
				args = append(args, "%"+strings.ToLower(term)+"%")
			}
			query = query.Where("("+strings.Join(conds, " OR ")+")", args...)
		}
	}

	total, err := query.Count()
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	sortBy := ctx.Request().Query("sortBy", "id")
	dir := strings.ToLower(ctx.Request().Query("sortDir", "asc"))
	if dir != "asc" && dir != "desc" {
		dir = "asc"
	}

	switch {
	case contains(spec.sortableColumns(), sortBy):
		query = query.OrderBy(sortBy, dir)
	default:
		// SortRaw aliases (developer-defined fragments) come next; anything
		// else falls back to id so external input can never reach ORDER BY.
		if frag, ok := spec.SortRaw[sortBy]; ok && frag != "" {
			query = query.OrderBy("("+frag+")", dir)
		} else {
			query = query.OrderBy("id", dir)
		}
	}

	// Recycle-bin listing: scanning into a typed destination would let gorm
	// re-apply its soft-delete scope, so trashed rows are read as raw maps.
	if trashedMode == "only" {
		var rows []map[string]any
		if err := query.Select("*").
			Offset((page - 1) * perPage).
			Limit(perPage).
			Get(&rows); err != nil {
			return httpx.Error(ctx, 500, err.Error())
		}
		for _, row := range rows {
			delete(row, "password")
			for _, h := range spec.Hidden {
				delete(row, h)
			}
		}
		return ctx.Response().Success().Json(http.Json{
			"items":      rows,
			"total":      total,
			"page":       page,
			"perPage":    perPage,
			"totalPages": maxInt64((total+int64(perPage)-1)/int64(perPage), 1),
		})
	}

	rows := spec.newSlice()
	if err := query.Offset((page - 1) * perPage).Limit(perPage).Get(rows); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	items := renderSlice(spec, rows)
	return ctx.Response().Success().Json(http.Json{
		"items":      items,
		"total":      total,
		"page":       page,
		"perPage":    perPage,
		"totalPages": maxInt64((total+int64(perPage)-1)/int64(perPage), 1),
	})
}

/* ---------------- POST /api/admin/{resource} ---------------- */

func (r *GatewayController) Store(ctx http.Context) http.Response {
	spec, ok := currentSpec(ctx)
	if !ok {
		return httpx.Error(ctx, 404, unknownResource(ctx))
	}
	if resp := requirePermission(ctx, spec.prefix()+".create"); resp != nil {
		return *resp
	}

	body := map[string]any{}
	if err := ctx.Request().Bind(&body); err != nil {
		return httpx.Error(ctx, 422, "Invalid JSON body")
	}
	data := spec.filterPayload(body)
	if spec.BeforeSave != nil {
		if err := spec.BeforeSave(data); err != nil {
			return httpx.Error(ctx, 422, err.Error())
		}
	}

	model := spec.newInstance()
	if err := applyData(model, data); err != nil {
		return httpx.Error(ctx, 422, err.Error())
	}
	if err := facadesQuery().Create(model); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	if resp := runAfterSave(ctx, spec, model, data); resp != nil {
		return *resp
	}
	auditservices.Record(ctx, CurrentUserID(ctx), "create", spec.Name, getID(model), http.StatusOK, data)

	return ctx.Response().Success().Json(spec.renderRow(deref(model)))
}

/* ---------------- GET /api/admin/{resource}/{id} ---------------- */

func (r *GatewayController) Show(ctx http.Context) http.Response {
	spec, ok := currentSpec(ctx)
	if !ok {
		return httpx.Error(ctx, 404, unknownResource(ctx))
	}
	if resp := requirePermission(ctx, spec.prefix()+".view"); resp != nil {
		return *resp
	}

	id, valid := parseID(ctx)
	if !valid {
		return httpx.Error(ctx, 404, notFoundMessage(spec))
	}

	model := spec.newInstance()
	// FindOrFail (not Find): gorm's Find does not error on empty results,
	// which would leak zero-value rows for deleted/unknown ids.
	if err := facadesQuery().FindOrFail(model, id); err != nil {
		return httpx.Error(ctx, 404, notFoundMessage(spec))
	}
	if spec.OnShow != nil {
		spec.OnShow(deref(model))
	}
	return ctx.Response().Success().Json(spec.renderRow(deref(model)))
}

/* ---------------- PUT /api/admin/{resource}/{id} ---------------- */

func (r *GatewayController) Update(ctx http.Context) http.Response {
	spec, ok := currentSpec(ctx)
	if !ok {
		return httpx.Error(ctx, 404, unknownResource(ctx))
	}
	if resp := requirePermission(ctx, spec.prefix()+".edit"); resp != nil {
		return *resp
	}

	id, valid := parseID(ctx)
	if !valid {
		return httpx.Error(ctx, 404, notFoundMessage(spec))
	}

	model := spec.newInstance()
	if err := facadesQuery().FindOrFail(model, id); err != nil {
		return httpx.Error(ctx, 404, notFoundMessage(spec))
	}

	body := map[string]any{}
	if err := ctx.Request().Bind(&body); err != nil {
		return httpx.Error(ctx, 422, "Invalid JSON body")
	}
	data := spec.filterPayload(body)
	if spec.BeforeSave != nil {
		if err := spec.BeforeSave(data); err != nil {
			return httpx.Error(ctx, 422, err.Error())
		}
	}
	if err := applyData(model, data); err != nil {
		return httpx.Error(ctx, 422, err.Error())
	}
	if err := facadesQuery().Save(model); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	if resp := runAfterSave(ctx, spec, model, data); resp != nil {
		return *resp
	}
	auditservices.Record(ctx, CurrentUserID(ctx), "update", spec.Name, id, http.StatusOK, data)

	return ctx.Response().Success().Json(spec.renderRow(deref(model)))
}

// runAfterSave executes the spec's post-persist hook, converting failures
// into a uniform error response.
func runAfterSave(ctx http.Context, spec *Spec, model any, data map[string]any) *http.Response {
	if spec.AfterSave == nil {
		return nil
	}
	if err := spec.AfterSave(deref(model), data); err != nil {
		resp := httpx.Error(ctx, 500, err.Error())
		return &resp
	}
	return nil
}

/* ---------------- DELETE /api/admin/{resource}/{id} ---------------- */

func (r *GatewayController) Destroy(ctx http.Context) http.Response {
	spec, ok := currentSpec(ctx)
	if !ok {
		return httpx.Error(ctx, 404, unknownResource(ctx))
	}
	if resp := requirePermission(ctx, spec.prefix()+".delete"); resp != nil {
		return *resp
	}

	id, valid := parseID(ctx)
	if !valid {
		return httpx.Error(ctx, 404, notFoundMessage(spec))
	}

	model := spec.newInstance()
	setID(model, id)
	if _, err := facadesQuery().Delete(model); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	auditservices.Record(ctx, CurrentUserID(ctx), "destroy", spec.Name, id, http.StatusOK, nil)

	return ctx.Response().Success().Json(http.Json{"deleted": true})
}

/* ---------------- POST /api/admin/{resource}/bulk-delete ---------------- */

func (r *GatewayController) BulkDelete(ctx http.Context) http.Response {
	return r.doBulkDelete(ctx, ctx.Request().Route("resource"))
}

// BulkDeleteForResource serves static re-binds (see RegisterBulkDeleteFor)
// where no {resource} wildcard segment exists.
func (r *GatewayController) BulkDeleteForResource(ctx http.Context, resource string) http.Response {
	return r.doBulkDelete(ctx, resource)
}

func (r *GatewayController) doBulkDelete(ctx http.Context, resource string) http.Response {
	spec, ok := Get(resource)
	if !ok {
		return httpx.Error(ctx, 404, "Unknown resource \""+resource+"\"")
	}
	if resp := requirePermission(ctx, spec.prefix()+".delete"); resp != nil {
		return *resp
	}

	var req struct {
		IDs []uint64 `json:"ids"`
	}
	if err := ctx.Request().Bind(&req); err != nil || len(req.IDs) == 0 {
		return httpx.Error(ctx, 422, "ids[] is required")
	}

	values := make([]any, len(req.IDs))
	for i, id := range req.IDs {
		values[i] = id
	}

	result, err := facadesQuery().Model(spec.newInstance()).
		WhereIn("id", values).
		Delete()
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	removed := int64(0)
	if result != nil {
		removed = result.RowsAffected
	}
	auditservices.Record(ctx, CurrentUserID(ctx), "bulk-delete", resource, 0, http.StatusOK, map[string]any{"ids": req.IDs})
	return ctx.Response().Success().Json(http.Json{"removed": removed})
}

/* ---------------- POST /api/admin/{resource}/{id}/restore ---------------- */

// Restore undoes a soft delete for specs that declare SoftDeletable.
func (r *GatewayController) Restore(ctx http.Context) http.Response {
	spec, ok := currentSpec(ctx)
	if !ok {
		return httpx.Error(ctx, 404, unknownResource(ctx))
	}
	if !spec.SoftDeletable {
		return httpx.Error(ctx, 404, "resource has no recycle bin")
	}
	if resp := requirePermission(ctx, spec.prefix()+".delete"); resp != nil {
		return *resp
	}

	id, valid := parseID(ctx)
	if !valid {
		return httpx.Error(ctx, 404, notFoundMessage(spec))
	}

	if _, err := facadesQuery().Model(spec.newInstance()).
		Where("id = ?", id).
		Update("deleted_at", nil); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	auditservices.Record(ctx, CurrentUserID(ctx), "restore", spec.Name, id, http.StatusOK, nil)

	model := spec.newInstance()
	if err := facadesQuery().FindOrFail(model, id); err != nil {
		return httpx.Error(ctx, 404, notFoundMessage(spec))
	}
	return ctx.Response().Success().Json(spec.renderRow(deref(model)))
}

/* ---------------- helpers ---------------- */

func currentSpec(ctx http.Context) (*Spec, bool) {
	return Get(ctx.Request().Route("resource"))
}

func requirePermission(ctx http.Context, required string) *http.Response {
	identity, ok := authcontrollers.RequireIdentity(ctx)
	if !ok {
		resp := httpx.Error(ctx, 401, "Unauthorized")
		return &resp
	}
	if !accessmatcher.Can(identity.Permissions, required) {
		resp := httpx.Error(ctx, 403, "Missing permission: "+required)
		return &resp
	}
	return nil
}

func parseID(ctx http.Context) (uint64, bool) {
	raw := ctx.Request().Route("id")
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, false
	}
	return id, true
}

func unknownResource(ctx http.Context) string {
	return "Unknown resource \"" + ctx.Request().Route("resource") + "\""
}

func notFoundMessage(spec *Spec) string { return spec.Name + " not found" }

// applyData overlays validated payload keys onto a model instance via a JSON
// round-trip, so only Fillable fields can ever reach persistence.
func applyData(model any, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, model)
}

func renderSlice(spec *Spec, rows any) []map[string]any {
	rv := reflect.ValueOf(rows)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	out := make([]map[string]any, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		out = append(out, spec.renderRow(deref(rv.Index(i).Interface())))
	}
	return out
}

func deref(v any) any {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		return rv.Elem().Interface()
	}
	return v
}

func setID(model any, id uint64) {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if f := v.FieldByName("ID"); f.IsValid() && f.CanSet() && f.Kind() == reflect.Uint64 {
		f.SetUint(id)
	}
}

// getID reads the primary key off a persisted model instance.
func getID(model any) uint64 {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if f := v.FieldByName("ID"); f.IsValid() && f.Kind() == reflect.Uint64 {
		return f.Uint()
	}
	return 0
}

func contains(list []string, want string) bool {
	for _, item := range list {
		if item == want {
			return true
		}
	}
	return false
}

func clampInt(raw string, minV, maxV int) int {
	n, err := strconv.Atoi(raw)
	if err != nil || n < minV {
		return minV
	}
	if n > maxV {
		return maxV
	}
	return n
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
