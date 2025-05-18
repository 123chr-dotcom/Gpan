package models

import (
	"gorm.io/gorm"
	"fmt"
	"gpan/database"
)

type File struct {
	gorm.Model
	Name        string `json:"name"`      // 原始文件名
	Path        string `json:"path"`      // 存储路径
	Size        int64  `json:"size"`      // 文件大小(字节)
	Extension   string `json:"extension"` // 文件扩展名
	DirectoryID uint   `json:"directory_id" gorm:"index"` // 所属目录ID
	UserID      uint   `json:"user_id" gorm:"not null;index"` // 所属用户ID
}

// CreateFile 创建文件记录
func CreateFile(file *File) error {
	if database.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return database.DB.Create(file).Error
}

// GetFileByID 根据ID获取文件
func GetFileByID(id uint) (*File, error) {
	if database.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	var file File
	err := database.DB.First(&file, id).Error
	return &file, err
}

// GetAllFiles 获取所有文件
func GetAllFiles() ([]File, error) {
	if database.DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	var files []File
	err := database.DB.Order("id desc").Find(&files).Error
	return files, err
}

// DeleteFile 删除文件记录
func DeleteFile(id uint) error {
	if database.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return database.DB.Delete(&File{}, id).Error
}
