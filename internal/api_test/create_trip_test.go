package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kostas-111/share_trip/internal/api"
	"github.com/stretchr/testify/require"
)

func TestServer_CreateTrip(t *testing.T) {
	t.Parallel()

	t.Run("success - создание поездки", func(t *testing.T) {
		payload := api.CreateTripRequest{
			ClientID:       uuid.New(),
			FromPoint:      "Москва",
			ToPoint:        "Тверь",
			DepartureTime:  time.Now().UTC().Add(24 * time.Hour),
			AvailableSeats: 3,
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPost,
			"/trip/create",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got api.CreateTripResponse
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)

		// id и created_at генерирует сервер, тест их заранее не знает:
		// значения выносятся из сравнения структур и проверяются отдельно
		id, createdAt := got.ID, got.CreatedAt
		got.ID, got.CreatedAt = uuid.Nil, time.Time{}

		expected := api.CreateTripResponse{
			ClientID:       payload.ClientID,
			FromPoint:      payload.FromPoint,
			ToPoint:        payload.ToPoint,
			DepartureTime:  payload.DepartureTime,
			AvailableSeats: payload.AvailableSeats,
			Status:         "draft",
		}
		require.Equal(t, expected, got)

		require.NotEqual(t, uuid.Nil, id)
		require.False(t, createdAt.IsZero())
	})
}
