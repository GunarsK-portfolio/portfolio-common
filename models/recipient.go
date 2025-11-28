package models

import "time"

// Recipient represents an email recipient for contact form notifications
type Recipient struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"column:email;uniqueIndex" binding:"required,email,max=255"`
	Name      string    `json:"name" gorm:"column:name" binding:"required,max=255"`
	IsActive  bool      `json:"isActive" gorm:"column:is_active;default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Recipient) TableName() string {
	return "messaging.recipients"
}

// RecipientCreate is the DTO for creating a new recipient
type RecipientCreate struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Name     string `json:"name" binding:"required,max=255"`
	IsActive *bool  `json:"isActive"` // Optional, defaults to true
}

// RecipientUpdate is the DTO for updating a recipient
type RecipientUpdate struct {
	Email    *string `json:"email,omitempty" binding:"omitempty,email,max=255"`
	Name     *string `json:"name,omitempty" binding:"omitempty,max=255"`
	IsActive *bool   `json:"isActive,omitempty"`
}
