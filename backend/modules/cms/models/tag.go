package models

import (
	"github.com/goravel/framework/database/orm"
)

type Tag struct {
	orm.Timestamps
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex;size:128" json:"name"`
	Slug string `gorm:"uniqueIndex;size:191" json:"slug"`
}

func (Tag) TableName() string { return "tags" }

// ArticleTag is the articles↔tags join pivot.
type ArticleTag struct {
	ArticleID uint64 `gorm:"primaryKey;autoIncrement:false" json:"article_id"`
	TagID     uint64 `gorm:"primaryKey;autoIncrement:false" json:"tag_id"`
}

func (ArticleTag) TableName() string { return "article_tags" }
