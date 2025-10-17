package repository

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ActionLog represents an entry in the audit.action_log table
type ActionLog struct {
	ID           int64           `json:"id" gorm:"primaryKey"`
	ActionType   string          `json:"action_type" gorm:"column:action_type"`
	ResourceType *string         `json:"resource_type,omitempty" gorm:"column:resource_type"`
	ResourceID   *int64          `json:"resource_id,omitempty" gorm:"column:resource_id"`
	UserID       *int64          `json:"user_id,omitempty" gorm:"column:user_id"`
	IPAddress    *string         `json:"ip_address,omitempty" gorm:"column:ip_address"`
	UserAgent    *string         `json:"user_agent,omitempty" gorm:"column:user_agent"`
	Metadata     json.RawMessage `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb"`
	CreatedAt    time.Time       `json:"created_at" gorm:"column:created_at"`
}

func (ActionLog) TableName() string {
	return "audit.action_log"
}

// ActionLogRepository handles action log database operations
type ActionLogRepository interface {
	LogAction(log *ActionLog) error
	GetActionsByType(actionType string, limit int) ([]ActionLog, error)
	GetActionsByResource(resourceType string, resourceID int64) ([]ActionLog, error)
	GetActionsByUser(userID int64, limit int) ([]ActionLog, error)
	CountActionsByResource(resourceType string, resourceID int64) (int64, error)
}

type actionLogRepository struct {
	db *gorm.DB
}

// NewActionLogRepository creates a new action log repository
func NewActionLogRepository(db *gorm.DB) ActionLogRepository {
	return &actionLogRepository{db: db}
}

// LogAction inserts a new action log entry
func (r *actionLogRepository) LogAction(log *ActionLog) error {
	return r.db.Create(log).Error
}

// GetActionsByType retrieves actions by type
func (r *actionLogRepository) GetActionsByType(actionType string, limit int) ([]ActionLog, error) {
	var logs []ActionLog
	err := r.db.Where("action_type = ?", actionType).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetActionsByResource retrieves actions for a specific resource
func (r *actionLogRepository) GetActionsByResource(resourceType string, resourceID int64) ([]ActionLog, error) {
	var logs []ActionLog
	err := r.db.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetActionsByUser retrieves actions by a specific user
func (r *actionLogRepository) GetActionsByUser(userID int64, limit int) ([]ActionLog, error) {
	var logs []ActionLog
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// CountActionsByResource counts actions for a specific resource
func (r *actionLogRepository) CountActionsByResource(resourceType string, resourceID int64) (int64, error) {
	var count int64
	err := r.db.Model(&ActionLog{}).
		Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Count(&count).Error
	return count, err
}
