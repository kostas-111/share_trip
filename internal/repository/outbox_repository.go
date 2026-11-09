package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/kostas-111/share_trip/internal/usecase"
)

var _ usecase.OutboxRepository = (*RepoPg)(nil)

// InsertEvent пишет событие в той же транзакции, что и изменение агрегата.
// created_at не передаём - его заполняет DEFAULT NOW() в таблице.
func (r *RepoPg) InsertEvent(ctx context.Context, tx pgx.Tx, event usecase.OutboxEvent) error {
	e, err := outboxToEntity(event)
	if err != nil {
		return err
	}

	const q = `
		INSERT INTO trip_domain.outbox_event
			(id, event_name, aggregate_id, payload)
		VALUES ($1, $2, $3, $4)`

	if _, err := tx.Exec(ctx, q, e.ID, e.EventName, e.AggregateID, e.Payload); err != nil {
		return fmt.Errorf("insert outbox event: %w", err)
	}
	return nil
}
