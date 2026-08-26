package models

import (
	"github.com/goravel/framework/database/orm"
)

type Widget struct {
	orm.Timestamps
	ID        uint64 `gorm:"primaryKey" json:"id"`
	Area      string `gorm:"size:32;index" json:"area"`
	Type      string `gorm:"column:widget_type;size:64" json:"widget_type"`
	Title     string `gorm:"size:255" json:"title"`
	Config    string `gorm:"type:text" json:"config"`
	Sort      int    `gorm:"default:100" json:"sort"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
}

func (Widget) TableName() string { return "widgets" }
