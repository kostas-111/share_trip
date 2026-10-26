package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/kostas-111/share_trip/internal/domain"
	"github.com/kostas-111/share_trip/internal/repository/entity"
)

type TripRepository interface {
	Create(ctx context.Context, trip *domain.Trip) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Trip, error)
	GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*domain.Trip, error)
	UpdateStatus(ctx context.Context, tx pgx.Tx, id uuid.UUID, status domain.TripStatus) error
	InsertHistory(ctx context.Context, tx pgx.Tx, tripID uuid.UUID, from, to domain.TripStatus) error
}

var _ TripRepository = (*RepoPg)(nil)

func (r *RepoPg) Create(ctx context.Context, trip *domain.Trip) error {
	e := tripToEntity(trip)

	const q = `
		INSERT INTO trip_domain.trips
			(id, client_id, from_point, to_point, departure_time, seats, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.pool.Exec(ctx, q,
		e.ID, e.ClientID, e.FromPoint, e.ToPoint,
		e.DepartureTime, e.AvailableSeats, e.Status, e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert trip: %w", err)
	}
	return nil
}

func (r *RepoPg) GetByID(ctx context.Context, id uuid.UUID) (*domain.Trip, error) {
	const q = `
		SELECT id, client_id, from_point, to_point, departure_time, seats, status, created_at
		FROM trip_domain.trips
		WHERE id = $1`

	var e entity.Trip
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&e.ID, &e.ClientID, &e.FromPoint, &e.ToPoint,
		&e.DepartureTime, &e.AvailableSeats, &e.Status, &e.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not found — не ошибка на уровне БД
		}
		return nil, fmt.Errorf("get trip: %w", err)
	}

	return tripFromEntity(e), nil
}

func (r *RepoPg) GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*domain.Trip, error) {
	const q = `
		SELECT id, client_id, from_point, to_point, departure_time, seats, status, created_at
		FROM trip_domain.trips
		WHERE id = $1
		FOR UPDATE`

	var e entity.Trip
	err := tx.QueryRow(ctx, q, id).Scan(
		&e.ID, &e.ClientID, &e.FromPoint, &e.ToPoint,
		&e.DepartureTime, &e.AvailableSeats, &e.Status, &e.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: trip %s", domain.ErrTripNotFound, id)
		}
		return nil, fmt.Errorf("get trip for update: %w", err)
	}
	return tripFromEntity(e), nil
}

func (r *RepoPg) UpdateStatus(ctx context.Context, tx pgx.Tx, id uuid.UUID, status domain.TripStatus) error {
	const q = `UPDATE trip_domain.trips SET status = $2 WHERE id = $1`
	if _, err := tx.Exec(ctx, q, id, status); err != nil {
		return fmt.Errorf("update trip status: %w", err)
	}
	return nil
}

func (r *RepoPg) InsertHistory(ctx context.Context, tx pgx.Tx, tripID uuid.UUID, from, to domain.TripStatus) error {
	const q = `
		INSERT INTO trip_domain.trip_history (id, trip_id, from_status, to_status)
		VALUES ($1, $2, $3, $4)`
	if _, err := tx.Exec(ctx, q, uuid.New(), tripID, from, to); err != nil {
		return fmt.Errorf("insert trip history: %w", err)
	}
	return nil
}
