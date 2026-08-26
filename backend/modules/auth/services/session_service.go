package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"reflexcms/backend/app/facades"
	accessmodels "reflexcms/backend/modules/access/models"
	accessservices "reflexcms/backend/modules/access/services"
	authmodels "reflexcms/backend/modules/auth/models"
)

const (
	sessionTTL       = accessmodels.AdminSessionsTTL
	loginWindow      = time.Minute
	loginMaxAttempts = 5
)

var (
	ErrRateLimited        = errors.New("too many login attempts")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionExpired     = errors.New("session expired")
)

// Identity is the wire shape of AuthUser in shared/types/api.ts.
type Identity struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

// Login validates credentials (with per-IP rate limiting AND per-account
// lockout) and issues an opaque session token. Only the SHA-256 hash of the
// token is persisted.
func Login(email, password, ip, userAgent string) (string, *Identity, error) {
	if err := checkLoginRate(ip); err != nil {
		return "", nil, err
	}
	if err := checkAccountLock(email); err != nil {
		return "", nil, err
	}

	user, err := accessservices.FindActiveUserByEmail(email)
	if err != nil {
		recordAccountFailure(email)
		return "", nil, ErrInvalidCredentials
	}

	ok := facades.Hash().Check(password, user.Password)
	if !ok {
		recordAccountFailure(email)
		return "", nil, ErrInvalidCredentials
	}

	facades.Cache().Forget(loginKey(ip))
	facades.Cache().Forget(accountFailKey(email))

	token, err := IssueSession(user.ID, ip, userAgent)
	if err != nil {
		return "", nil, err
	}

	return token, IdentityOf(user), nil
}

// Resolve maps a raw bearer token back to a live identity, enforcing expiry
// and account status. Sessions are revocable by simply deleting their row.
func Resolve(token string) (*Identity, error) {
	if token == "" {
		return nil, ErrSessionExpired
	}

	var session authmodels.AdminSession
	// FirstOrFail: goravel's First silently zero-fills on empty results.
	if err := facades.Orm().Query().
		Where("token_hash = ?", HashToken(token)).
		FirstOrFail(&session); err != nil {
		return nil, ErrSessionExpired
	}

	now := time.Now()
	if now.After(session.ExpiresAt) {
		_, _ = facades.Orm().Query().Where("id = ?", session.ID).Delete(&authmodels.AdminSession{})
		return nil, ErrSessionExpired
	}

	var user accessmodels.User
	if err := facades.Orm().Query().
		Where("id = ? AND status = ?", session.UserID, "active").
		FirstOrFail(&user); err != nil {
		return nil, ErrSessionExpired
	}

	if user.RoleID > 0 {
		var role accessmodels.Role
		if err := facades.Orm().Query().Where("id = ?", user.RoleID).FirstOrFail(&role); err == nil {
			user.Role = &role
		}
	}

	if now.Sub(session.LastUsedAt) > time.Hour {
		_, _ = facades.Orm().Query().
			Model(&authmodels.AdminSession{}).
			Update("last_used_at", now)
	}

	return IdentityOf(&user), nil
}

func Revoke(token string) error {
	_, err := facades.Orm().Query().
		Where("token_hash = ?", HashToken(token)).
		Delete(&authmodels.AdminSession{})
	return err
}

func IssueSession(userID uint64, ip, userAgent string) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	token := hex.EncodeToString(raw)

	now := time.Now()
	session := authmodels.AdminSession{
		UserID:     userID,
		TokenHash:  HashToken(token),
		ExpiresAt:  now.Add(sessionTTL),
		IP:         ip,
		UserAgent:  userAgent,
		LastUsedAt: now,
		CreatedAt:  now,
	}
	if err := facades.Orm().Query().Create(&session); err != nil {
		return "", err
	}

	return token, nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func IdentityOf(u *accessmodels.User) *Identity {
	identity := &Identity{
		ID:          u.ID,
		Name:        u.Username,
		Email:       u.Email,
		Role:        "",
		Permissions: []string{},
	}
	if u.Role != nil {
		identity.Role = u.Role.Name
		identity.Permissions = []string(u.Role.Permissions)
	}
	return identity
}

/* ---------------- login rate limiting ---------------- */

const (
	accountFailLimit  = 10
	accountLockWindow = 15 * time.Minute
)

var ErrAccountLocked = errors.New("account temporarily locked, try again later")

// checkAccountLock enforces account-level lockout: after accountFailLimit
// consecutive failures the email is locked for accountLockWindow regardless
// of IP, defeating distributed password spraying against one account.
func checkAccountLock(email string) error {
	if parseCount(facades.Cache().Get(accountFailKey(email))) >= accountFailLimit {
		return ErrAccountLocked
	}
	return nil
}

func recordAccountFailure(email string) {
	key := accountFailKey(email)
	n := parseCount(facades.Cache().Get(key))
	_ = facades.Cache().Put(key, n+1, accountLockWindow)
}

func accountFailKey(email string) string { return "login:account:" + email }

func checkLoginRate(ip string) error {
	key := loginKey(ip)
	attempts := parseCount(facades.Cache().Get(key))
	if attempts >= loginMaxAttempts {
		return ErrRateLimited
	}
	return facades.Cache().Put(key, attempts+1, loginWindow)
}

func loginKey(ip string) string { return "login:" + ip }

func parseCount(raw any) int64 {
	s := fmt.Sprint(raw)
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return n
}
