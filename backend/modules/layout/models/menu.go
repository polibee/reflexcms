package models

import (
	"github.com/goravel/framework/database/orm"
)

type Menu struct {
	orm.Timestamps
	ID        uint64 `gorm:"primaryKey" json:"id"`
	Location  string `gorm:"size:32;index" json:"location"`
	Label     string `gorm:"size:128" json:"label"`
	URL       string `gorm:"column:url;size:512" json:"url"`
	Sort      int    `gorm:"default:100" json:"sort"`
	ParentID  uint64 `gorm:"default:0;index" json:"parent_id"`
	IsVisible bool   `gorm:"default:true" json:"is_visible"`
}

func (Menu) TableName() string { return "menus" }
