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
	validUserID := createTestUser(t, testDBPool)

	senderAccountID := createTestAccount(t, testDBPool, senderID, "RUB")
	receiverAccountID := createTestAccount(t, testDBPool, receiverID, "RUB")

	usdAccountID := createTestAccount(t, testDBPool, validUserID, "USD")
	eurAccountID := createTestAccount(t, testDBPool, validUserID, "EUR")

	blockedAccountID := createTestAccount(t, testDBPool, validUserID, "RUB")
	frozenAccountID := createTestAccount(t, testDBPool, validUserID, "RUB")

	const blockedStatus = 3
	updateAccountStatus(t, testDBPool, blockedAccountID, blockedStatus)
	const frozenStatus = 2
	updateAccountStatus(t, testDBPool, frozenAccountID, frozenStatus)

	creditAccountID := createTestAccount(t, testDBPool, validUserID, "RUB")
	setAccountCreditLimit(t, testDBPool, creditAccountID, decimal.NewFromInt(50000))

	duplicateKey := uuid.New().String()

	tests := []struct {
		name                string
		senderPathID        string
		idempotencyKey      string
		payload             map[string]interface{}
		expectedCode        int
		expectedEventsCount int
		expectedSenderBal   decimal.Decimal
		expectedReceiverBal decimal.Decimal
	}{
		{
			name:           "Success - Valid Transfer",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: duplicateKey,
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "15000",
				"currency_code": "RUB",
				"description":   "Rent payment",
			},
			expectedCode:        http.StatusOK,
			expectedEventsCount: 1,
			expectedSenderBal:   decimal.NewFromInt(35000),
			expectedReceiverBal: decimal.NewFromInt(15000),
		},

		{
			name:           "Fail - Duplicate Transaction (Idempotency)",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: duplicateKey,
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "15000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusConflict,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:           "Fail - Insufficient Funds (Not enough money)",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "99900000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:           "Fail - Unknown Receiver",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": uuid.New().String(),
				"amount":        "10000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusNotFound,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},

		{
			name:           "Fail - Missing Idempotency Key",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: "",
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "1000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:           "Fail - Negative Amount",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "-5000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:           "Fail - Zero Amount",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "0",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:           "Fail - Transfer to Self",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": senderAccountID.String(),
				"amount":        "5000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: InitialSenderBalance(),
		},
		{
			name:           "Fail - Invalid Amount Format (String)",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "five-thousand",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:           "Fail - Invalid Receiver UUID Format",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": "not-a-valid-uuid",
				"amount":        "1000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},

		{
			name:           "Success - Cross-Currency Transfer using Credit Limit (RUB to USD)",
			senderPathID:   creditAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": usdAccountID.String(),
				"amount":        "20000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusOK,
			expectedEventsCount: 1,
			expectedSenderBal:   decimal.NewFromInt(-20000),
			expectedReceiverBal: decimal.NewFromInt(216),
		},

		{
			name:           "Success - Cross-Currency Transfer (RUB to USD)",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": usdAccountID.String(),
				"amount":        "10000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusOK,
			expectedEventsCount: 1,
			expectedSenderBal:   decimal.NewFromInt(40000),
			expectedReceiverBal: decimal.NewFromInt(108),
		},
		{
			name:           "Fail - Exchange Service Unavailable (Target EUR)",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": eurAccountID.String(),
				"amount":        "10000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusInternalServerError,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:           "Fail - Invalid Sender UUID in URL",
			senderPathID:   "not-a-valid-uuid",
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "1000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:           "Success - Transfer using Credit Limit",
			senderPathID:   creditAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "10000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusOK,
			expectedEventsCount: 1,
			expectedSenderBal:   decimal.NewFromInt(-10000),
			expectedReceiverBal: decimal.NewFromInt(10000),
		},
		{
			name:           "Fail - Exceeding Credit Limit",
			senderPathID:   creditAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "50001",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   decimal.Zero,
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:           "Fail - Sender Account Inactive",
			senderPathID:   frozenAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "1000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:           "Fail - Receiver Account Inactive",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": frozenAccountID.String(),
				"amount":        "1000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: InitialSenderBalance(),
		},
		{
			name:                "Fail - Empty Payload Body",
			senderPathID:        senderAccountID.String(),
			idempotencyKey:      uuid.New().String(),
			payload:             map[string]interface{}{},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:           "Fail - Amount Overflow (Astronomical Number)",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"to_account_id": receiverAccountID.String(),
				"amount":        "9999999999999999999999999999999999",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
		},
		{
			name:           "Fail - Missing Required Fields",
			senderPathID:   senderAccountID.String(),
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"amount":        "1000",
				"currency_code": "RUB",
			},
			expectedCode:        http.StatusBadRequest,
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
			setAccountBalance(t, testDBPool, creditAccountID, decimal.Zero)
			setAccountBalance(t, testDBPool, frozenAccountID, InitialSenderBalance())
			setAccountBalance(t, testDBPool, usdAccountID, decimal.Zero)
			setAccountBalance(t, testDBPool, eurAccountID, decimal.Zero)

			bodyBytes, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			url := fmt.Sprintf("/api/accounts/%s/transfer", tt.senderPathID)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			if tt.idempotencyKey != "" {
				req.Header.Set("Idempotency-Key", tt.idempotencyKey)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "Response body: %s", rr.Body.String())
			require.Len(t, mockPub.Events, tt.expectedEventsCount)

			if parsedSenderUUID, err := uuid.Parse(tt.senderPathID); err == nil {
				actualSenderBal := getAccountBalance(t, testDBPool, parsedSenderUUID)
				require.True(t, tt.expectedSenderBal.Equal(actualSenderBal),
					"Sender balance mismatch: expected %s, got %s", tt.expectedSenderBal.String(), actualSenderBal.String())
			}

			if receiverStr, ok := tt.payload["to_account_id"].(string); ok {
				if parsedReceiverUUID, err := uuid.Parse(receiverStr); err == nil {
					if tt.expectedCode != http.StatusNotFound {
						actualReceiverBal := getAccountBalance(t, testDBPool, parsedReceiverUUID)
						require.True(t, tt.expectedReceiverBal.Equal(actualReceiverBal),
							"Receiver balance mismatch: expected %s, got %s", tt.expectedReceiverBal.String(), actualReceiverBal.String())
					}
				}
			}
		})
	}
}
