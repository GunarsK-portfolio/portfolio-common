package models

import (
	"strings"
	"time"
)

// ContactMessage represents a contact form submission
type ContactMessage struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"column:name" binding:"required,max=255"`
	Email     string     `json:"email" gorm:"column:email" binding:"required,email,max=255"`
	Subject   string     `json:"subject" gorm:"column:subject" binding:"required,max=500"`
	Message   string     `json:"message" gorm:"column:message" binding:"required"`
	Honeypot  string     `json:"-" gorm:"column:honeypot"` // Hidden from JSON, spam detection
	Status    string     `json:"status" gorm:"column:status;default:pending"`
	Attempts  int        `json:"attempts" gorm:"column:attempts;default:0"`
	LastError *string    `json:"lastError,omitempty" gorm:"column:last_error"`
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at"`
	SentAt    *time.Time `json:"sentAt,omitempty" gorm:"column:sent_at"`
}

func (ContactMessage) TableName() string {
	return "messaging.contact_messages"
}

// ContactMessageStatus constants
const (
	MessageStatusPending = "pending"
	MessageStatusQueued  = "queued"
	MessageStatusSent    = "sent"
	MessageStatusFailed  = "failed"
)

// ContactMessageCreate is the DTO for creating a new contact message (public endpoint)
type ContactMessageCreate struct {
	Name     string `json:"name" binding:"required,max=255"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Subject  string `json:"subject" binding:"required,max=500"`
	Message  string `json:"message" binding:"required,max=10000"`
	Honeypot string `json:"website" swaggerignore:"true"`
}

// IsSpam checks if the honeypot field is filled (indicates bot submission)
func (c *ContactMessageCreate) IsSpam() bool {
	return strings.TrimSpace(c.Honeypot) != ""
}
