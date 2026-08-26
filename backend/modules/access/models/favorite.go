package models

import "github.com/goravel/framework/support/carbon"

// Favorite marks a user's bookmark of a piece of content. The target is
// type-discriminated ("topic" today, "article" later) so one table serves
// every domain without cross-module foreign keys.
type Favorite struct {
	ID      uint64 `gorm:"primaryKey" json:"id"`
	UserID  uint64 `gorm:"index" json:"user_id"`
	FavType string `gorm:"size:16;default:topic" json:"fav_type"`
	FavID   uint64 `gorm:"index" json:"fav_id"`
	// CreatedAt is hand-managed (no UpdatedAt): favorites are immutable
	// bookmarks, toggled by insert/delete.
	CreatedAt *carbon.DateTime `gorm:"" json:"created_at,omitempty"`
}

func (Favorite) TableName() string { return "favorites" }
