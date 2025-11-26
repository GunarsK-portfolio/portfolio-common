package models

import "time"

type PortfolioProject struct {
	ID              int64     `json:"id" gorm:"primaryKey"`
	Title           string    `json:"title" binding:"required"`
	Category        string    `json:"category,omitempty"`
	Description     string    `json:"description"`
	LongDescription string    `json:"longDescription,omitempty" gorm:"column:long_description"`
	ImageFileID     *int64    `json:"imageFileId,omitempty" gorm:"column:image_file_id"`
	GithubURL       string    `json:"githubUrl,omitempty" gorm:"column:github_url" binding:"omitempty,url"`
	LiveURL         string    `json:"liveUrl,omitempty" gorm:"column:live_url" binding:"omitempty,url"`
	StartDate       *string   `json:"startDate,omitempty" gorm:"column:start_date"`
	EndDate         *string   `json:"endDate,omitempty" gorm:"column:end_date"`
	IsOngoing       bool      `json:"isOngoing" gorm:"column:is_ongoing"`
	TeamSize        *int      `json:"teamSize,omitempty" gorm:"column:team_size"`
	Role            string    `json:"role,omitempty"`
	Featured        bool      `json:"featured"`
	Features        []string  `json:"features,omitempty" gorm:"type:jsonb;serializer:json"`
	Challenges      []string  `json:"challenges,omitempty" gorm:"type:jsonb;serializer:json"`
	Learnings       []string  `json:"learnings,omitempty" gorm:"type:jsonb;serializer:json"`
	DisplayOrder    int       `json:"displayOrder,omitempty" gorm:"column:display_order"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"column:updated_at"`

	// Associations
	ImageFile    *StorageFile `json:"imageFile,omitempty" gorm:"foreignKey:ImageFileID"`
	Technologies []Skill      `json:"technologies,omitempty" gorm:"many2many:portfolio.project_technologies;joinForeignKey:ProjectID;joinReferences:SkillID"`
}

func (PortfolioProject) TableName() string {
	return "portfolio.portfolio_projects"
}
