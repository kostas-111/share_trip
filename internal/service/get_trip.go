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
	// service -> domain
	domainResp, err := domain.GetTrip(domain.GetTripRequest{ID: req.ID})
	if err != nil {
		return nil, err
	}

	trip, err := s.repo.GetByID(ctx, domainResp.ID)
	if err != nil {
		return nil, fmt.Errorf("get trip: %w", err) // внутренняя ошибка → 500
	}
	if trip == nil {
		return nil, fmt.Errorf("%w: id=%s", domain.ErrTripNotFound, req.ID)
	}

	return &GetTripResponse{Trip: trip}, nil
}
