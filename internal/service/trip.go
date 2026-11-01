package service

import "github.com/kostas-111/share_trip/internal/usecase"

// TripService — оркестратор сценариев: принимает модели своего слоя,
// вызывает сервисы usecase-слоя и собирает ответ.
type TripService struct {
	creator   *usecase.TripCreator
	reader    *usecase.TripReader
	publisher *usecase.TripPublisher
}

func NewTripService(
	creator *usecase.TripCreator,
	reader *usecase.TripReader,
	publisher *usecase.TripPublisher,
) *TripService {
	return &TripService{
		creator:   creator,
		reader:    reader,
		publisher: publisher,
	}
}
