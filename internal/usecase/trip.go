package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TripStatus string

// Значения совпадают с enum trip_domain.trip_status в миграциях.
const (
	StatusDraft     TripStatus = "draft"
	StatusPublished TripStatus = "published"
	StatusCanceled  TripStatus = "canceled"
	StatusCompleted TripStatus = "completed"
)

type Trip struct {
	ID             uuid.UUID
	ClientID       uuid.UUID
	FromPoint      string
	ToPoint        string
	DepartureTime  time.Time
	AvailableSeats int
	Status         TripStatus
	CreatedAt      time.Time
}

// Порт объявлен в слое, который им пользуется; реализация — *repository.RepoPg.
// Отсюда направление зависимости repository -> usecase.
type TripRepository interface {
	Create(ctx context.Context, trip *Trip) error
	GetByID(ctx context.Context, id uuid.UUID) (*Trip, error)
	GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Trip, error)
	UpdateStatus(ctx context.Context, tx pgx.Tx, id uuid.UUID, status TripStatus) error
	InsertHistory(ctx context.Context, tx pgx.Tx, tripID uuid.UUID, from, to TripStatus) error
}
