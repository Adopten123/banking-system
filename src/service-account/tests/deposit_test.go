package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func InitialBalance() decimal.Decimal {
	return decimal.NewFromInt(10000)
}

func TestDeposit_TableDriven(t *testing.T) {
	router, mockPub := setupTestEnv(testDBPool)

	validUserID := createTestUser(t, testDBPool)
	validAccountID := createTestAccount(t, testDBPool, validUserID, "RUB")

	tests := []struct {
		name                string
		accountPathID       string
		payload             map[string]interface{}
		expectedCode        int
		expectedEventsCount int
		expectedBalance     decimal.Decimal
	}{
		{
			name:          "Success - Valid Deposit",
			accountPathID: validAccountID.String(),
			payload: map[string]interface{}{
				"amount": "5000",
			},
			expectedCode:        http.StatusCreated,
			expectedEventsCount: 1,
			expectedBalance:     decimal.NewFromInt(15000),
		},
		{
			name:          "Fail - Negative Amount",
			accountPathID: validAccountID.String(),
			payload: map[string]interface{}{
				"amount": "-1000",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedBalance:     InitialBalance(),
		},
		{
			name:          "Fail - Zero Amount",
			accountPathID: validAccountID.String(),
			payload: map[string]interface{}{
				"amount": "0",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedBalance:     InitialBalance(),
		},
		{
			name:          "Fail - Unknown Account",
			accountPathID: uuid.New().String(),
			payload: map[string]interface{}{
				"amount": "5000",
			},
			expectedCode:        http.StatusNotFound,
			expectedEventsCount: 0,
			expectedBalance:     InitialBalance(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPub.Clear()

			setAccountBalance(t, testDBPool, validAccountID, InitialBalance())

			bodyBytes, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/accounts/%s/deposit", tt.accountPathID)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Idempotency-Key", uuid.New().String())

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "Response body: %s", rr.Body.String())

			require.Len(t, mockPub.Events, tt.expectedEventsCount)

			if tt.accountPathID == validAccountID.String() {
				actualBalance := getAccountBalance(t, testDBPool, validAccountID)
				require.True(t, tt.expectedBalance.Equal(actualBalance),
					"Balance mismatch: expected %s, got %s", tt.expectedBalance.String(), actualBalance.String())
			}
		})
	}
}
