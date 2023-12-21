package models

import (
	"time"

	"gorm.io/gorm"
)

type GormComment struct {
	gorm.Model
	UserID      uint   `gorm:"index;not null"`
	PostID      uint   `gorm:"index;not null"`
	Content     string `gorm:"type:text"`
	PublishedAt time.Time
	User        *GormUser `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Post        *GormPost `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
