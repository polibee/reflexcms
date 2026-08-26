package models

import (
	"time"
)

// AdminSession backs the opaque admin session token. Only the SHA-256 hash of
// the token is stored; the raw token exists solely in the login response and
// in the httpOnly cookie held by the admin BFF.
type AdminSession struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	UserID     uint64    `gorm:"index" json:"user_id"`
	TokenHash  string    `gorm:"uniqueIndex;size:64" json:"-"`
	ExpiresAt  time.Time `json:"expires_at"`
	IP         string    `gorm:"size:64" json:"ip"`
	UserAgent  string    `gorm:"size:255" json:"user_agent"`
	LastUsedAt time.Time `json:"last_used_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AdminSession) TableName() string { return "admin_sessions" }
