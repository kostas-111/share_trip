package domain

import (
	"github.com/google/uuid"
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
