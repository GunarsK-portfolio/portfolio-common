# Testing Guide

## Overview

The portfolio-common library uses Go's standard `testing` package for unit tests.

## Quick Commands

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific test
go test -v -run TestNewService ./jwt/

# Run all JWT tests
go test -v ./jwt/

# Run all Health tests
go test -v ./health/
```

## Test Files

**`jwt/jwt_test.go`** - 35 tests

| Category | Tests | Coverage |
|----------|-------|----------|
| Constructor | 7 | NewService, NewValidatorOnly, validation |
| Token Generation | 12 | Access tokens, refresh tokens, claims |
| Token Validation | 11 | Valid, expired, tampered, malformed |
| Expiry Handling | 5 | TTL, boundary conditions |
| Concurrency | 2 | Thread-safety verification |

**`health/health_test.go`** - 14 tests

| Category | Tests | Coverage |
|----------|-------|----------|
| Aggregator | 5 | Constructor, Register, no checkers, timeout validation |
| Health Status | 4 | Healthy, unhealthy, degraded, priority |
| Timeout | 1 | Context cancellation |
| HTTP Handler | 3 | 200 OK, 503 responses |
| Concurrency | 1 | Thread-safe Register |

**`health/postgres_test.go`** - 3 tests

| Category | Tests | Coverage |
|----------|-------|----------|
| Constructor | 2 | NewPostgresChecker, Name |
| Error Handling | 1 | Nil database |

**`health/rabbitmq_test.go`** - 3 tests

| Category | Tests | Coverage |
|----------|-------|----------|
| Constructor | 2 | NewRabbitMQChecker, Name |
| Error Handling | 1 | Nil connection |

**`health/redis_test.go`** - 3 tests

| Category | Tests | Coverage |
|----------|-------|----------|
| Constructor | 2 | NewRedisChecker, Name |
| Error Handling | 1 | Nil client |

**`health/minio_test.go`** - 4 tests

| Category | Tests | Coverage |
|----------|-------|----------|
| Constructor | 2 | NewMinIOChecker, Name |
| Error Handling | 2 | Nil client with/without bucket |

## Key Testing Patterns

**Mock Checker**: Function fields allow per-test behavior customization

```go
checker := &mockChecker{
    name:   "db",
    result: CheckResult{Status: StatusHealthy, Latency: "1ms"},
}
```

**HTTP Testing**: Uses `httptest.ResponseRecorder` with Gin router

```go
w := httptest.NewRecorder()
req, _ := http.NewRequest(http.MethodGet, "/health", nil)
router.ServeHTTP(w, req)
```

**Table-driven tests**: Multiple scenarios with `tests := []struct{...}`

**Concurrency**: Goroutines + channels for thread-safety verification

## Contributing Tests

1. Follow naming: `Test<FunctionName>_<Scenario>`
2. Organize by function with section markers
3. Use table-driven tests for multiple scenarios
4. Account for JWT second precision in timing tests
5. Clean up resources with `defer` or `t.Cleanup()`
6. Verify: `go test -cover ./...`
