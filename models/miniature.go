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

	// Associations
	CoverImageFile *StorageFile `json:"coverImageFile,omitempty" gorm:"foreignKey:CoverImageID"`
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
	Theme          *MiniatureTheme             `json:"theme,omitempty" gorm:"foreignKey:ThemeID"`
	MiniatureFiles []MiniatureFile             `json:"-" gorm:"foreignKey:MiniatureProjectID"`
	Techniques     []MiniatureProjectTechnique `json:"techniques,omitempty" gorm:"foreignKey:MiniatureProjectID"`
	Paints         []MiniatureProjectPaint     `json:"paints,omitempty" gorm:"foreignKey:MiniatureProjectID"`

	// Computed field (populated by repository layer - requires URL building)
	Images []Image `json:"images,omitempty" gorm:"-"`
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

// MiniatureTechnique represents a painting technique in the master technique list (cl_techniques table).
// This model is used for managing the technique catalog with CRUD operations.
type MiniatureTechnique struct {
	ID              int64     `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" binding:"required"`
	Description     string    `json:"description"`
	DifficultyLevel string    `json:"difficultyLevel,omitempty" gorm:"column:difficulty_level"`
	DisplayOrder    int       `json:"displayOrder,omitempty" gorm:"column:display_order"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (MiniatureTechnique) TableName() string {
	return "miniatures.cl_techniques"
}

// MiniatureProjectTechnique is the junction table linking projects to techniques
type MiniatureProjectTechnique struct {
	ID                 int64     `json:"id" gorm:"primaryKey"`
	MiniatureProjectID int64     `json:"miniatureProjectId" gorm:"column:miniature_project_id"`
	TechniqueID        int64     `json:"techniqueId" gorm:"column:technique_id"`
	Notes              string    `json:"notes"`
	CreatedAt          time.Time `json:"createdAt" gorm:"column:created_at"`

	// Associations
	Technique *MiniatureTechnique `json:"technique,omitempty" gorm:"foreignKey:TechniqueID"`
}

func (MiniatureProjectTechnique) TableName() string {
	return "miniatures.miniature_techniques"
}

// MiniatureProjectPaint is the junction table linking projects to paints
type MiniatureProjectPaint struct {
	ID                 int64     `json:"id" gorm:"primaryKey"`
	MiniatureProjectID int64     `json:"miniatureProjectId" gorm:"column:miniature_project_id"`
	PaintID            int64     `json:"paintId" gorm:"column:paint_id"`
	UsageNotes         string    `json:"usageNotes,omitempty" gorm:"column:usage_notes"`
	CreatedAt          time.Time `json:"createdAt" gorm:"column:created_at"`

	// Associations
	Paint *MiniaturePaint `json:"paint,omitempty" gorm:"foreignKey:PaintID"`
}

func (MiniatureProjectPaint) TableName() string {
	return "miniatures.miniature_paints"
}

// MiniaturePaint represents a paint in the master paint list (cl_paints table).
// This model is used for managing the paint catalog with CRUD operations.
// Paints are categorized by type and identified by the combination of name and manufacturer.
type MiniaturePaint struct {
	ID           int64  `json:"id" gorm:"primaryKey"`
	Name         string `json:"name" binding:"required"`
	Manufacturer string `json:"manufacturer" binding:"required"`
	// ColorHex is the hexadecimal color code in #RRGGBB or #RGB format (e.g., #FF5733, #F00)
	ColorHex *string `json:"colorHex,omitempty" gorm:"column:color_hex" binding:"omitempty,hexcolor"`
	// PaintType categorizes the paint (Base, Layer, Shade, Wash, Contrast, Dry, Technical, Metallic, Air, Primer, Edge, Glaze, Ink)
	// Database enforces these values via CHECK constraint
	PaintType *string   `json:"paintType,omitempty" gorm:"column:paint_type"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (MiniaturePaint) TableName() string {
	return "miniatures.cl_paints"
}
