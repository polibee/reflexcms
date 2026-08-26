package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type Notification struct {
	orm.Timestamps
	ID     uint64     `gorm:"primaryKey" json:"id"`
	UserID uint64     `gorm:"index" json:"user_id"`
	Type   string     `gorm:"size:64" json:"type"`
	Data   string     `gorm:"type:text" json:"data"`
	ReadAt *time.Time `json:"read_at"`
}

func (Notification) TableName() string { return "notifications" }
