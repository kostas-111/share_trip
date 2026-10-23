package service

import "github.com/kostas-111/share_trip/internal/repository"

type TripService struct {
	repo repository.TripRepository
}

func NewTripService(repo repository.TripRepository) *TripService {
	return &TripService{repo: repo}
}
