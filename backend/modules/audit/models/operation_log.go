package models

import (
	"time"
)

// AdminOperationLog is the audit trail for every admin mutation. Rows are
// written by the platform (never via the gateway), exposed read-only.
type AdminOperationLog struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	UserID     uint64    `gorm:"index" json:"user_id"`
	Action     string    `gorm:"size:32" json:"action"`
	Resource   string    `gorm:"size:64" json:"resource"`
	ResourceID string    `gorm:"size:32" json:"resource_id"`
	StatusCode int       `json:"status_code"`
	IP         string    `gorm:"size:64" json:"ip"`
	UserAgent  string    `gorm:"size:255" json:"user_agent"`
	Payload    string    `gorm:"type:text" json:"payload"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

func (AdminOperationLog) TableName() string { return "admin_operation_logs" }
