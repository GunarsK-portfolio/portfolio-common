package models

import "time"

type Profile struct {
	ID           int64        `json:"id" gorm:"primaryKey"`
	FullName     string       `json:"name" gorm:"column:full_name" binding:"required"`
	Title        string       `json:"title"`
	Bio          string       `json:"tagline"`
	Email        string       `json:"email"`
	Phone        string       `json:"phone,omitempty"`
	Location     string       `json:"location,omitempty"`
	AvatarFileID *int64       `json:"avatarFileId,omitempty" gorm:"column:avatar_file_id"`
	AvatarFile   *StorageFile `json:"-" gorm:"foreignKey:AvatarFileID"`
	AvatarURL    string       `json:"avatarUrl,omitempty" gorm:"-"`
	ResumeFileID *int64       `json:"resumeFileId,omitempty" gorm:"column:resume_file_id"`
	ResumeFile   *StorageFile `json:"-" gorm:"foreignKey:ResumeFileID"`
	ResumeURL    string       `json:"resumeUrl,omitempty" gorm:"-"`
	CreatedAt    time.Time    `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time    `json:"updatedAt" gorm:"column:updated_at"`
}

func (Profile) TableName() string {
	return "portfolio.profile"
}
