package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) Route(route fiber.Router) {
	route.Get("/check-health", s.CheckDbConn)
	route.Post("/trip/create", s.trip.CreateTrip)
	route.Post("/trip/publish", s.trip.PublishTrip)
	route.Get("/trip/:tripId", s.trip.GetTrip)
	route.Get("/metrics", adaptor.HTTPHandler(promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{})))
}
