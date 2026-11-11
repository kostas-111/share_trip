package repository

import (
	"github.com/kostas-111/share_trip/internal/repository/entity"
	"github.com/kostas-111/share_trip/internal/usecase"
)

func tripToEntity(t *usecase.Trip) entity.Trip {
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

func tripFromEntity(e entity.Trip) *usecase.Trip {
	return &usecase.Trip{
		ID:             e.ID,
		ClientID:       e.ClientID,
		FromPoint:      e.FromPoint,
		ToPoint:        e.ToPoint,
		DepartureTime:  e.DepartureTime,
		AvailableSeats: e.AvailableSeats,
		Status:         usecase.TripStatus(e.Status),
		CreatedAt:      e.CreatedAt,
	}
}
