package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/kostas-111/share_trip/internal/usecase"
)

var _ usecase.OutboxRepository = (*RepoPg)(nil)

// InsertEvent пишет событие в той же транзакции, что и изменение агрегата.
// created_at не передаём - его заполняет DEFAULT NOW() в таблице.
func (r *RepoPg) InsertEvent(ctx context.Context, tx pgx.Tx, event usecase.OutboxEvent) error {
	started := time.Now()
	result := "success"

	defer func() {
		r.metrics.RepositoryQueryTotal.WithLabelValues("outbox_insert_event", result).Inc()
		r.metrics.RepositoryQueryDuration.WithLabelValues("outbox_insert_event", result).
			Observe(time.Since(started).Seconds())
	}()

	e, err := outboxToEntity(event)
	if err != nil {
		result = "error"
		return err
	}

	const q = `
		INSERT INTO trip_domain.outbox_event
			(id, event_name, aggregate_id, payload)
		VALUES ($1, $2, $3, $4)`

	if _, err := tx.Exec(ctx, q, e.ID, e.EventName, e.AggregateID, e.Payload); err != nil {
		result = "error"
		return fmt.Errorf("insert outbox event: %w", err)
	}
	return nil
}
