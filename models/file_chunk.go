package models

import (
	"time"
	"gorm.io/gorm"
)

type FileChunk struct {
	gorm.Model
	ChunkNumber int    `gorm:"not null"`
	Size        int64  `gorm:"not null"`
	Hash        string `gorm:"not null"`
	Path        string `gorm:"not null"`
	FileID      uint   `gorm:"not null;index"`
	Status      ChunkStatus `gorm:"default:1"` // 1-上传中 2-上传完成 3-校验失败
	UploadedAt  time.Time
}

type ChunkStatus int

const (
	ChunkStatusUploading ChunkStatus = 1
	ChunkStatusCompleted ChunkStatus = 2
	ChunkStatusFailed    ChunkStatus = 3
)
