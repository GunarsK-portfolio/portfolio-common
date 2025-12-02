# logger

Structured logging with slog and Gin middleware.

## Usage

```go
import "github.com/GunarsK-portfolio/portfolio-common/logger"

appLogger := logger.New(logger.Config{
    Level:       "info",
    Format:      "json",
    ServiceName: "my-service",
    AddSource:   false,
})

// Use middleware
router.Use(logger.Recovery(appLogger))
router.Use(logger.RequestLogger(appLogger))

// Log messages
appLogger.Info("Server started", "port", 8080)
appLogger.Error("Database error", "error", err)
```

## Log Levels

- `debug` - Detailed debugging info
- `info` - General operational info
- `warn` - Warning conditions
- `error` - Error conditions

## Formats

- `json` - Structured JSON (production)
- `text` - Human-readable (development)
