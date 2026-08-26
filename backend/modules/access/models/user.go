package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type User struct {
	orm.Timestamps
	orm.SoftDeletes
	ID       uint64 `gorm:"primaryKey" json:"id"`
	Username string `gorm:"uniqueIndex;size:64" json:"username"`
	Email    string `gorm:"uniqueIndex;size:255" json:"email"`
	// password carries an omitempty json key (not "-") because the gateway
	// hydrates models through a JSON round-trip; output is stripped by the
	// users spec Hidden list instead.
	Password string           `gorm:"size:255" json:"password,omitempty"`
	Avatar   string           `gorm:"size:512" json:"avatar"`
	Bio      string           `gorm:"type:text" json:"bio"`
	Website  string           `gorm:"size:255" json:"website"`
	RoleID   uint64           `gorm:"index" json:"role_id"`
	Status   string           `gorm:"size:32;default:active;index" json:"status"`
	BannedAt *carbon.DateTime `gorm:"" json:"banned_at,omitempty"`
	Role     *Role            `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (User) TableName() string { return "users" }

// IsAdminActive reports whether the account may sign in to the admin panel.
func (u *User) IsAdminActive() bool { return u.Status == "active" }

// AdminSessionsTTL is how long an issued admin session stays valid.
const AdminSessionsTTL = 8 * time.Hour
