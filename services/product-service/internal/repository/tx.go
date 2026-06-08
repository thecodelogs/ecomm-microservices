package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

// WithTx returns a new context with the transaction.
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// GetTx returns the transaction from context, or nil if none.
func GetTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}

// DBExecutor defines common methods of pgxpool.Pool and pgx.Tx
type DBExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// getDb returns the transaction from the context if present, otherwise returns the pool.
func getDb(ctx context.Context, db *pgxpool.Pool) DBExecutor {
	if tx := GetTx(ctx); tx != nil {
		return tx
	}
	return db
}
