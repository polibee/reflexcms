package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type Topic struct {
	orm.Timestamps
	orm.SoftDeletes
	ID              uint64           `gorm:"primaryKey" json:"id"`
	UserID          uint64           `gorm:"index" json:"user_id"`
	ForumCategoryID uint64           `gorm:"index" json:"forum_category_id"`
	Title           string           `gorm:"size:512" json:"title"`
	Content         string           `gorm:"type:text" json:"content"`
	Status          string           `gorm:"size:32;default:open;index" json:"status"`
	IsPinned        bool             `gorm:"default:false" json:"is_pinned"`
	IsFeatured      bool             `gorm:"default:false" json:"is_featured"`
	ReplyCount      uint64           `gorm:"default:0" json:"reply_count"`
	LikeCount       uint64           `gorm:"default:0" json:"like_count"`
	ViewCount       uint64           `gorm:"default:0" json:"view_count"`
	LastReplyAt     *carbon.DateTime `json:"last_reply_at"`
	LastReplyUserID uint64           `json:"last_reply_user_id"`
	BestReplyID     uint64           `json:"best_reply_id"`
}

func (Topic) TableName() string { return "topics" }
