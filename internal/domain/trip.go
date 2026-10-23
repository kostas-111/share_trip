package domain

import (
	"time"

	"github.com/google/uuid"
)

type TripStatus string

const (
	StatusDraft     TripStatus = "draft"
	StatusPublished TripStatus = "published"
	StatusCanceled  TripStatus = "canceled"
	StatusCompleted TripStatus = "completed"
)

type Trip struct {
	ID              uuid.UUID
	ClientID        uuid.UUID
	FromPoint       string
	ToPoint         string
	DepartureTime   time.Time
	AvailableSeats  int
	Status          TripStatus
	CreatedAt       time.Time
}