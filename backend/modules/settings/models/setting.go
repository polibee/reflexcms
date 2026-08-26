package models

import (
	"github.com/goravel/framework/database/orm"
)

type Setting struct {
	orm.Timestamps
	ID    uint64 `gorm:"primaryKey" json:"id"`
	Key   string `gorm:"uniqueIndex;size:191" json:"key"`
	Value string `gorm:"type:text" json:"value"`
	Group string `gorm:"size:64;default:general" json:"group"`
}

func (Setting) TableName() string { return "settings" }
