package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/kostas-111/share_trip/internal/service"
)

type GetTripResponse struct {
	ID             uuid.UUID `json:"id"`
	ClientID       uuid.UUID `json:"client_id"`
	FromPoint      string    `json:"from_point"`
	ToPoint        string    `json:"to_point"`
	DepartureTime  time.Time `json:"departure_time"`
	AvailableSeats int       `json:"available_seats"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

func (h *TripHandler) GetTrip(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("tripId"))
	if err != nil {
		return writeBadRequest(c, "invalid trip id: "+err.Error())
	}

	resp, err := h.svc.GetTrip(c.UserContext(), service.GetTripRequest{ID: id})
	if err != nil {
		return writeError(c, err)
	}

	t := resp.Trip
	return c.Status(fiber.StatusOK).JSON(GetTripResponse{
		ID:             t.ID,
		ClientID:       t.ClientID,
		FromPoint:      t.FromPoint,
		ToPoint:        t.ToPoint,
		DepartureTime:  t.DepartureTime,
		AvailableSeats: t.AvailableSeats,
		Status:         string(t.Status),
		CreatedAt:      t.CreatedAt,
	})
}
