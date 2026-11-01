package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kostas-111/share_trip/internal/domain"
	"github.com/kostas-111/share_trip/internal/repository"
)

type TripPublisher struct {
	repo   repository.TripRepository
	runner repository.Runner
}

func NewTripPublisher(repo repository.TripRepository, runner repository.Runner) *TripPublisher {
	return &TripPublisher{repo: repo, runner: runner}
}

func (p *TripPublisher) Publish(ctx context.Context, tripID uuid.UUID) (*domain.Trip, error) {
	if tripID == uuid.Nil {
		return nil, fmt.Errorf("%w: trip_id is required", domain.ErrValidation)
	}

	var published *domain.Trip
	err := p.runner.WithTx(ctx, func(tx pgx.Tx) error {
		trip, err := p.repo.GetByIDForUpdate(ctx, tx, tripID)
		if err != nil {
			return err
		}

		domResp, err := domain.PublishTrip(trip)
		if err != nil {
			return err
		}

		if err := p.repo.UpdateStatus(ctx, tx, trip.ID, domResp.ToStatus); err != nil {
			return err
		}
		if err := p.repo.InsertHistory(ctx, tx, trip.ID, domResp.FromStatus, domResp.ToStatus); err != nil {
			return err
		}

		domResp.Trip.Status = domResp.ToStatus
		published = domResp.Trip
		return nil
	})
	if err != nil {
		return nil, err
	}

	return published, nil
}
