package api

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type Server struct {
	db       Pinger
	trip     *TripHandler
	registry prometheus.Gatherer
}

func NewServer(db Pinger, trip *TripHandler, reg prometheus.Gatherer) *Server {
	return &Server{
		db:       db,
		trip:     trip,
		registry: reg,
	}
}
