# database

Database connection utilities with GORM and connection pooling.

## Usage

```go
import "github.com/GunarsK-portfolio/portfolio-common/database"

db, err := database.Connect(database.PostgresConfig{
    Host:     "localhost",
    Port:     "5432",
    User:     "user",
    Password: "password",
    DBName:   "portfolio",
    SSLMode:  "disable",
    TimeZone: "UTC",
})
if err != nil {
    log.Fatal(err)
}
defer database.CloseDB(db)
```

## Functions

- `Connect(cfg PostgresConfig) (*gorm.DB, error)` - Connect to PostgreSQL
- `CloseDB(db *gorm.DB) error` - Close database connection
