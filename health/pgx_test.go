package health

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestNewPgxChecker(t *testing.T) {
	checker := NewPgxChecker(nil)

	if checker == nil {
		t.Fatal("expected checker to not be nil")
	}
}

func TestPgxChecker_Name(t *testing.T) {
	checker := NewPgxChecker(nil)

	if checker.Name() != "postgres" {
		t.Errorf("expected name 'postgres', got %s", checker.Name())
	}
}

func TestPgxChecker_Check_NilPool(t *testing.T) {
	checker := NewPgxChecker(nil)

	result := checker.Check(context.Background())

	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy status for nil pool, got %s", result.Status)
	}
	if result.Error != "pool is nil" {
		t.Errorf("expected 'pool is nil' error, got %s", result.Error)
	}
}

func TestPgxChecker_Check_PingFailure(t *testing.T) {
	// pgxpool connects lazily, so a pool pointing at a closed port builds
	// fine and only the ping fails.
	poolCfg, err := pgxpool.ParseConfig("postgres://user:pass@127.0.0.1:1/db?sslmode=disable")
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	defer pool.Close()

	result := NewPgxChecker(pool).Check(ctx)

	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy status for unreachable database, got %s", result.Status)
	}
	if !strings.HasPrefix(result.Error, "ping failed:") {
		t.Errorf("expected error to start with 'ping failed:', got %s", result.Error)
	}
	if result.Latency == "" {
		t.Error("expected latency to be reported")
	}
}
