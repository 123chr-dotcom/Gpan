package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex"`
	IsActive     bool   `gorm:"default:true"`
	LastLogin    *time.Time
	Files        []File      `gorm:"foreignKey:UserID"`
	Directories  []Directory `gorm:"foreignKey:UserID"`
}
