package domain

import (
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

type CreateTripResponse struct {
	Trip *Trip
}

/*
Здесь только бизнес-правила.
*
*/
func CreateTrip(req CreateTripRequest) (CreateTripResponse, error) {
	if req.ClientID == uuid.Nil {
		return CreateTripResponse{}, fmt.Errorf("%w: client_id is required", ErrValidation)
	}
	if req.FromPoint == "" {
		return CreateTripResponse{}, fmt.Errorf("%w: from_point is required", ErrValidation)
	}
	if req.ToPoint == "" {
		return CreateTripResponse{}, fmt.Errorf("%w: to_point is required", ErrValidation)
	}
	if req.DepartureTime.IsZero() {
		return CreateTripResponse{}, fmt.Errorf("%w: departure_time is required", ErrValidation)
	}
	if !req.DepartureTime.After(time.Now()) {
		return CreateTripResponse{}, fmt.Errorf("%w: departure_time must be in the future", ErrValidation)
	}
	if req.AvailableSeats <= 0 {
		return CreateTripResponse{}, fmt.Errorf("%w: available_seats must be positive", ErrValidation)
	}

	return CreateTripResponse{
		Trip: &Trip{
			ID:             uuid.New(),
			ClientID:       req.ClientID,
			FromPoint:      req.FromPoint,
			ToPoint:        req.ToPoint,
			DepartureTime:  req.DepartureTime.UTC(),
			AvailableSeats: req.AvailableSeats,
			Status:         StatusDraft,
			CreatedAt:      time.Now().UTC(),
		},
	}, nil
}
