package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type GetTripRequest struct {
	ID uuid.UUID
}

type GetTripResponse struct {
	ID uuid.UUID
}

func GetTrip(req GetTripRequest) (GetTripResponse, error) {
	if req.ID == uuid.Nil {
		return GetTripResponse{}, fmt.Errorf("%w: trip_id is required", ErrValidation)
	}
	return GetTripResponse{ID: req.ID}, nil
}
