package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/delivery/worker"
	"github.com/Adopten123/banking-system/service-account/internal/repository/postgres"
	"github.com/Adopten123/banking-system/service-account/internal/service"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func getRecurringPayment(t *testing.T, id uuid.UUID) postgres.RecurringPayment {
	var p postgres.RecurringPayment
	query := `SELECT id, source_type_id, source_id, amount, is_active, next_execution_time FROM recurring_payments WHERE id = $1`
	err := testDBPool.QueryRow(context.Background(), query, id).Scan(
		&p.ID, &p.SourceTypeID, &p.SourceID, &p.Amount, &p.IsActive, &p.NextExecutionTime,
	)
	require.NoError(t, err)
	return p
}

func TestRecurringPayments_API_TableDriven(t *testing.T) {
	router, _, _ := setupTestEnv(testDBPool)

	validUserID := createTestUser(t, testDBPool)
	sourceAccountID := createTestAccount(t, testDBPool, validUserID, "RUB")
	destAccountID := createTestAccount(t, testDBPool, validUserID, "RUB")

	tests := []struct {
		name         string
		payload      map[string]interface{}
		expectedCode int
	}{
		{
			name: "Success - Create valid recurring payment (Transfer)",
			payload: map[string]interface{}{
				"source_type":      "account",
				"source_id":        sourceAccountID.String(),
				"destination_type": "account",
				"destination_id":   destAccountID.String(),
				"amount":           "500",
				"currency_code":    "RUB",
				"cron_expression":  "0 12 1 * *",
				"description":      "Internet Bill",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "Fail - Invalid CRON Expression",
			payload: map[string]interface{}{
				"source_type":     "account",
				"source_id":       sourceAccountID.String(),
				"amount":          "500",
				"cron_expression": "invalid-cron",
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Fail - Source Account Not Found",
			payload: map[string]interface{}{
				"source_type":     "account",
				"source_id":       uuid.New().String(),
				"amount":          "500",
				"cron_expression": "0 12 1 * *",
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/recurring-payments", bytes.NewReader(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "Response body: %s", rr.Body.String())
		})
	}
}

func TestRecurringWorker_ProcessDuePayments(t *testing.T) {

	repo := postgres.NewAccountRepo(testDBPool)
	mockPub := NewMockPublisher()
	mockVault := &MockCardVaultClient{}
	mockExch := &MockExchangeClient{}
	svc := service.NewAccountService(repo, mockPub, mockExch, mockVault)

	recWorker := worker.NewRecurringWorker(repo, svc)

	validUserID := createTestUser(t, testDBPool)
	sourceAccID := createTestAccount(t, testDBPool, validUserID, "RUB")
	destAccID := createTestAccount(t, testDBPool, validUserID, "RUB")

	setAccountBalance(t, testDBPool, sourceAccID, decimal.NewFromInt(5000))

	paymentID := uuid.New()
	pastTime := time.Now().Add(-1 * time.Hour).UTC()

	cronExpr := "* * * * *"

	_, err := testDBPool.Exec(context.Background(), `
		INSERT INTO recurring_payments (
			id, source_type_id, source_id, destination_type_id, destination_id, 
			amount, currency_code, cron_expression, next_execution_time, is_active
		) VALUES (
			$1, 1, $2, 1, $3, 1000, 'RUB', $4, $5, true
		)
	`, paymentID, sourceAccID, destAccID, cronExpr, pastTime)
	require.NoError(t, err)

	recWorker.ProcessDuePayments(context.Background())

	senderBalance := getAccountBalance(t, testDBPool, sourceAccID)
	require.True(t, decimal.NewFromInt(4000).Equal(senderBalance), "Sender balance should be deducted")

	receiverBalance := getAccountBalance(t, testDBPool, destAccID)
	require.True(t, decimal.NewFromInt(1000).Equal(receiverBalance), "Receiver balance should be incremented")

	updatedPayment := getRecurringPayment(t, paymentID)
	require.True(t, updatedPayment.NextExecutionTime.Time.After(time.Now().UTC()), "Next execution time should be in the future")

	require.Len(t, mockPub.Events, 1)
	require.Equal(t, "TransferCreatedEvent", mockPub.Events[0].EventName())
}

func TestCancelRecurringPayment_API(t *testing.T) {
	router, _, _ := setupTestEnv(testDBPool)

	validUserID := createTestUser(t, testDBPool)
	sourceAccID := createTestAccount(t, testDBPool, validUserID, "RUB")

	paymentID := uuid.New()

	_, err := testDBPool.Exec(context.Background(), `
		INSERT INTO recurring_payments (
			id, source_type_id, source_id, amount, currency_code, cron_expression, next_execution_time, is_active
		) VALUES (
			$1, 1, $2, 500, 'RUB', '* * * * *', now(), true
		)
	`, paymentID, sourceAccID)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodDelete, "/api/recurring-payments/"+paymentID.String(), nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "Response body: %s", rr.Body.String())

	canceledPayment := getRecurringPayment(t, paymentID)
	require.False(t, canceledPayment.IsActive.Bool, "Payment should be inactive after cancellation")
}

func TestRecurringWorker_ProcessDuePayments_InsufficientFunds(t *testing.T) {
	repo := postgres.NewAccountRepo(testDBPool)
	mockPub := NewMockPublisher()
	mockVault := &MockCardVaultClient{}
	mockExch := &MockExchangeClient{}
	svc := service.NewAccountService(repo, mockPub, mockExch, mockVault)

	recWorker := worker.NewRecurringWorker(repo, svc)

	validUserID := createTestUser(t, testDBPool)
	sourceAccID := createTestAccount(t, testDBPool, validUserID, "RUB")
	destAccID := createTestAccount(t, testDBPool, validUserID, "RUB")

	setAccountBalance(t, testDBPool, sourceAccID, decimal.Zero)

	paymentID := uuid.New()
	pastTime := time.Now().Add(-1 * time.Hour).UTC()

	_, err := testDBPool.Exec(context.Background(), `
		INSERT INTO recurring_payments (
			id, source_type_id, source_id, destination_type_id, destination_id, 
			amount, currency_code, cron_expression, next_execution_time, is_active
		) VALUES (
			$1, 1, $2, 1, $3, 1000, 'RUB', '* * * * *', $4, true
		)
	`, paymentID, sourceAccID, destAccID, pastTime)
	require.NoError(t, err)

	recWorker.ProcessDuePayments(context.Background())

	senderBalance := getAccountBalance(t, testDBPool, sourceAccID)
	require.True(t, decimal.Zero.Equal(senderBalance), "Sender balance should remain 0")

	updatedPayment := getRecurringPayment(t, paymentID)
	require.True(t, updatedPayment.NextExecutionTime.Time.After(time.Now().UTC()), "Next execution time MUST be moved forward even if payment failed")
}
