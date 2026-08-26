package models

import (
	"github.com/goravel/framework/database/orm"
)

type ArticleCategory struct {
	orm.Timestamps
	ID          uint64 `gorm:"primaryKey" json:"id"`
	ParentID    uint64 `gorm:"index" json:"parent_id"`
	Name        string `gorm:"size:128" json:"name"`
	Slug        string `gorm:"uniqueIndex;size:191" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Sort        int    `gorm:"default:100" json:"sort"`
}

func (ArticleCategory) TableName() string { return "article_categories" }
