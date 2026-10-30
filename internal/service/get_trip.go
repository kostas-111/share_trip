package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kostas-111/share_trip/internal/domain"
)

type GetTripRequest struct {
	ID uuid.UUID
}

type GetTripResponse struct {
	Trip *domain.Trip
}

func (s *TripService) GetTrip(ctx context.Context, req GetTripRequest) (*GetTripResponse, error) {
	if err := (domain.GetTripParams{ID: req.ID}).Validate(); err != nil {
		return nil, err
	}

	trip, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("TripService.GetTrip: %w", err)
	}
	if trip == nil {
		return nil, fmt.Errorf("%w: id=%s", domain.ErrTripNotFound, req.ID)
	}

	return &GetTripResponse{Trip: trip}, nil
}
