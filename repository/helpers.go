package repository

import (
	"context"

	"gorm.io/gorm"
)

// CheckRowsAffected returns gorm.ErrRecordNotFound if no rows were affected.
// Use for delete/update operations that should fail if target doesn't exist.
func CheckRowsAffected(result *gorm.DB) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// SafeUpdater provides safe update operations for GORM repositories.
// Checks existence before updating to ensure idempotent updates (no false 404s).
type SafeUpdater struct {
	db *gorm.DB
}

// NewSafeUpdater creates a new SafeUpdater instance.
func NewSafeUpdater(db *gorm.DB) *SafeUpdater {
	return &SafeUpdater{db: db}
}

// Update performs an update excluding system fields (ID, CreatedAt, UpdatedAt).
// Uses Updates to avoid zero-value overwrites unlike Save.
// Checks existence first to ensure idempotent updates (no false 404s).
func (s *SafeUpdater) Update(ctx context.Context, model interface{}, id int64) error {
	return s.UpdateWithOptions(ctx, model, id, nil)
}

// UpdateWithAssociations performs safe update with association handling.
// Use for models with has-many or many-to-many relationships.
func (s *SafeUpdater) UpdateWithAssociations(ctx context.Context, model interface{}, id int64) error {
	return s.UpdateWithOptions(ctx, model, id, &gorm.Session{FullSaveAssociations: true})
}

// UpdateWithOptions is the internal implementation with optional session config.
func (s *SafeUpdater) UpdateWithOptions(ctx context.Context, model interface{}, id int64, session *gorm.Session) error {
	// Check existence using COUNT to avoid loading data
	var count int64
	if err := s.db.WithContext(ctx).Model(model).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	// Build update query
	db := s.db.WithContext(ctx).Model(model).Where("id = ?", id)

	// Apply session options if provided
	if session != nil {
		db = db.Session(session)
	}

	// Now perform update - RowsAffected=0 is OK (idempotent)
	return db.Omit("ID", "CreatedAt", "UpdatedAt").Updates(model).Error
}
