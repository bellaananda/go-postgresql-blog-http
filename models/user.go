package models

import (
	"gorm.io/gorm"
)

type GormUser struct {
	gorm.Model
	Name     string         `gorm:"size:255"`
	Email    string         `gorm:"unique;size:255;not null"`
	Password string         `gorm:"size:255"`
	Username string         `gorm:"unique;size:255;not null"`
	Posts    []*GormPost    `gorm:"foreignkey:UserID"`
	Comments []*GormComment `gorm:"foreignkey:UserID"`
}
