package models

import "gorm.io/gorm"

type Directory struct {
	gorm.Model
	Name       string `gorm:"not null"`
	ParentID   *uint  `gorm:"index"`
	UserID     uint   `gorm:"not null;index"`
	Subdirectories []Directory `gorm:"foreignKey:ParentID"`
	Files      []File `gorm:"foreignKey:DirectoryID"`
}

func (d *Directory) IsRoot() bool {
	return d.ParentID == nil
}
