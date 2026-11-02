package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type TripHistory struct {
	ID         uuid.UUID      `db:"id"`
	TripID     uuid.UUID      `db:"trip_id"`
	FromStatus sql.NullString `db:"from_status"`
	ToStatus   string         `db:"to_status"`
	CreatedAt  time.Time      `db:"created_at"`
}
