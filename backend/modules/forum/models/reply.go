package models

import (
	"github.com/goravel/framework/database/orm"
)

type Reply struct {
	orm.Timestamps
	orm.SoftDeletes
	ID        uint64 `gorm:"primaryKey" json:"id"`
	TopicID   uint64 `gorm:"index" json:"topic_id"`
	UserID    uint64 `gorm:"index" json:"user_id"`
	ParentID  uint64 `gorm:"index" json:"parent_id"`
	Content   string `gorm:"type:text" json:"content"`
	Floor     int    `json:"floor"`
	LikeCount uint64 `gorm:"default:0" json:"like_count"`
}

func (Reply) TableName() string { return "replies" }
