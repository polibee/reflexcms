package adminhub

import (
	"path"
	"strings"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	auditservices "reflexcms/backend/modules/audit/services"
)

var allowedUploadExt = map[string]string{
	".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png",
	".webp": "image/webp", ".gif": "image/gif",
}

const maxUploadBytes = 5 << 20 // 5 MiB

// Upload handles POST /api/admin/uploads (multipart, field "file").
// Files are stored under the local public disk with random names — original
// filenames are never reused, so no traversal or overwrite is possible.
// They are served from /storage/media/* via the static route.
func Upload(ctx http.Context) http.Response {
	if resp := CheckPermission(ctx, "articles.edit"); resp != nil {
		return *resp
	}

	file, err := ctx.Request().File("file")
	if err != nil {
		return httpx.Error(ctx, 422, "file is required")
	}

	size, err := file.Size()
	if err != nil {
		return httpx.Error(ctx, 422, "cannot read file size")
	}
	if size > maxUploadBytes {
		return httpx.Error(ctx, 422, "file exceeds 5MB limit")
	}

	ext := strings.ToLower(path.Ext(file.GetClientOriginalName()))
	if _, ok := allowedUploadExt[ext]; !ok {
		return httpx.Error(ctx, 422, "unsupported file type: "+ext)
	}

	name := randHex(12) + ext
	saved, err := facades.Storage().PutFileAs("media", file, name)
	if err != nil {
		return httpx.Error(ctx, 500, "store failed: "+err.Error())
	}

	auditservices.Record(ctx, CurrentUserID(ctx), "upload", "media", 0, http.StatusOK,
		map[string]any{"path": saved, "size": size})

	return ctx.Response().Success().Json(http.Json{
		"url":  "/storage/" + saved,
		"path": saved,
	})
}
