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
