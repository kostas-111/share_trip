package repository

import (
	"github.com/kostas-111/share_trip/internal/domain"
	"github.com/kostas-111/share_trip/internal/repository/entity"
)

func tripToEntity(t *domain.Trip) entity.Trip {
	return entity.Trip{
		ID:             t.ID,
		ClientID:       t.ClientID,
		FromPoint:      t.FromPoint,
		ToPoint:        t.ToPoint,
		DepartureTime:  t.DepartureTime,
		AvailableSeats: t.AvailableSeats,
		Status:         string(t.Status),
		CreatedAt:      t.CreatedAt,
	}
}

func tripFromEntity(e entity.Trip) *domain.Trip {
	return &domain.Trip{
		ID:             e.ID,
		ClientID:       e.ClientID,
		FromPoint:      e.FromPoint,
		ToPoint:        e.ToPoint,
		DepartureTime:  e.DepartureTime,
		AvailableSeats: e.AvailableSeats,
		Status:         domain.TripStatus(e.Status),
		CreatedAt:      e.CreatedAt,
	}
}
