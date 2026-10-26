package domain

import (
	"fmt"
	"time"
)

type PublishTripResponse struct {
	Trip       *Trip
	FromStatus TripStatus
	ToStatus   TripStatus
}

// PublishTrip — только бизнес-правила перехода DRAFT -> PUBLISHED.
// Не делает I/O. Ожидает, что trip уже прочитан (под блокировкой) в repo.
func PublishTrip(trip *Trip) (PublishTripResponse, error) {
	if trip == nil {
		return PublishTripResponse{}, fmt.Errorf("%w: trip not found", ErrTripNotFound)
	}
	if trip.Status != StatusDraft {
		return PublishTripResponse{}, fmt.Errorf(
			"%w: only draft trips can be published (current status: %s)",
			ErrInvalidStatus, trip.Status,
		)
	}
	if !trip.DepartureTime.After(time.Now()) {
		return PublishTripResponse{}, fmt.Errorf(
			"%w: cannot publish a trip with departure_time in the past",
			ErrValidation,
		)
	}

	return PublishTripResponse{
		Trip:       trip,
		FromStatus: trip.Status,
		ToStatus:   StatusPublished,
	}, nil
}
