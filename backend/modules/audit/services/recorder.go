// Package audit records every admin mutation into admin_operation_logs.
// The writer is fire-and-forget: audit failures must never break the
// request, only log.
package services

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	auditmodels "reflexcms/backend/modules/audit/models"
)

var sensitiveKeys = map[string]bool{
	"password": true, "password_hash": true, "token": true,
	"secret": true, "authorization": true,
}

// Record persists one audit entry asynchronously. Callers pass the already
// parsed request body because the transport layer may not allow re-reads.
func Record(ctx http.Context, userID uint64, action, resource string, resourceID uint64, statusCode int, body map[string]any) {
	entry := auditmodels.AdminOperationLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: strconv.FormatUint(resourceID, 10),
		StatusCode: statusCode,
		IP:         truncate(ctx.Request().Ip(), 64),
		UserAgent:  truncate(ctx.Request().Header("User-Agent", ""), 255),
		Payload:    sanitize(body),
		CreatedAt:  time.Now(),
	}

	go func() {
		if err := facades.Orm().Query().Create(&entry); err != nil {
			facades.Log().Warning("audit: write failed: " + err.Error())
		}
	}()
}

func sanitize(body map[string]any) string {
	if len(body) == 0 {
		return ""
	}
	clean := make(map[string]any, len(body))
	for k, v := range body {
		if sensitiveKeys[strings.ToLower(k)] {
			clean[k] = "[redacted]"
			continue
		}
		if s, ok := v.(string); ok && len(s) > 500 {
			v = truncate(s, 500)
		}
		clean[k] = v
	}
	b, err := json.Marshal(clean)
	if err != nil {
		return ""
	}
	return truncate(string(b), 2000)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}
