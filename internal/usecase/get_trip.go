package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type TripReader struct {
	repo TripRepository
}

func NewTripReader(repo TripRepository) *TripReader {
	return &TripReader{repo: repo}
}

func (r *TripReader) GetByID(ctx context.Context, id uuid.UUID) (*Trip, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: trip_id is required", ErrValidation)
	}

	trip, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("TripReader.GetByID: %w", err)
	}

	// repository не считает отсутствие строки ошибкой - решение принимает сценарий
	if trip == nil {
		return nil, fmt.Errorf("%w: id=%s", ErrTripNotFound, id)
	}

	return trip, nil
}
