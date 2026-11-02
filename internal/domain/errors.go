package domain

import "errors"

var (
	// 400 - вход неверный
	ErrValidation = errors.New("validation error")

	// 404 - объект не найден
	ErrTripNotFound = errors.New("trip not found")

	// 403 - действие запрещено
	ErrForbidden = errors.New("forbidden")

	// 409 - конфликт состояния
	ErrTripCanceled     = errors.New("trip is canceled")
	ErrNoAvailableSeats = errors.New("no available seats")
	ErrInvalidStatus    = errors.New("invalid trip status transition")

	// 422 - бизнес-правило нарушено при валидных данных
	ErrUnprocessable = errors.New("unprocessable entity")
)
