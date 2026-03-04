package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount_TableDriven(t *testing.T) {
	// 1. Make test env
	router, mockPub := setupTestEnv(testDBPool)

	// 2. Make user
	validUserID := createTestUser(t, testDBPool)

	// 3. Make test table
	tests := []struct {
		name                string
		payload             map[string]interface{}
		expectedCode        int
		expectedEventsCount int
	}{
		{
			name: "Success - Valid Account",
			payload: map[string]interface{}{
				"user_id":       validUserID.String(),
				"currency_code": "USD",
				"type_id":       1,
				"name":          "My Valid Account",
			},
			expectedCode:        http.StatusCreated,
			expectedEventsCount: 1,
		},
		{
			name: "Fail - Missing User ID",
			payload: map[string]interface{}{
				"user_id":       " ",
				"currency_code": "EUR",
				"type_id":       1,
				"name":          "No User Account",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
		},
		{
			name: "Fail - Invalid UUID Format",
			payload: map[string]interface{}{
				"user_id":       "invalid-uuid-string",
				"currency_code": "RUB",
				"type_id":       1,
				"name":          "Bad UUID",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
		},
		{
			name: "Fail - Unknown Currency",
			payload: map[string]interface{}{
				"user_id":       validUserID.String(),
				"currency_code": "XYZ",
				"type_id":       1,
				"name":          "Wrong Currency",
			},
			expectedCode:        http.StatusInternalServerError,
			expectedEventsCount: 0,
		},
		{
			name: "Fail - User Does Not Exist in DB",
			payload: map[string]interface{}{
				"user_id":       uuid.New().String(),
				"currency_code": "USD",
				"type_id":       1,
				"name":          "Ghost User",
			},
			expectedCode:        http.StatusInternalServerError,
			expectedEventsCount: 0,
		},
	}

	// 4. Testing loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear mock before every sub-test
			mockPub.Clear()

			bodyBytes, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/accounts", bytes.NewReader(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			// Exec HTTP-query
			router.ServeHTTP(rr, req)

			// Check status code
			require.Equal(t, tt.expectedCode, rr.Code, "Response body: %s", rr.Body.String())

			// Check event in RabbitMQ
			require.Len(t, mockPub.Events, tt.expectedEventsCount)

			if tt.expectedEventsCount > 0 {
				require.Equal(t, "AccountCreatedEvent", mockPub.Events[0].EventName())
			}
		})
	}
}
