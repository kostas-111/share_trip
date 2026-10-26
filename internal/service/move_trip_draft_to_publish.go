package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kostas-111/share_trip/internal/domain"
)

type PublishTripRequest struct {
	TripID uuid.UUID
}

type PublishTripResponse struct {
	Trip *domain.Trip
}

func (s *TripService) PublishTrip(ctx context.Context, req PublishTripRequest) (*PublishTripResponse, error) {
	if req.TripID == uuid.Nil {
		return nil, fmt.Errorf("%w: trip_id is required", domain.ErrValidation)
	}

  var resp *PublishTripResponse
	err := s.runner.WithTx(ctx, func(tx pgx.Tx) error {
		// 1) Читаем с блокировкой.
		trip, err := s.repo.GetByIDForUpdate(ctx, tx, req.TripID)
		if err != nil {
			return err // ошибка ErrNotFound уже обёрнута в repo
		}

		// 2) Бизнес-правило перехода (чистая функция домена).
		domResp, err := domain.PublishTrip(trip)
		if err != nil {
			return err // ErrConflict / ErrValidation
		}

		// 3) Обновляем статус и пишем историю в этой же транзакции.
		if err := s.repo.UpdateStatus(ctx, tx, trip.ID, domResp.ToStatus); err != nil {
			return err
		}
		if err := s.repo.InsertHistory(ctx, tx, trip.ID, domResp.FromStatus, domResp.ToStatus); err != nil {
			return err
		}

    domResp.Trip.Status = domResp.ToStatus
    resp = &PublishTripResponse{Trip: domResp.Trip}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return resp, nil
}
