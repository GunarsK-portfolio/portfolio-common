# models

Shared GORM database models for portfolio services.

## Models

- `Profile` - User profile information
- `WorkExperience` - Work history entries
- `Certification` - Professional certifications
- `Skill` - Technical skills with proficiency levels
- `Project` - Portfolio projects
- `Miniature` - Miniature painting showcase
- `MiniatureTheme` - Miniature themes/categories
- `StorageFile` - File metadata for MinIO storage
- `User` - User accounts
- `ActionLog` - Audit log entries
- `Recipient` - Email recipients for messaging
- `ContactMessage` - Contact form submissions

## Usage

```go
import "github.com/GunarsK-portfolio/portfolio-common/models"

var profile models.Profile
db.First(&profile)
```
