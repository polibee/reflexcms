package models

import (
	"github.com/goravel/framework/database/orm"
)

const (
	TopicOpen   = "open"
	TopicClosed = "closed"
)

type ForumCategory struct {
	orm.Timestamps
	ID          uint64 `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:128" json:"name"`
	Slug        string `gorm:"uniqueIndex;size:191" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Icon        string `gorm:"size:64" json:"icon"`
	Sort        int    `gorm:"default:100" json:"sort"`
	TopicCount  uint64 `gorm:"default:0" json:"topic_count"`
}

func (ForumCategory) TableName() string { return "forum_categories" }
