package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	srv := New(":8080")

	t.Run("status endpoint", func(t *testing.T) {
		// Set some test metrics
		srv.IncrementBotRequests()
		srv.IncrementBotRequests()
		srv.SetBotConnected(true)
		srv.IncrementAMLRequests()
		srv.SetAMLConnected(true)

		req := httptest.NewRequest(http.MethodGet, "/status", nil)
		w := httptest.NewRecorder()

		srv.handleStatus(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response StatusResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "ok", response.Status)
		assert.True(t, time.Since(response.Timestamp) < time.Second)
		assert.Equal(t, int64(2), response.Bot.RequestsCount)
		assert.True(t, response.Bot.IsConnected)
		assert.Equal(t, int64(1), response.AML.RequestsCount)
		assert.True(t, response.AML.IsConnected)
	})

	t.Run("invalid method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/status", nil)
		w := httptest.NewRecorder()

		srv.handleStatus(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func TestServerStartStop(t *testing.T) {
	srv := New(":0") // Use port 0 to get a random available port

	// Start server in background
	errChan := make(chan error)
	go func() {
		errChan <- srv.Start()
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Stop server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Stop(ctx)
	assert.NoError(t, err)

	// Check if server stopped
	select {
	case err := <-errChan:
		assert.Error(t, err, "http: Server closed")
	case <-time.After(1 * time.Second):
		t.Fatal("server did not stop in time")
	}
}
