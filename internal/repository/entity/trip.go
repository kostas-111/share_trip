package entity

import (
	"time"

	"github.com/google/uuid"
)

type Trip struct {
	ID             uuid.UUID `db:"id"`
	ClientID       uuid.UUID `db:"client_id"`
	FromPoint      string    `db:"from_point"`
	ToPoint        string    `db:"to_point"`
	DepartureTime  time.Time `db:"departure_time"`
	AvailableSeats int       `db:"seats"`
	Status         string    `db:"status"`
	CreatedAt      time.Time `db:"created_at"`
}
