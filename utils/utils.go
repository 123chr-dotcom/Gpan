package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// GenerateUUID 生成UUID字符串
func GenerateUUID() string {
	return uuid.New().String()
}

// DeleteFile 删除文件
func DeleteFile(path string) error {
	return os.Remove(path)
}

// EnsureDir 确保目录存在
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// GetFileExtension 获取文件扩展名
func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.ToLower(ext)
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
