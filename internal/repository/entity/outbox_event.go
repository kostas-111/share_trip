package entity

import (
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID          uuid.UUID `db:"id"`
	EventName   string    `db:"event_name"`
	AggregateID uuid.UUID `db:"aggregate_id"`
	Payload     []byte    `db:"payload"`
	CreatedAt   time.Time `db:"created_at"`
}
