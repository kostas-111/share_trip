package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	EventTripPublished = "trip_published"
)

type OutboxEvent struct {
	ID          uuid.UUID
	EventName   string
	AggregateID uuid.UUID
	Payload     any
}

type TripPublishedPayload struct {
	TripID uuid.UUID `json:"trip_id"`
}

func NewTripPublishedEvent(tripID uuid.UUID) OutboxEvent {
	return OutboxEvent{
		ID:          uuid.New(),
		EventName:   EventTripPublished,
		AggregateID: tripID,
		Payload:     TripPublishedPayload{TripID: tripID},
	}
}

// OutboxRepository пишет исходящее событие в той же транзакции, что и изменение агрегата.
type OutboxRepository interface {
	InsertEvent(ctx context.Context, tx pgx.Tx, event OutboxEvent) error
}
