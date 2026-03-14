package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func InitialSenderBalance() decimal.Decimal {
	return decimal.NewFromInt(50000)
}

func TestTransfer_TableDriven(t *testing.T) {
	router, mockPub, mockVault := setupTestEnv(testDBPool)

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

	senderCardUUID := uuid.New()
	senderPAN := "4276000011112222"
	_, err := testDBPool.Exec(context.Background(),
		"INSERT INTO cards (id, account_id, pan_mask, status) VALUES ($1, (SELECT id FROM accounts WHERE public_id = $2), $3, 'active')",
		senderCardUUID, senderAccountID, "4276 **** **** 2222")
	require.NoError(t, err)

	duplicateKey := uuid.New().String()

	tests := []struct {
		name                string
		idempotencyKey      string
		payload             map[string]interface{}
		setupMock           func()
		expectedCode        int
		expectedEventsCount int
		expectedSenderBal   decimal.Decimal
		expectedReceiverBal decimal.Decimal
		checkSenderAccID    uuid.UUID
		checkReceiverAccID  uuid.UUID
	}{
		{
			name:           "Success - Valid Transfer Account to Account",
			idempotencyKey: duplicateKey,
			payload: map[string]interface{}{
				"source_type":      "account",
				"source_id":        senderAccountID.String(),
				"destination_type": "account",
				"destination_id":   receiverAccountID.String(),
				"amount":           "15000",
				"currency_code":    "RUB",
				"description":      "Rent payment",
			},
			expectedCode:        http.StatusOK,
			expectedEventsCount: 1,
			expectedSenderBal:   decimal.NewFromInt(35000),
			expectedReceiverBal: decimal.NewFromInt(15000),
			checkSenderAccID:    senderAccountID,
			checkReceiverAccID:  receiverAccountID,
		},
		{
			name:           "Success - Valid Transfer Card to Account (PAN to UUID)",
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"source_type":      "card",
				"source_id":        senderPAN,
				"destination_type": "account",
				"destination_id":   receiverAccountID.String(),
				"amount":           "5000",
				"currency_code":    "RUB",
			},
			setupMock: func() {
				mockVault.GetTokenByPanFunc = func(ctx context.Context, pan string) (string, error) {
					if pan == senderPAN {
						return senderCardUUID.String(), nil
					}
					return "", fmt.Errorf("card not found")
				}
			},
			expectedCode:        http.StatusOK,
			expectedEventsCount: 1,
			expectedSenderBal:   decimal.NewFromInt(45000),
			expectedReceiverBal: decimal.NewFromInt(5000),
			checkSenderAccID:    senderAccountID,
			checkReceiverAccID:  receiverAccountID,
		},
		{
			name:           "Fail - Card Not Found in Vault (Wrong PAN)",
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"source_type":      "card",
				"source_id":        "4276000099998888",
				"destination_type": "account",
				"destination_id":   receiverAccountID.String(),
				"amount":           "1000",
				"currency_code":    "RUB",
			},
			setupMock: func() {
				mockVault.GetTokenByPanFunc = func(ctx context.Context, pan string) (string, error) {
					return "", domain.ErrCardNotFound
				}
			},
			expectedCode:        http.StatusNotFound,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
			checkSenderAccID:    senderAccountID,
			checkReceiverAccID:  receiverAccountID,
		},
		{
			name:           "Fail - Duplicate Transaction (Idempotency)",
			idempotencyKey: duplicateKey,
			payload: map[string]interface{}{
				"source_type":      "account",
				"source_id":        senderAccountID.String(),
				"destination_type": "account",
				"destination_id":   receiverAccountID.String(),
				"amount":           "15000",
				"currency_code":    "RUB",
			},
			expectedCode:        http.StatusConflict,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
			checkSenderAccID:    senderAccountID,
			checkReceiverAccID:  receiverAccountID,
		},
		{
			name:           "Fail - Insufficient Funds",
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"source_type":      "account",
				"source_id":        senderAccountID.String(),
				"destination_type": "account",
				"destination_id":   receiverAccountID.String(),
				"amount":           "99900000",
				"currency_code":    "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
			checkSenderAccID:    senderAccountID,
			checkReceiverAccID:  receiverAccountID,
		},
		{
			name:           "Fail - Missing Idempotency Key",
			idempotencyKey: "",
			payload: map[string]interface{}{
				"source_type":      "account",
				"source_id":        senderAccountID.String(),
				"destination_type": "account",
				"destination_id":   receiverAccountID.String(),
				"amount":           "1000",
				"currency_code":    "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
			checkSenderAccID:    senderAccountID,
			checkReceiverAccID:  receiverAccountID,
		},
		{
			name:           "Fail - Transfer to Self",
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"source_type":      "account",
				"source_id":        senderAccountID.String(),
				"destination_type": "account",
				"destination_id":   senderAccountID.String(),
				"amount":           "5000",
				"currency_code":    "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: InitialSenderBalance(),
			checkSenderAccID:    senderAccountID,
			checkReceiverAccID:  senderAccountID,
		},
		{
			name:           "Fail - Invalid Source Format",
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"source_type":      "account",
				"source_id":        "not-a-valid-uuid",
				"destination_type": "account",
				"destination_id":   receiverAccountID.String(),
				"amount":           "1000",
				"currency_code":    "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
			checkSenderAccID:    senderAccountID,
			checkReceiverAccID:  receiverAccountID,
		},
		{
			name:           "Success - Cross-Currency Transfer using Credit Limit (RUB to USD)",
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"source_type":      "account",
				"source_id":        creditAccountID.String(),
				"destination_type": "account",
				"destination_id":   usdAccountID.String(),
				"amount":           "20000",
				"currency_code":    "RUB",
			},
			expectedCode:        http.StatusOK,
			expectedEventsCount: 1,
			expectedSenderBal:   decimal.NewFromInt(-20000),
			expectedReceiverBal: decimal.NewFromInt(216),
			checkSenderAccID:    creditAccountID,
			checkReceiverAccID:  usdAccountID,
		},
		{
			name:           "Fail - Source Account Inactive",
			idempotencyKey: uuid.New().String(),
			payload: map[string]interface{}{
				"source_type":      "account",
				"source_id":        frozenAccountID.String(),
				"destination_type": "account",
				"destination_id":   receiverAccountID.String(),
				"amount":           "1000",
				"currency_code":    "RUB",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedSenderBal:   InitialSenderBalance(),
			expectedReceiverBal: decimal.Zero,
			checkSenderAccID:    frozenAccountID,
			checkReceiverAccID:  receiverAccountID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPub.Clear()
			if tt.setupMock != nil {
				tt.setupMock()
			}

			setAccountBalance(t, testDBPool, senderAccountID, InitialSenderBalance())
			setAccountBalance(t, testDBPool, receiverAccountID, decimal.Zero)
			setAccountBalance(t, testDBPool, creditAccountID, decimal.Zero)
			setAccountBalance(t, testDBPool, frozenAccountID, InitialSenderBalance())
			setAccountBalance(t, testDBPool, usdAccountID, decimal.Zero)
			setAccountBalance(t, testDBPool, eurAccountID, decimal.Zero)

			bodyBytes, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/transfers", bytes.NewReader(bodyBytes))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			if tt.idempotencyKey != "" {
				req.Header.Set("Idempotency-Key", tt.idempotencyKey)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "Response body: %s", rr.Body.String())
			require.Len(t, mockPub.Events, tt.expectedEventsCount)

			if tt.checkSenderAccID != uuid.Nil {
				actualSenderBal := getAccountBalance(t, testDBPool, tt.checkSenderAccID)
				require.True(t, tt.expectedSenderBal.Equal(actualSenderBal),
					"Sender balance mismatch: expected %s, got %s", tt.expectedSenderBal.String(), actualSenderBal.String())
			}

			if tt.checkReceiverAccID != uuid.Nil && tt.expectedCode != http.StatusNotFound && tt.expectedCode != http.StatusBadRequest {
				actualReceiverBal := getAccountBalance(t, testDBPool, tt.checkReceiverAccID)
				require.True(t, tt.expectedReceiverBal.Equal(actualReceiverBal),
					"Receiver balance mismatch: expected %s, got %s", tt.expectedReceiverBal.String(), actualReceiverBal.String())
			}
		})
	}
}
