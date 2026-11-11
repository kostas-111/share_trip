package api

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kostas-111/share_trip/internal/observability/logctx"
	"github.com/kostas-111/share_trip/internal/service"
)

type CreateTripRequest struct {
	ClientID       uuid.UUID `json:"client_id"`
	FromPoint      string    `json:"from_point"`
	ToPoint        string    `json:"to_point"`
	DepartureTime  time.Time `json:"departure_time"`
	AvailableSeats int       `json:"available_seats"`
}

type CreateTripResponse struct {
	ID             uuid.UUID `json:"id"`
	ClientID       uuid.UUID `json:"client_id"`
	FromPoint      string    `json:"from_point"`
	ToPoint        string    `json:"to_point"`
	DepartureTime  time.Time `json:"departure_time"`
	AvailableSeats int       `json:"available_seats"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

func (h *TripHandler) CreateTrip(c *fiber.Ctx) error {
	ctx := c.UserContext()

	logger := logctx.Logger(ctx).With(
		slog.String("server", "TripServer"),
		slog.String("handler", "CreateTrip"),
	)

	var req CreateTripRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warn(
			"create trip failed: invalid JSON body",
			slog.Any("error", err),
		)
		return writeBadRequest(c, "invalid JSON body: "+err.Error())
	}

	logger = logger.With(slog.String("client_id", req.ClientID.String()))
	ctx = logctx.WithLogger(ctx, logger)

	logger.Info("create trip request accepted")

	resp, err := h.svc.CreateTrip(ctx, service.CreateTripRequest{
		ClientID:       req.ClientID,
		FromPoint:      req.FromPoint,
		ToPoint:        req.ToPoint,
		DepartureTime:  req.DepartureTime,
		AvailableSeats: req.AvailableSeats,
	})
	if err != nil {
		return writeError(c, err)
	}

	t := resp.Trip

	logger.Info(
		"create trip completed",
		slog.String("trip_id", t.ID.String()),
	)

	return c.Status(fiber.StatusCreated).JSON(CreateTripResponse{
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
