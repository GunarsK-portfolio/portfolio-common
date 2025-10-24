package models

import "time"

type StorageFile struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	S3Key     string    `json:"-" gorm:"column:s3_key"`
	S3Bucket  string    `json:"-" gorm:"column:s3_bucket"`
	FileName  string    `json:"fileName" gorm:"column:file_name"`
	FileSize  int64     `json:"fileSize" gorm:"column:file_size"`
	MimeType  string    `json:"mimeType" gorm:"column:mime_type"`
	FileType  string    `json:"fileType" gorm:"column:file_type"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (StorageFile) TableName() string {
	return "storage.files"
}
