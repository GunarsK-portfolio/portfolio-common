package models

import "time"

type Certification struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" binding:"required"`
	Issuer        string    `json:"issuer" binding:"required"`
	IssueDate     string    `json:"issueDate" gorm:"column:issue_date" binding:"required"`
	ExpiryDate    *string   `json:"expiryDate,omitempty" gorm:"column:expiry_date"`
	CredentialID  string    `json:"credentialId,omitempty" gorm:"column:credential_id"`
	CredentialURL string    `json:"credentialUrl,omitempty" gorm:"column:credential_url"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Certification) TableName() string {
	return "portfolio.certifications"
}
