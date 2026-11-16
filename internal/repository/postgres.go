package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kostas-111/share_trip/internal/observability/metrics"
)

type RepoPg struct {
	pool    *pgxpool.Pool
	metrics *metrics.Metrics
}

func NewRepoPg(pool *pgxpool.Pool, m *metrics.Metrics) *RepoPg {
	return &RepoPg{pool: pool, metrics: m}
}

func (r *RepoPg) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}
