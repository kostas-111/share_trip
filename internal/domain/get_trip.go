package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type GetTripParams struct {
	ID uuid.UUID
}

func (r GetTripParams) Validate() error {
	if r.ID == uuid.Nil {
		return fmt.Errorf("%w: trip_id is required", ErrValidation)
	}
	return nil
}
