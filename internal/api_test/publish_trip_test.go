package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/kostas-111/share_trip/internal/api"
)

func doJSON(t *testing.T, method, path string, payload any) (int, []byte) {
	t.Helper()

	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		require.NoError(t, err)
		body = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, path, body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := testApp.Test(req, -1)
	require.NoError(t, err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("close response body: %v", err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, respBody
}

// createTrip создаёт поездку в состоянии draft и возвращает ответ api.
func createTrip(t *testing.T, departureTime time.Time) api.CreateTripResponse {
	t.Helper()

	status, body := doJSON(t, http.MethodPost, "/trip/create", api.CreateTripRequest{
		ClientID:       uuid.New(),
		FromPoint:      "Москва",
		ToPoint:        "Тверь",
		DepartureTime:  departureTime,
		AvailableSeats: 3,
	})
	require.Equal(t, http.StatusCreated, status)

	var created api.CreateTripResponse
	require.NoError(t, json.Unmarshal(body, &created))
	return created
}

func publishTrip(t *testing.T, tripID uuid.UUID) (int, []byte) {
	t.Helper()
	return doJSON(t, http.MethodPost, "/trip/publish", api.PublishTripRequest{TripID: tripID})
}

func TestServer_PublishTrip(t *testing.T) {
	t.Parallel()

	t.Run("success - переход draft в published", func(t *testing.T) {
		t.Parallel()

		created := createTrip(t, time.Now().UTC().Add(24*time.Hour))

		status, body := publishTrip(t, created.ID)
		require.Equal(t, http.StatusOK, status)

		var published api.PublishTripResponse
		require.NoError(t, json.Unmarshal(body, &published))

		require.Equal(t, created.ID, published.ID)
		require.Equal(t, created.ClientID, published.ClientID)
		require.Equal(t, created.FromPoint, published.FromPoint)
		require.Equal(t, created.ToPoint, published.ToPoint)
		require.WithinDuration(t, created.DepartureTime, published.DepartureTime, time.Microsecond)
		require.Equal(t, created.AvailableSeats, published.AvailableSeats)
		require.Equal(t, "published", published.Status)

		// отдельный GET проверяет, что статус лег в БД, а не только в ответ публикации
		status, body = doJSON(t, http.MethodGet, "/trip/"+created.ID.String(), nil)
		require.Equal(t, http.StatusOK, status)

		var stored api.GetTripResponse
		require.NoError(t, json.Unmarshal(body, &stored))
		require.Equal(t, created.ID, stored.ID)
		require.Equal(t, "published", stored.Status)
	})

	t.Run("success - событие trip_published в outbox", func(t *testing.T) {
		t.Parallel()

		created := createTrip(t, time.Now().UTC().Add(48*time.Hour))

		status, _ := publishTrip(t, created.ID)
		require.Equal(t, http.StatusOK, status)

		var (
			eventName     string
			tripIDPayload string
		)
		err := testPool.QueryRow(testCtx, `
			SELECT event_name, payload->>'trip_id'
			FROM trip_domain.outbox_event
			WHERE aggregate_id = $1`, created.ID,
		).Scan(&eventName, &tripIDPayload)
		require.NoError(t, err)

		require.Equal(t, "trip_published", eventName)
		require.Equal(t, created.ID.String(), tripIDPayload)
	})

	t.Run("conflict - повторная публикация", func(t *testing.T) {
		t.Parallel()

		created := createTrip(t, time.Now().UTC().Add(24*time.Hour))

		status, _ := publishTrip(t, created.ID)
		require.Equal(t, http.StatusOK, status)

		status, body := publishTrip(t, created.ID)
		require.Equal(t, http.StatusConflict, status)

		var errResp api.ErrorResponse
		require.NoError(t, json.Unmarshal(body, &errResp))
		require.Equal(t, api.CodeConflict, errResp.Code)
		require.Contains(t, errResp.Message, "only draft trips can be published")
	})

	t.Run("not found - несуществующая поездка", func(t *testing.T) {
		t.Parallel()

		status, body := publishTrip(t, uuid.New())
		require.Equal(t, http.StatusNotFound, status)

		var errResp api.ErrorResponse
		require.NoError(t, json.Unmarshal(body, &errResp))
		require.Equal(t, api.CodeNotFound, errResp.Code)
	})
}
