package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

		// timestamptz хранит микросекунды: без усечения значение из ответа create
		// (наносекунды) не совпадёт с тем, что вернёт БД при публикации
		departure := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Microsecond)
		created := createTrip(t, departure)

		status, body := publishTrip(t, created.ID)
		require.Equal(t, http.StatusOK, status)

		var published api.PublishTripResponse
		require.NoError(t, json.Unmarshal(body, &published))

		// created_at заполняет сервер, а время из БД приходит в локальной таймзоне:
		// приводим к UTC и обнуляем серверное поле, чтобы сравнить ответ целиком
		publishedCreatedAt := published.CreatedAt
		published.DepartureTime = published.DepartureTime.UTC()
		published.CreatedAt = time.Time{}

		expectedPublished := api.PublishTripResponse{
			ID:             created.ID,
			ClientID:       created.ClientID,
			FromPoint:      "Москва",
			ToPoint:        "Тверь",
			DepartureTime:  departure,
			AvailableSeats: 3,
			Status:         "published",
		}
		require.Equal(t, expectedPublished, published)
		require.False(t, publishedCreatedAt.IsZero())

		// отдельный GET проверяет, что статус лег в БД, а не только в ответ публикации
		status, body = doJSON(t, http.MethodGet, "/trip/"+created.ID.String(), nil)
		require.Equal(t, http.StatusOK, status)

		var stored api.GetTripResponse
		require.NoError(t, json.Unmarshal(body, &stored))

		storedCreatedAt := stored.CreatedAt
		stored.DepartureTime = stored.DepartureTime.UTC()
		stored.CreatedAt = time.Time{}

		expectedStored := api.GetTripResponse{
			ID:             created.ID,
			ClientID:       created.ClientID,
			FromPoint:      "Москва",
			ToPoint:        "Тверь",
			DepartureTime:  departure,
			AvailableSeats: 3,
			Status:         "published",
		}
		require.Equal(t, expectedStored, stored)
		require.False(t, storedCreatedAt.IsZero())
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
		require.Equal(t, api.ErrorResponse{
			Code: api.CodeConflict,
			Message: "invalid trip status transition: only draft trips can be published " +
				"(current status: published)",
		}, errResp)
	})

	t.Run("not found - несуществующая поездка", func(t *testing.T) {
		t.Parallel()

		tripID := uuid.New()

		status, body := publishTrip(t, tripID)
		require.Equal(t, http.StatusNotFound, status)

		var errResp api.ErrorResponse
		require.NoError(t, json.Unmarshal(body, &errResp))
		require.Equal(t, api.ErrorResponse{
			Code:    api.CodeNotFound,
			Message: fmt.Sprintf("trip not found: trip %s", tripID),
		}, errResp)
	})
}
