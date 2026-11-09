package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/kostas-111/share_trip/internal/usecase"
)

type PublishTripRequest struct {
	TripID   uuid.UUID
	ClientID uuid.UUID
}

type PublishTripResponse struct {
	Trip             *usecase.Trip
	AlreadyPublished bool
}

func (s *TripService) PublishTrip(ctx context.Context, req PublishTripRequest) (*PublishTripResponse, error) {
	result, err := s.publisher.Publish(ctx, req.TripID, req.ClientID)
	if err != nil {
		return nil, err
	}

	return &PublishTripResponse{
		Trip:             result.Trip,
		AlreadyPublished: result.AlreadyPublished,
	}, nil
}
