# metrics

Prometheus metrics collection with Gin middleware.

## Usage

```go
import "github.com/GunarsK-portfolio/portfolio-common/metrics"

metricsCollector := metrics.New(metrics.Config{
    ServiceName: "admin",
    Namespace:   "portfolio",
})

// HTTP metrics middleware
router.Use(metricsCollector.Middleware())

// Expose metrics endpoint
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

## Collected Metrics

- `http_requests_total` - Total HTTP requests by method, path, status
- `http_request_duration_seconds` - Request latency histogram
