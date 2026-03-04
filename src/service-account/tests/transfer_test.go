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

func InitialSenderBalance() decimal.Decimal {
	return decimal.NewFromInt(50000)
}

func TestTransfer_TableDriven(t *testing.T) {
	router, mockPub := setupTestEnv(testDBPool)

	senderID := createTestUser(t, testDBPool)
	receiverID := createTestUser(t, testDBPool)

	senderAccountID := createTestAccount(t, testDBPool, senderID, "RUB")
	receiverAccountID := createTestAccount(t, testDBPool, receiverID, "RUB")

	tests := []struct {
		name                string
		senderPathID        string
		payload             map[string]interface{}
		expectedCode        int
		expectedEventsCount int
		expectedSenderBal   decimal.Decimal
		expectedReceiverBal decimal.Decimal
	}{
		{
			name:         "Success - Valid Transfer",
			senderPathID: senderAccountID.String(),
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "15000",
				"currency_code": "RUB",
				"description":   "string",
			},
			expectedCode:        http.StatusOK,
			expectedEventsCount: 1,
			expectedSenderBal:   decimal.NewFromInt(35000),
			expectedReceiverBal: decimal.NewFromInt(15000),
		},
		{
			name:         "Fail - Insufficient Funds (Not enough money)",
			senderPathID: senderAccountID.String(),
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "999000",
				"currency_code": "RUB",
				"description":   "string",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:         "Fail - Negative Amount",
			senderPathID: senderAccountID.String(),
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "-5000",
				"currency_code": "RUB",
				"description":   "string",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:         "Fail - Unknown Receiver",
			senderPathID: senderAccountID.String(),
			payload: map[string]interface{}{
				"to_account_id": uuid.New().String(),
				"amount":        "10000",
				"currency_code": "RUB",
				"description":   "string",
			},
			expectedCode:        http.StatusNotFound,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPub.Clear()

			setAccountBalance(t, testDBPool, senderAccountID, InitialSenderBalance())
			setAccountBalance(t, testDBPool, receiverAccountID, decimal.Zero)

			bodyBytes, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/accounts/%s/transfer", tt.senderPathID)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Idempotency-Key", uuid.New().String())

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "Response body: %s", rr.Body.String())

			require.Len(t, mockPub.Events, tt.expectedEventsCount)

			actualSenderBal := getAccountBalance(t, testDBPool, senderAccountID)
			actualReceiverBal := getAccountBalance(t, testDBPool, receiverAccountID)

			require.True(t, tt.expectedSenderBal.Equal(actualSenderBal),
				"Sender balance mismatch: expected %s, got %s", tt.expectedSenderBal.String(), actualSenderBal.String())

			require.True(t, tt.expectedReceiverBal.Equal(actualReceiverBal),
				"Receiver balance mismatch: expected %s, got %s", tt.expectedReceiverBal.String(), actualReceiverBal.String())
		})
	}
}
