package domain

import (
	"github.com/google/uuid"
)

// Имена событий, которые пишутся в outbox.
const (
	EventTripPublished = "trip_published"
)

// OutboxEvent — доменное представление строки таблицы outbox_event.
// Payload остаётся структурой; в JSONB его сериализует repository.
type OutboxEvent struct {
	ID          uuid.UUID
	EventName   string
	AggregateID uuid.UUID
	Payload     any
}

// TripPublishedPayload - полезная нагрузка события trip_published.
type TripPublishedPayload struct {
	TripID uuid.UUID `json:"trip_id"`
}

// NewTripPublishedEvent - событие успешного перевода поездки в published.
func NewTripPublishedEvent(tripID uuid.UUID) OutboxEvent {
	return OutboxEvent{
		ID:          uuid.New(),
		EventName:   EventTripPublished,
		AggregateID: tripID,
		Payload:     TripPublishedPayload{TripID: tripID},
	}
}
