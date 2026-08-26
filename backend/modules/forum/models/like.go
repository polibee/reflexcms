package models

import (
	"time"
)

const (
	LikeableTopic = "topics"
	LikeableReply = "replies"
)

// Like is the polymorphic like record. The composite unique index on
// (user_id, likeable_type, likeable_id) makes toggles idempotent.
type Like struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	UserID       uint64    `gorm:"index" json:"user_id"`
	LikeableType string    `gorm:"size:64" json:"likeable_type"`
	LikeableID   uint64    `json:"likeable_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (Like) TableName() string { return "likes" }
