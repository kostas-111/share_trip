package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func WithTx[T any](
	ctx context.Context,
	pool *pgxpool.Pool,
	block func(tx pgx.Tx) (*T, error),
) (*T, error) {
	txBegin, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = txBegin.Rollback(ctx) }()

	res, err := block(txBegin)
	if err != nil {
		return nil, err // ошибку уже обернули слои ниже
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := txBegin.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}
	return res, nil
}

func WithTxErr(
	ctx context.Context,
	pool *pgxpool.Pool,
	block func(tx pgx.Tx) error,
) error {
	_, err := WithTx(ctx, pool, func(tx pgx.Tx) (*struct{}, error) {
		if err := block(tx); err != nil {
			return nil, err
		}
		return &struct{}{}, nil
	})
	return err
}

type Runner interface {
	WithTx(ctx context.Context, block func(tx pgx.Tx) error) error
}

func (r *RepoPg) WithTx(ctx context.Context, block func(tx pgx.Tx) error) error {
	return WithTxErr(ctx, r.pool, block)
}
