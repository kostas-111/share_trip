package service

import (
	"context"
	"errors"
	"time"

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
	started := time.Now()
	result := "success"

	defer func() {
		s.metrics.TripPublishTotal.WithLabelValues(result).Inc()
		s.metrics.TripPublishDuration.WithLabelValues(result).
			Observe(time.Since(started).Seconds())
	}()

	publishResult, err := s.publisher.Publish(ctx, req.TripID, req.ClientID)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrValidation):
			result = "validation_error"
		case errors.Is(err, usecase.ErrTripNotFound):
			result = "not_found"
		case errors.Is(err, usecase.ErrForbidden):
			result = "forbidden"
		case errors.Is(err, usecase.ErrInvalidStatus):
			result = "conflict"
		default:
			result = "internal_error"
		}

		return nil, err
	}

	if publishResult.AlreadyPublished {
		result = "already_published"
	}

	return &PublishTripResponse{
		Trip:             publishResult.Trip,
		AlreadyPublished: publishResult.AlreadyPublished,
	}, nil
}
