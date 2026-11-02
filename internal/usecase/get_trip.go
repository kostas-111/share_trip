package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kostas-111/share_trip/internal/domain"
	"github.com/kostas-111/share_trip/internal/repository"
)

type TripReader struct {
	repo repository.TripRepository
}

func NewTripReader(repo repository.TripRepository) *TripReader {
	return &TripReader{repo: repo}
}

func (r *TripReader) GetByID(ctx context.Context, id uuid.UUID) (*domain.Trip, error) {
	if err := (domain.GetTripParams{ID: id}).Validate(); err != nil {
		return nil, err
	}

	trip, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("TripReader.GetByID: %w", err)
	}

	// repository не считает отсутствие строки ошибкой — решение принимает сценарий
	if trip == nil {
		return nil, fmt.Errorf("%w: id=%s", domain.ErrTripNotFound, id)
	}

	return trip, nil
}
