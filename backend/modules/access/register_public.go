package access

import (
	"fmt"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"

	"reflexcms/backend/app/facades"
	httpx "reflexcms/backend/app/support/httpx"
	authservices "reflexcms/backend/modules/auth/services"
	accessmodels "reflexcms/backend/modules/access/models"
	settingsservices "reflexcms/backend/modules/settings/services"
)

/* Public registration: POST /api/v1/auth/register. Gated by
 * reg.enabled / reg.invite_required settings; consumes one use of an
 * invite code when required. Returns a session token like login so the
 * new user is signed in immediately. */

// RegisterPublic handles public account creation.
func RegisterPublic(ctx http.Context) http.Response {
	if !regEnabled() {
		return httpx.Error(ctx, 403, "本站暂未开放注册")
	}

	var req struct {
		Username   string `json:"username"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		InviteCode string `json:"invite_code"`
	}
	if err := ctx.Request().Bind(&req); err != nil {
		return httpx.Error(ctx, 422, "invalid request body")
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return httpx.Error(ctx, 422, "用户名、邮箱和密码不能为空")
	}
	if len([]rune(req.Username)) > 32 {
		return httpx.Error(ctx, 422, "用户名过长（最多 32 字）")
	}
	if len(req.Password) < 8 {
		return httpx.Error(ctx, 422, "密码至少 8 位")
	}
	if !strings.Contains(req.Email, "@") {
		return httpx.Error(ctx, 422, "邮箱格式不正确")
	}
	if err := settingsservices.ValidateRegistrationEmail(req.Email); err != nil {
		return httpx.Error(ctx, 422, err.Error())
	}

	// Invite code gate
	inviteRequired := regInviteRequired()
	inviteCode := strings.TrimSpace(req.InviteCode)
	if inviteRequired && inviteCode == "" {
		return httpx.Error(ctx, 422, "本站为邀请制注册，请填写邀请码")
	}

	// Uniqueness checks
	var dupName, dupEmail int64
	dupName, _ = facades.Orm().Query().Table("users").
		Where("username = ? AND deleted_at IS NULL", req.Username).Count()
	if dupName > 0 {
		return httpx.Error(ctx, 422, "用户名已被占用")
	}
	dupEmail, _ = facades.Orm().Query().Table("users").
		Where("email = ? AND deleted_at IS NULL", req.Email).Count()
	if dupEmail > 0 {
		return httpx.Error(ctx, 422, "邮箱已被注册")
	}

	hashed, err := facades.Hash().Make(req.Password)
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	newUser := accessmodels.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
		Status:   "active",
	}
	// Production bring-up path: on an empty site the very first registered
	// account becomes the super-admin, so a fresh deployment never ships
	// seeded credentials. Later sign-ups keep the default member role.
	userCount, countErr := facades.Orm().Query().Table("users").Count()
	if countErr == nil && userCount == 0 {
		newUser.RoleID = firstUserSuperAdminRoleID()
	}
	if err := facades.Orm().Query().Create(&newUser); err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}
	newUserID := newUser.ID

	// Consume invite code (max_uses / used_count / expires_at enforced).
	// On failure the fresh account is removed so the code isn't wasted.
	if inviteCode != "" {
		if err := consumeInvite(inviteCode, newUserID); err != nil {
			_, _ = facades.Orm().Query().Exec(`DELETE FROM users WHERE id = ?`, newUserID)
			return httpx.Error(ctx, 422, err.Error())
		}
	}

	token, err := authservices.IssueSession(newUserID, ctx.Request().Ip(), ctx.Request().Header("User-Agent", ""))
	if err != nil {
		return httpx.Error(ctx, 500, err.Error())
	}

	return ctx.Response().Success().Json(http.Json{
		"token":   token,
		"user":    map[string]any{"id": newUserID, "name": req.Username, "email": req.Email},
		"message": "注册成功",
	})
}

func regEnabled() bool {
	v := settingsservices.Get("reg.enabled")
	return v == "true" || v == "1"
}

func regInviteRequired() bool {
	v := settingsservices.Get("reg.invite_required")
	return v == "true" || v == "1"
}

// firstUserSuperAdminRoleID returns the super-admin role id, creating the
// role when the deployment skipped db:seed. Returns 0 (default member)
// when even that fails so registration never breaks.
func firstUserSuperAdminRoleID() uint64 {
	var role accessmodels.Role
	if err := facades.Orm().Query().Where("name = ?", "super-admin").FirstOrFail(&role); err == nil {
		return role.ID
	}
	role = accessmodels.Role{
		Name:        "super-admin",
		DisplayName: "超级管理员",
		Permissions: accessmodels.StringList{"*"},
		Sort:        10,
	}
	if err := facades.Orm().Query().Create(&role); err != nil {
		return 0
	}
	return role.ID
}

// consumeInvite validates and increments an invite code.
func consumeInvite(code string, userID uint64) error {
	var rows []map[string]any
	if err := facades.Orm().Query().Table("invite_codes").
		Where("code = ?", code).Get(&rows); err != nil || len(rows) == 0 {
		return fmt.Errorf("邀请码不存在")
	}
	inv := rows[0]
	maxUses := 1
	if v, ok := inv["max_uses"].(int64); ok {
		maxUses = int(v)
	}
	used := 0
	if v, ok := inv["used_count"].(int64); ok {
		used = int(v)
	}
	if used >= maxUses {
		return fmt.Errorf("邀请码已被用完")
	}
	if exp, ok := inv["expires_at"].(string); ok && exp != "" && exp < time.Now().UTC().Format("2006-01-02T15:04:05") {
		return fmt.Errorf("邀请码已过期")
	}
	id := uint64(0)
	if v, ok := inv["id"].(int64); ok {
		id = uint64(v)
	}
	if id == 0 {
		return fmt.Errorf("邀请码无效")
	}
	if _, err := facades.Orm().Query().Exec(`
		UPDATE invite_codes SET used_count = used_count + 1, updated_at = NOW()
		WHERE id = ? AND used_count < max_uses
	`, id); err != nil {
		return err
	}
	_ = userID
	return nil
}