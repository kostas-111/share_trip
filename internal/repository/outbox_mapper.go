package repository

import (
	"encoding/json"
	"fmt"

	"github.com/kostas-111/share_trip/internal/domain"
	"github.com/kostas-111/share_trip/internal/repository/entity"
)

func outboxToEntity(e domain.OutboxEvent) (entity.OutboxEvent, error) {
	payload, err := json.Marshal(e.Payload)
	if err != nil {
		return entity.OutboxEvent{}, fmt.Errorf("marshal outbox payload: %w", err)
	}

	return entity.OutboxEvent{
		ID:          e.ID,
		EventName:   e.EventName,
		AggregateID: e.AggregateID,
		Payload:     payload,
	}, nil
}
