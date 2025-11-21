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
```

## Test Files

**`jwt/jwt_test.go`** - 35 tests

- Constructor validation (7)
- Token generation (12)
- Token validation (11)
- Expiry handling (5)
- Concurrency (2)
- Security testing (signatures, tampering, algorithm confusion)

## Key Testing Patterns

**JWT Timestamp Precision**: JWT uses second precision, not milliseconds.
Tests use multi-second delays (1001ms) to ensure different timestamps.

**Table-driven tests**: Multiple scenarios with `tests := []struct{...}`

**Concurrency**: Goroutines + channels for thread-safety verification

**Error checking**: Tests verify specific error types
(`ErrSecretTooShort`, `ErrInvalidUserID`, etc.)

## Test Constants

```go
testSecret        = "test-secret-key-at-least-32-chars-long"
testAccessExpiry  = 15 * time.Minute
testRefreshExpiry = 168 * time.Hour  // 7 days
```

## Contributing Tests

1. Follow naming: `Test<FunctionName>_<Scenario>`
2. Organize by function with section markers
3. Use table-driven tests for multiple scenarios
4. Account for JWT second precision in timing tests
5. Clean up resources with `defer` or `t.Cleanup()`
6. Verify: `go test -cover ./...`
