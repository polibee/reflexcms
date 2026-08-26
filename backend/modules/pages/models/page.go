package models

import (
	"github.com/goravel/framework/database/orm"
)

// Page is an admin-authored static page (隐私政策、服务条款、关于我们…).
// Content is rich-text HTML authored in the admin editor and rendered
// as-is on the public site — admins are trusted authors.
type Page struct {
	orm.Timestamps
	ID      uint64 `gorm:"primaryKey" json:"id"`
	Title   string `gorm:"size:255" json:"title"`
	Slug    string `gorm:"uniqueIndex;size:191" json:"slug"`
	Content string `gorm:"type:text" json:"content"`
	Status  string `gorm:"size:32;default:draft;index" json:"status"`
}

func (Page) TableName() string { return "pages" }
