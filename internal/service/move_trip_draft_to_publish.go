package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/kostas-111/share_trip/internal/domain"
)

type PublishTripRequest struct {
	TripID uuid.UUID
}

type PublishTripResponse struct {
	Trip *domain.Trip
}

func (s *TripService) PublishTrip(ctx context.Context, req PublishTripRequest) (*PublishTripResponse, error) {
	trip, err := s.publisher.Publish(ctx, req.TripID)
	if err != nil {
		return nil, err
	}

	return &PublishTripResponse{Trip: trip}, nil
}
