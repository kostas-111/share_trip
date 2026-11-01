package usecase

import (
	"context"
	"fmt"

	"github.com/kostas-111/share_trip/internal/domain"
	"github.com/kostas-111/share_trip/internal/repository"
)

type TripCreator struct {
	repo repository.TripRepository
}

func NewTripCreator(repo repository.TripRepository) *TripCreator {
	return &TripCreator{repo: repo}
}

func (c *TripCreator) Create(ctx context.Context, req domain.CreateTripRequest) (*domain.Trip, error) {
	domainResp, err := domain.CreateTrip(req)
	if err != nil {
		return nil, err
	}

	if err := c.repo.Create(ctx, domainResp.Trip); err != nil {
		return nil, fmt.Errorf("TripCreator.Create: %w", err)
	}

	return domainResp.Trip, nil
}
