package api

import "github.com/gofiber/fiber/v2"

func (s *Server) Route(route fiber.Router) {
	route.Get("/check-health", s.CheckDbConn)
}
