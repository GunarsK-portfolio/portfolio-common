package models

import "time"

type MiniatureTheme struct {
	ID           int64              `json:"id" gorm:"primaryKey"`
	Name         string             `json:"name" binding:"required"`
	Description  string             `json:"description"`
	CoverImageID *int64             `json:"coverImageId,omitempty" gorm:"column:cover_image_id"`
	DisplayOrder int                `json:"displayOrder,omitempty" gorm:"column:display_order"`
	Miniatures   []MiniatureProject `json:"miniatures,omitempty" gorm:"foreignKey:ThemeID"`
	CreatedAt    time.Time          `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time          `json:"updatedAt" gorm:"column:updated_at"`

	// Computed fields
	CoverImage string `json:"coverImage,omitempty" gorm:"-"`
}

func (MiniatureTheme) TableName() string {
	return "miniatures.miniature_themes"
}

type MiniatureProject struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	ThemeID       *int64    `json:"themeId,omitempty" gorm:"column:theme_id"`
	Title         string    `json:"name" gorm:"column:title" binding:"required"`
	Description   string    `json:"description"`
	CompletedDate *string   `json:"completedDate,omitempty" gorm:"column:completed_date"`
	Scale         string    `json:"scale,omitempty"`
	Manufacturer  string    `json:"manufacturer,omitempty"`
	TimeSpent     *float64  `json:"timeSpent,omitempty" gorm:"column:time_spent"`
	Difficulty    string    `json:"difficulty,omitempty"`
	DisplayOrder  int       `json:"displayOrder,omitempty" gorm:"column:display_order"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at"`

	// Associations
	Theme          *MiniatureTheme `json:"theme,omitempty" gorm:"foreignKey:ThemeID"`
	MiniatureFiles []MiniatureFile `json:"-" gorm:"foreignKey:MiniatureProjectID"`

	// Computed fields (populated by repository layer)
	Files      []MiniatureFile  `json:"-" gorm:"-"`
	Images     []Image          `json:"images,omitempty" gorm:"-"`
	Techniques []string         `json:"techniques,omitempty" gorm:"-"`
	Paints     []MiniaturePaint `json:"paints,omitempty" gorm:"-"`
}

func (MiniatureProject) TableName() string {
	return "miniatures.miniature_projects"
}

type MiniatureFile struct {
	ID                 int64        `json:"id" gorm:"primaryKey"`
	MiniatureProjectID int64        `json:"miniatureProjectId" gorm:"column:miniature_project_id"`
	FileID             int64        `json:"fileId" gorm:"column:file_id"`
	Caption            string       `json:"caption"`
	DisplayOrder       int          `json:"displayOrder,omitempty" gorm:"column:display_order"`
	File               *StorageFile `json:"file,omitempty" gorm:"foreignKey:FileID"`
	CreatedAt          time.Time    `json:"createdAt" gorm:"column:created_at"`
}

func (MiniatureFile) TableName() string {
	return "miniatures.miniature_files"
}

// Image is the simplified view for frontend
type Image struct {
	ID      int64  `json:"id"`
	URL     string `json:"url"`
	Caption string `json:"caption"`
}

// MiniaturePaint represents a paint in the master paint list (cl_paints table)
type MiniaturePaint struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" binding:"required"`
	Manufacturer string    `json:"manufacturer" binding:"required"`
	ColorHex     *string   `json:"colorHex,omitempty" gorm:"column:color_hex"`
	PaintType    *string   `json:"paintType,omitempty" gorm:"column:paint_type"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (MiniaturePaint) TableName() string {
	return "miniatures.cl_paints"
}
