package service

import (
	"github.com/kostas-111/share_trip/internal/repository"
)

type TripService struct {
	repo   repository.TripRepository
	runner repository.Runner
}

func NewTripService(repo repository.TripRepository, runner repository.Runner) *TripService {
	return &TripService{repo: repo, runner: runner}
}
