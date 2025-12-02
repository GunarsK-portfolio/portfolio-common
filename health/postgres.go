package health

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PostgresChecker checks PostgreSQL database connectivity
type PostgresChecker struct {
	db *gorm.DB
}

// NewPostgresChecker creates a new PostgreSQL health checker
func NewPostgresChecker(db *gorm.DB) Checker {
	return &PostgresChecker{db: db}
}

// Name returns the name of this checker
func (c *PostgresChecker) Name() string {
	return "postgres"
}

// Check verifies the database connection is alive
func (c *PostgresChecker) Check(ctx context.Context) CheckResult {
	start := time.Now()

	if c.db == nil {
		return CheckResult{
			Status:  StatusUnhealthy,
			Latency: time.Since(start).String(),
			Error:   "database is nil",
		}
	}

	sqlDB, err := c.db.DB()
	if err != nil {
		return CheckResult{
			Status:  StatusUnhealthy,
			Latency: time.Since(start).String(),
			Error:   fmt.Sprintf("failed to get database instance: %v", err),
		}
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return CheckResult{
			Status:  StatusUnhealthy,
			Latency: time.Since(start).String(),
			Error:   fmt.Sprintf("ping failed: %v", err),
		}
	}

	return CheckResult{
		Status:  StatusHealthy,
		Latency: time.Since(start).String(),
	}
}
