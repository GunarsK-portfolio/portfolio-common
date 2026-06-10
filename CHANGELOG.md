# Changelog

## v0.51.0

Queue stack upgrade for long-running workers.

**Behavior changes** (no source changes required, review before upgrading):

- Publisher and consumer now reconnect automatically with exponential
  backoff after a dropped connection or channel, re-declaring topology.
  `Consume` no longer returns when the delivery channel closes; it resumes
  consumption and only returns on context cancellation, `Close()`, or after
  `RABBITMQ_RECONNECT_MAX_ATTEMPTS`. Set `RABBITMQ_RECONNECT=false`
  (`DisableReconnect`) for the old fail-fast behavior.
- `consumer.Close()` now waits for in-flight handlers to finish.
- If the consume context is cancelled while a message is being handled, the
  message is requeued without consuming a retry attempt (previously it
  entered the retry ladder).
- `Publish` uses the context correlation ID (`logger.GetCorrelationID`) as
  the message CorrelationId when present; the consumer injects it back into
  the handler context.
- Invalid numeric `RABBITMQ_*` values now panic at startup instead of
  falling back to defaults with a warning.
- `GetRetryCount` accepts int8/int16 headers and clamps negatives to 0.

**Additive**:

- `queue.Permanent(err)` / `queue.ErrPermanent` route a message straight to
  the DLQ without retries.
- Optional publisher confirms (`RABBITMQ_PUBLISHER_CONFIRMS`), returning
  `queue.ErrPublishNotConfirmed` on broker NACK.
- Optional retry jitter (`RABBITMQ_RETRY_JITTER`, 0-1 fraction) via
  per-message TTLs.
- Optional concurrent consumption (`RABBITMQ_CONSUMER_CONCURRENCY`).
- Configurable heartbeat (`RABBITMQ_HEARTBEAT`, default 10s) and client
  connection name (from `RABBITMQ_CONSUMER_TAG`).
- `config.NewRabbitMQConfigWithPrefix(prefix)` for per-queue env config
  with fallback to un-prefixed names.
- `queue.WithPublisherLogger/WithPublisherMetrics/WithConsumerMetrics`
  options; `metrics.NewQueueMetrics` Prometheus recorder (publish/consume
  totals and durations, reconnects, queue depth gauge).
- `health.NewRabbitMQCheckerWithProvider(provider)` - required with
  reconnection, the fixed-connection checker goes stale after a reconnect.
- `health.NewQueueDepthChecker(provider, queue, threshold)` exposing DLQ
  depth under `details.messages`; `CheckResult` gained a `Details` field.
- `config.RabbitMQConfig.WithDefaults()` normalizing zero-valued optional
  fields; applied automatically by the queue constructors so struct-literal
  configs behave like env-loaded ones.
- Integration test suite (testcontainers-go) covering publish/consume,
  retry ladder, DLQ, permanent errors, reconnection, confirms, shutdown
  semantics, and concurrency.

**Call-site recommendation**: replace
`health.NewRabbitMQChecker(publisher.Connection())` with
`health.NewRabbitMQCheckerWithProvider(publisher.Connection)`.

## v0.34.0

**BREAKING**: JWT token generation methods now require a `scopes` parameter.

```go
// Before (v0.33.0)
token, err := jwtService.GenerateAccessToken(userID, username)
refreshToken, err := jwtService.GenerateRefreshToken(userID, username)

// After (v0.34.0)
scopes := map[string]string{"profile": "read", "projects": "edit"}
token, err := jwtService.GenerateAccessToken(userID, username, scopes)
refreshToken, err := jwtService.GenerateRefreshToken(userID, username, scopes)

// For nil scopes (no permissions)
token, err := jwtService.GenerateAccessToken(userID, username, nil)
```

- `GenerateAccessToken(userID, username)` now requires third `scopes` param
- `GenerateRefreshToken(userID, username)` now requires third `scopes` param
- Add `Scopes` field to JWT `Claims` struct
- Add `middleware/permission.go` with `RequirePermission()` middleware
- Add permission level constants: `LevelNone`, `LevelRead`, `LevelEdit`, `LevelDelete`
- Auth middleware now extracts scopes from JWT into Gin context

## v0.33.0

- Add `health` package for dependency health checking
- Add `Connection()` method to RabbitMQPublisher for health checks
- Add per-package README files
- Restructure main README to link to package docs

## v0.32.0

- Add `CloseDB` helper function to database package
- Add `queue` package with RabbitMQ publisher, retry queues, and DLQ support

## v0.21.0

**BREAKING**: AuthMiddleware now uses local JWT validation.

```go
// Before (v0.20.0)
authMiddleware := middleware.NewAuthMiddleware("http://auth-service:8084/api/v1")

// After (v0.21.0)
jwtService, _ := jwt.NewValidatorOnly(os.Getenv("JWT_SECRET"))
authMiddleware := middleware.NewAuthMiddleware(jwtService)
```

- `NewAuthMiddleware(authServiceURL string, opts...)` changed to
  `NewAuthMiddleware(jwtService jwt.Service)`
- Add `jwt` package for local token validation
- Remove `WithTimeout` option (no longer needed)
- Services must provide `JWT_SECRET` environment variable

## v0.19.0

- Add `SSLMode` field to `DatabaseConfig` with validation
- New environment variable: `DB_SSLMODE` (optional, default: `disable`)
- Valid values: `disable`, `allow`, `prefer`, `require`, `verify-ca`, `verify-full`

## v0.12.0

**BREAKING**: `ALLOWED_ORIGINS` environment variable is now required.

- No default value provided for security reasons
- Services will panic on startup if not configured
- Use comma-separated list: `ALLOWED_ORIGINS=http://localhost:8080,https://example.com`
