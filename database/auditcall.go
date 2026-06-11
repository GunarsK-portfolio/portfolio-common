package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthContext is the audit.set_context() payload: the session identity
// applied before every audited SQL function call.
type AuthContext struct {
	UserID    int64
	Username  string
	ClientIP  string
	UserAgent string
}

// CallInto runs `SELECT audit.set_context(...)` and the given single-row
// query in one pgx batch, scanning the result into dest. A batch executes in
// one implicit transaction, so the transaction-local audit settings cover the
// query and an error in either statement aborts both - the semantics of an
// explicit BEGIN/COMMIT exchange in one network round trip.
func CallInto(ctx context.Context, pool *pgxpool.Pool, auth AuthContext, dest any, query string, args ...any) error {
	b := &pgx.Batch{}
	b.Queue("SELECT audit.set_context($1, $2, $3, $4)",
		auth.UserID, auth.Username, auth.ClientIP, auth.UserAgent)
	b.Queue(query, args...)

	br := pool.SendBatch(ctx, b)
	defer func() { _ = br.Close() }()

	if _, err := br.Exec(); err != nil {
		return fmt.Errorf("set audit context: %w", err)
	}
	if err := br.QueryRow().Scan(dest); err != nil {
		return err
	}
	return br.Close()
}

// CallJSON runs an audited query returning JSONB (the SELECT schema.fn(...)
// convention).
func CallJSON(ctx context.Context, pool *pgxpool.Pool, auth AuthContext, query string, args ...any) (json.RawMessage, error) {
	var result json.RawMessage
	err := CallInto(ctx, pool, auth, &result, query, args...)
	return result, err
}

// CallBool runs an audited query returning a boolean (delete functions).
func CallBool(ctx context.Context, pool *pgxpool.Pool, auth AuthContext, query string, args ...any) (bool, error) {
	var result bool
	err := CallInto(ctx, pool, auth, &result, query, args...)
	return result, err
}

// CallDiscard runs an audited query and discards its scalar result.
func CallDiscard(ctx context.Context, pool *pgxpool.Pool, auth AuthContext, query string, args ...any) error {
	var discard any
	return CallInto(ctx, pool, auth, &discard, query, args...)
}
