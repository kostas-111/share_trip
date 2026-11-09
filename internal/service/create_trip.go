package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/kostas-111/share_trip/internal/usecase"
)

type CreateTripRequest struct {
	ClientID       uuid.UUID
	FromPoint      string
	ToPoint        string
	DepartureTime  time.Time
	AvailableSeats int
}

type CreateTripResponse struct {
	Trip *usecase.Trip
}

func (s *TripService) CreateTrip(ctx context.Context, req CreateTripRequest) (*CreateTripResponse, error) {
	// service.CreateTripRequest -> usecase.CreateTripRequest
	trip, err := s.creator.Create(ctx, usecase.CreateTripRequest{
		ClientID:       req.ClientID,
		FromPoint:      req.FromPoint,
		ToPoint:        req.ToPoint,
		DepartureTime:  req.DepartureTime,
		AvailableSeats: req.AvailableSeats,
	})
	if err != nil {
		return nil, err
	}

	// usecase.Trip -> service.CreateTripResponse
	return &CreateTripResponse{Trip: trip}, nil
}
