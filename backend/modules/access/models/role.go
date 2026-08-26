package models

import (
	"github.com/goravel/framework/database/orm"

	strlist "reflexcms/backend/app/support/strlist"
)

// StringList is the shared JSONB string-array type (kept as an alias so
// existing references keep compiling).
type StringList = strlist.List

type Role struct {
	orm.Timestamps
	ID          uint64     `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"uniqueIndex;size:64" json:"name"`
	DisplayName string     `gorm:"size:128" json:"display_name"`
	Permissions StringList `gorm:"type:jsonb;default:'[]'" json:"permissions"`
	Sort        int        `gorm:"default:100" json:"sort"`
}

func (Role) TableName() string { return "roles" }
