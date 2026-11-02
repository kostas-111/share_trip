package service

import (
	"context"

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
	trip, err := s.reader.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &GetTripResponse{Trip: trip}, nil
}
