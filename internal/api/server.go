package api

import "github.com/kostas-111/share_trip/internal/repository"

type Server struct {
	Repository *repository.RepoPg
}

func NewServer(repo *repository.RepoPg) *Server {
	return &Server{Repository: repo}
}
