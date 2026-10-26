package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kostas-111/share_trip/internal/service"
)

type PublishTripRequest struct {
	TripID uuid.UUID `json:"trip_id"`
}

type PublishTripResponse struct {
	ID             uuid.UUID `json:"id"`
	ClientID       uuid.UUID `json:"client_id"`
	FromPoint      string    `json:"from_point"`
	ToPoint        string    `json:"to_point"`
	DepartureTime  time.Time `json:"departure_time"`
	AvailableSeats int       `json:"available_seats"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

func (h *TripHandler) PublishTrip(c *fiber.Ctx) error {
	var req PublishTripRequest
	if err := c.BodyParser(&req); err != nil {
		return writeBadRequest(c, "invalid JSON body: "+err.Error())
	}

	resp, err := h.svc.PublishTrip(c.UserContext(), service.PublishTripRequest{
		TripID: req.TripID,
	})
	if err != nil {
		return writeError(c, err)
	}

	t := resp.Trip
	return c.Status(fiber.StatusOK).JSON(PublishTripResponse{
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
