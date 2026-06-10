# database

Database connection utilities for GORM and pgx with connection pooling.

## GORM Usage

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

## pgx Usage

For services using pgx directly instead of GORM. Builds the pool from the
shared `config.DatabaseConfig` and verifies connectivity with a ping. The
appName argument is reported as the PostgreSQL `application_name`.

```go
import "github.com/GunarsK-portfolio/portfolio-common/database"

pool, err := database.NewPgxPool(ctx, cfg.DatabaseConfig, "my-service")
if err != nil {
    log.Fatal(err)
}
defer pool.Close()

// Larger pool for a busier service
pool, err = database.NewPgxPool(ctx, cfg.DatabaseConfig, "my-api",
    database.WithPoolSize(25, 5))
```

Defaults: MaxConns 10, MinConns 2, MaxConnLifetime 1h, MaxConnIdleTime 10m,
HealthCheckPeriod 30s, sslmode `disable` when the config leaves it empty.
Options receive the `*pgxpool.Config` after defaults are applied and may
change any field.

## Functions

- `Connect(cfg PostgresConfig) (*gorm.DB, error)` - Connect to PostgreSQL via GORM
- `CloseDB(db *gorm.DB) error` - Close database connection
- `NewPgxPool(ctx, cfg, appName, opts...)` - Create a pgx connection pool
  from `config.DatabaseConfig`
- `WithPoolSize(maxConns, minConns int32) PgxPoolOption` - Set pool sizing
