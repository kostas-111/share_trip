package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/kostas-111/share_trip/internal/observability/logctx"
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
	logger := logctx.Logger(ctx).With(
		slog.String("service", "TripService"),
		slog.String("operation", "CreateTrip"),
		slog.String("client_id", req.ClientID.String()),
	)
	ctx = logctx.WithLogger(ctx, logger)

	logger.Info("create trip started")

	// service.CreateTripRequest -> usecase.CreateTripRequest
	trip, err := s.creator.Create(ctx, usecase.CreateTripRequest{
		ClientID:       req.ClientID,
		FromPoint:      req.FromPoint,
		ToPoint:        req.ToPoint,
		DepartureTime:  req.DepartureTime,
		AvailableSeats: req.AvailableSeats,
	})
	if err != nil {
		logger.Error("create trip failed", slog.Any("error", err))
		return nil, err
	}

	logger.Info("create trip completed", slog.String("trip_id", trip.ID.String()))

	// usecase.Trip -> service.CreateTripResponse
	return &CreateTripResponse{Trip: trip}, nil
}
