package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const headerClientID = "X-Client-Id"

func clientIDFromRequest(c *fiber.Ctx) (uuid.UUID, error) {
	raw := c.Get(headerClientID)
	if raw == "" {
		return uuid.Nil, fmt.Errorf("%s is required", headerClientID)
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s is not a valid uuid", headerClientID)
	}

	return id, nil
}
