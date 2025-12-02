# server

HTTP server utilities with graceful shutdown.

## Usage

```go
import "github.com/GunarsK-portfolio/portfolio-common/server"

cfg := server.DefaultConfig("8080")
if err := server.Run(router, cfg, logger); err != nil {
    log.Fatal(err)
}
```

## Features

- Graceful shutdown on SIGINT/SIGTERM
- Configurable timeouts
- Structured logging integration
