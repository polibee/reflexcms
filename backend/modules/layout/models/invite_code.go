package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type InviteCode struct {
	orm.Timestamps
	ID        uint64     `gorm:"primaryKey" json:"id"`
	Code      string     `gorm:"uniqueIndex;size:32" json:"code"`
	CreatorID uint64     `gorm:"index" json:"creator_id"`
	MaxUses   int        `gorm:"default:1" json:"max_uses"`
	UsedCount int        `gorm:"default:0" json:"used_count"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (InviteCode) TableName() string { return "invite_codes" }
