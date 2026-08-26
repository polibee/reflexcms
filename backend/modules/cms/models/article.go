package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"

	strlist "reflexcms/backend/app/support/strlist"
)

const (
	ArticleDraft     = "draft"
	ArticleScheduled = "scheduled"
	ArticlePublished = "published"
	ArticleArchived  = "archived"
)

type Article struct {
	orm.Timestamps
	orm.SoftDeletes
	ID           uint64           `gorm:"primaryKey" json:"id"`
	AuthorID     uint64           `gorm:"index" json:"author_id"`
	CategoryID   uint64           `gorm:"index" json:"category_id"`
	Title        string           `gorm:"size:512" json:"title"`
	Slug         string           `gorm:"uniqueIndex;size:191" json:"slug"`
	Content      string           `gorm:"type:text" json:"content"`
	Summary      string           `gorm:"type:text" json:"summary"`
	Cover        string           `gorm:"size:512" json:"cover"`
	Status       string           `gorm:"size:32;default:draft;index" json:"status"`
	IsPinned     bool             `gorm:"default:false;index" json:"is_pinned"`
	PublishedAt  *carbon.DateTime `json:"published_at"`
	ViewCount    uint64           `gorm:"default:0" json:"view_count"`
	CommentCount uint64           `gorm:"default:0" json:"comment_count"`
	AISummary    string           `gorm:"type:text" json:"ai_summary"`
	AIKeywords   strlist.List     `gorm:"type:jsonb" json:"ai_keywords"`
}

func (Article) TableName() string { return "articles" }

// IsPublished reports whether the article is publicly visible.
func (a *Article) IsPublished() bool { return a.Status == ArticlePublished }
