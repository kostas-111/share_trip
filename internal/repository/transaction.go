package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kostas-111/share_trip/internal/observability/logctx"
)

func WithTx[T any](
	ctx context.Context,
	pool *pgxpool.Pool,
	block func(tx pgx.Tx) (*T, error),
) (*T, error) {
	logger := logctx.Logger(ctx).With(
		slog.String("layer", "transaction"),
	)

	logger.Info("begin transaction")

	txBegin, err := pool.Begin(ctx)
	if err != nil {
		logger.Error("failed to begin transaction", slog.Any("error", err))
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err := txBegin.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logger.Error("rollback transaction failed", slog.Any("error", err))
		}
	}()

	res, err := block(txBegin)
	if err != nil {
		logger.Error("transaction block failed", slog.Any("error", err))
		return nil, err // ошибку уже обернули слои ниже
	}

	if err := ctx.Err(); err != nil {
		logger.Error("transaction aborted before commit", slog.Any("error", err))
		return nil, err
	}

	if err := txBegin.Commit(ctx); err != nil {
		logger.Error("failed to commit transaction", slog.Any("error", err))
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	logger.Info("commit transaction")

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

func (r *RepoPg) WithTx(ctx context.Context, block func(tx pgx.Tx) error) error {
	return WithTxErr(ctx, r.pool, block)
}
