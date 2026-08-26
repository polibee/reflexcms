package models

import (
	"github.com/goravel/framework/database/orm"
)

type Carousel struct {
	orm.Timestamps
	ID       uint64 `gorm:"primaryKey" json:"id"`
	Title    string `gorm:"size:255" json:"title"`
	ImageURL string `gorm:"column:image_url;size:1024" json:"image_url"`
	LinkURL  string `gorm:"column:link_url;size:1024" json:"link_url"`
	Sort     int    `gorm:"default:100" json:"sort"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

func (Carousel) TableName() string { return "carousels" }
