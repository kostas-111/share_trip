package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Runner выполняет блок в одной транзакции; реализация — *repository.RepoPg.
type Runner interface {
	WithTx(ctx context.Context, block func(tx pgx.Tx) error) error
}

type TripPublisher struct {
	repo   TripRepository
	runner Runner
	outbox OutboxRepository
}

func NewTripPublisher(
	repo TripRepository,
	runner Runner,
	outbox OutboxRepository,
) *TripPublisher {
	return &TripPublisher{repo: repo, runner: runner, outbox: outbox}
}

// Publish переводит поездку DRAFT -> PUBLISHED одной транзакцией:
// статус, запись истории перехода и событие в outbox.
func (p *TripPublisher) Publish(ctx context.Context, tripID uuid.UUID) (*Trip, error) {
	if tripID == uuid.Nil {
		return nil, fmt.Errorf("%w: trip_id is required", ErrValidation)
	}

	var published *Trip
	err := p.runner.WithTx(ctx, func(tx pgx.Tx) error {
		trip, err := p.repo.GetByIDForUpdate(ctx, tx, tripID)
		if err != nil {
			return err
		}
		if trip == nil {
			return fmt.Errorf("%w: trip not found", ErrTripNotFound)
		}

		from := trip.Status
		if from != StatusDraft {
			return fmt.Errorf(
				"%w: only draft trips can be published (current status: %s)",
				ErrInvalidStatus, from,
			)
		}
		if !trip.DepartureTime.After(time.Now()) {
			return fmt.Errorf(
				"%w: cannot publish a trip with departure_time in the past",
				ErrValidation,
			)
		}
		to := StatusPublished

		if err := p.repo.UpdateStatus(ctx, tx, trip.ID, to); err != nil {
			return err
		}
		if err := p.repo.InsertHistory(ctx, tx, trip.ID, from, to); err != nil {
			return err
		}
		if err := p.outbox.InsertEvent(ctx, tx, NewTripPublishedEvent(trip.ID)); err != nil {
			return err
		}

		trip.Status = to
		published = trip
		return nil
	})
	if err != nil {
		return nil, err
	}

	return published, nil
}
