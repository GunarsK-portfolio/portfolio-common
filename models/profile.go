package models

import "time"

type Profile struct {
	ID           int64        `json:"id" gorm:"primaryKey"`
	FullName     string       `json:"name" gorm:"column:full_name" binding:"required"`
	Title        string       `json:"title"`
	Bio          string       `json:"tagline"`
	Email        string       `json:"email" binding:"omitempty,email"`
	Phone        string       `json:"phone,omitempty"`
	Location     string       `json:"location,omitempty"`
	Github       string       `json:"github,omitempty" binding:"omitempty,url"`
	Linkedin     string       `json:"linkedin,omitempty" binding:"omitempty,url"`
	AvatarFileID *int64       `json:"avatarFileId,omitempty" gorm:"column:avatar_file_id"`
	AvatarFile   *StorageFile `json:"avatarFile,omitempty" gorm:"foreignKey:AvatarFileID"`
	ResumeFileID *int64       `json:"resumeFileId,omitempty" gorm:"column:resume_file_id"`
	ResumeFile   *StorageFile `json:"resumeFile,omitempty" gorm:"foreignKey:ResumeFileID"`
	CreatedAt    time.Time    `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time    `json:"updatedAt" gorm:"column:updated_at"`
}

func (Profile) TableName() string {
	return "portfolio.profile"
}
