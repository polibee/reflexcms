package models

import (
	"github.com/goravel/framework/database/orm"
)

const (
	CommentPending  = "pending"
	CommentApproved = "approved"
	CommentRejected = "rejected"
	CommentSpam     = "spam"
)

type Comment struct {
	orm.Timestamps
	orm.SoftDeletes
	ID        uint64 `gorm:"primaryKey" json:"id"`
	ArticleID uint64 `gorm:"index" json:"article_id"`
	UserID    uint64 `gorm:"index" json:"user_id"`
	ParentID  uint64 `gorm:"index" json:"parent_id"`
	Content   string `gorm:"type:text" json:"content"`
	Status    string `gorm:"size:32;default:pending;index" json:"status"`
}

func (Comment) TableName() string { return "comments" }
