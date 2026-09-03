package models

import (
	"github.com/goravel/framework/database/orm"
)

// BoardModerator assigns a moderator to a forum board; moderators can pin /
// feature / close / best-reply within their boards via the public API.
type BoardModerator struct {
	orm.Timestamps
	ID      uint64 `gorm:"primaryKey" json:"id"`
	BoardID uint64 `gorm:"uniqueIndex:idx_board_user" json:"board_id"`
	UserID  uint64 `gorm:"uniqueIndex:idx_board_user" json:"user_id"`
}

func (BoardModerator) TableName() string { return "board_moderators" }
