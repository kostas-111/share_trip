package api

import "github.com/gofiber/fiber/v2"

func (s *Server) Route(route fiber.Router) {
	route.Get("/check-health", s.CheckDbConn)
	route.Post("/trip/create", s.trip.CreateTrip)
	route.Post("/trip/publish", s.trip.PublishTrip)
	route.Get("/trip/:tripId", s.trip.GetTrip)
}
