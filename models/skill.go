package models

import "time"

type SkillType struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description,omitempty"`
	DisplayOrder int       `json:"displayOrder,omitempty" gorm:"column:display_order"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (SkillType) TableName() string {
	return "portfolio.cl_skill_types"
}

type Skill struct {
	ID           int64      `json:"id" gorm:"primaryKey"`
	Skill        string     `json:"skill" binding:"required"`
	SkillTypeID  int64      `json:"skillTypeId" gorm:"column:skill_type_id" binding:"required"`
	SkillType    *SkillType `json:"skillType,omitempty" gorm:"foreignKey:SkillTypeID"`
	IsVisible    bool       `json:"isVisible" gorm:"column:is_visible"`
	DisplayOrder int        `json:"displayOrder,omitempty" gorm:"column:display_order"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"column:updated_at"`

	// Computed field
	Type string `json:"type" gorm:"-"`
}

func (Skill) TableName() string {
	return "portfolio.skills"
}
