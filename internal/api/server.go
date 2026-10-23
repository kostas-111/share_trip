package api

import (
	"context"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type Server struct {
	db   Pinger
	trip *TripHandler
}

func NewServer(db Pinger, trip *TripHandler) *Server {
	return &Server{
		db:   db,
		trip: trip,
	}
}
