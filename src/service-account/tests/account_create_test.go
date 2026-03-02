package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount_Success(t *testing.T) {
	router, mockPub := setupTestEnv(testDBPool)
	mockPub.Clear()

	// Make user
	userID := createTestUser(t, testDBPool)

	// Make query-body
	reqBody := map[string]interface{}{
		"user_id":       userID.String(),
		"currency_code": "USD",
		"type_id":       1,
		"name":          "My Test Account",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(http.MethodPost, "/api/accounts", bytes.NewReader(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Waiting success status
	require.Equal(t, http.StatusCreated, rr.Code)

	// Check response
	require.Contains(t, rr.Body.String(), userID.String())
	require.Contains(t, rr.Body.String(), "USD")

	// Check RabbitMQ
	require.Len(t, mockPub.Events, 1, "Expected exactly 1 event published to RabbitMQ")
	require.Equal(t, "AccountCreatedEvent", mockPub.Events[0].EventName())
}
