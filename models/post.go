package models

import (
	"time"

	"gorm.io/gorm"
)

type GormPost struct {
	gorm.Model
	UserID      uint   `gorm:"index;not null"`
	Title       string `gorm:"size:255"`
	Content     string `gorm:"type:text"`
	Thumbnail   string `gorm:"type:text"`
	IsPublished bool   `gorm:"default:false"`
	PublishedAt time.Time
	User        *GormUser      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Comments    []*GormComment `gorm:"foreignkey:PostID"`
}
