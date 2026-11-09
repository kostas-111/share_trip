package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreateTripRequest struct {
	ClientID       uuid.UUID
	FromPoint      string
	ToPoint        string
	DepartureTime  time.Time
	AvailableSeats int
}

type TripCreator struct {
	repo TripRepository
}

func NewTripCreator(repo TripRepository) *TripCreator {
	return &TripCreator{repo: repo}
}

// Create проверяет вход, собирает черновик поездки и сохраняет его.
func (c *TripCreator) Create(ctx context.Context, req CreateTripRequest) (*Trip, error) {
	if req.ClientID == uuid.Nil {
		return nil, fmt.Errorf("%w: client_id is required", ErrValidation)
	}
	if req.FromPoint == "" {
		return nil, fmt.Errorf("%w: from_point is required", ErrValidation)
	}
	if req.ToPoint == "" {
		return nil, fmt.Errorf("%w: to_point is required", ErrValidation)
	}
	if req.DepartureTime.IsZero() {
		return nil, fmt.Errorf("%w: departure_time is required", ErrValidation)
	}
	if !req.DepartureTime.After(time.Now()) {
		return nil, fmt.Errorf("%w: departure_time must be in the future", ErrValidation)
	}
	if req.AvailableSeats <= 0 {
		return nil, fmt.Errorf("%w: available_seats must be positive", ErrValidation)
	}

	trip := &Trip{
		ID:             uuid.New(),
		ClientID:       req.ClientID,
		FromPoint:      req.FromPoint,
		ToPoint:        req.ToPoint,
		DepartureTime:  req.DepartureTime.UTC(),
		AvailableSeats: req.AvailableSeats,
		Status:         StatusDraft,
		CreatedAt:      time.Now().UTC(),
	}

	if err := c.repo.Create(ctx, trip); err != nil {
		return nil, fmt.Errorf("TripCreator.Create: %w", err)
	}

	return trip, nil
}
