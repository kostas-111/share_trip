package api

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/kostas-111/share_trip/internal/observability/logctx"
	"github.com/kostas-111/share_trip/internal/usecase"
)

type ErrorCode string

const (
	CodeValidation    ErrorCode = "VALIDATION_ERROR"
	CodeNotFound      ErrorCode = "NOT_FOUND"
	CodeForbidden     ErrorCode = "FORBIDDEN"
	CodeConflict      ErrorCode = "CONFLICT"
	CodeUnprocessable ErrorCode = "UNPROCESSABLE"
	CodeInternal      ErrorCode = "INTERNAL_ERROR"
)

type ErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func errorStatus(err error) (int, ErrorCode, string) {
	switch {
	case errors.Is(err, usecase.ErrValidation):
		return fiber.StatusBadRequest, CodeValidation, err.Error()
	case errors.Is(err, usecase.ErrTripNotFound):
		return fiber.StatusNotFound, CodeNotFound, err.Error()
	case errors.Is(err, usecase.ErrForbidden):
		return fiber.StatusForbidden, CodeForbidden, err.Error()
	case errors.Is(err, usecase.ErrTripCanceled),
		errors.Is(err, usecase.ErrNoAvailableSeats),
		errors.Is(err, usecase.ErrInvalidStatus):
		return fiber.StatusConflict, CodeConflict, err.Error()
	case errors.Is(err, usecase.ErrUnprocessable):
		return fiber.StatusUnprocessableEntity, CodeUnprocessable, err.Error()
	default:
		return fiber.StatusInternalServerError, CodeInternal, "internal server error"
	}
}

// отправляет ErrorResponse клиенту.
func writeError(c *fiber.Ctx, err error) error {
	status, code, msg := errorStatus(err)
	if status == fiber.StatusInternalServerError {
		// клиенту уходит обезличенное сообщение: детали ошибки остаются только в логе
		logctx.Logger(c.UserContext()).Error(
			"unhandled error",
			slog.Any("error", err),
		)
	}
	return c.Status(status).JSON(ErrorResponse{Code: code, Message: msg})
}

// для ошибок парсинга/валидации тела на уровне api.
func writeBadRequest(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		Code:    CodeValidation,
		Message: msg,
	})
}
