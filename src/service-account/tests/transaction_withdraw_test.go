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

// Начальный баланс для тестов на снятие
func InitialWithdrawBalance() decimal.Decimal {
	return decimal.NewFromInt(10000)
}

func TestWithdraw_TableDriven(t *testing.T) {
	router, mockPub, mockVault := setupTestEnv(testDBPool)

	validUserID := createTestUser(t, testDBPool)
	validAccountID := createTestAccount(t, testDBPool, validUserID, "RUB")

	const closedStatus = 4
	closedAccountID := createTestAccount(t, testDBPool, validUserID, "RUB")
	updateAccountStatus(t, testDBPool, closedAccountID, closedStatus)

	// --- СОЗДАЕМ ТЕСТОВУЮ КАРТУ ДЛЯ ОТПРАВИТЕЛЯ (ДЛЯ СНЯТИЯ) ---
	sourceCardUUID := uuid.New()
	sourcePAN := "4276000011112222"
	_, err := testDBPool.Exec(context.Background(),
		"INSERT INTO cards (id, account_id, pan_mask, status) VALUES ($1, (SELECT id FROM accounts WHERE public_id = $2), $3, 'active')",
		sourceCardUUID, validAccountID, "4276 **** **** 2222")
	require.NoError(t, err)

	tests := []struct {
		name                string
		payload             map[string]interface{}
		setupMock           func() // Настройка Сейфа
		expectedCode        int
		expectedEventsCount int
		expectedBalance     decimal.Decimal
		checkAccID          uuid.UUID
	}{
		{
			name: "Success - Valid Withdraw from Account",
			payload: map[string]interface{}{
				"source_type": "account",
				"source_id":   validAccountID.String(),
				"amount":      "3000",
			},
			expectedCode:        http.StatusOK, // У тебя в хендлере стоит StatusOK (200) для снятия
			expectedEventsCount: 1,
			expectedBalance:     decimal.NewFromInt(7000), // 10000 - 3000
			checkAccID:          validAccountID,
		},
		{
			name: "Success - Valid Withdraw from Card (PAN to UUID)",
			payload: map[string]interface{}{
				"source_type": "card",
				"source_id":   sourcePAN,
				"amount":      "5000",
			},
			setupMock: func() {
				mockVault.GetTokenByPanFunc = func(ctx context.Context, pan string) (string, error) {
					if pan == sourcePAN {
						return sourceCardUUID.String(), nil
					}
					return "", fmt.Errorf("card not found")
				}
			},
			expectedCode:        http.StatusOK,
			expectedEventsCount: 1,
			expectedBalance:     decimal.NewFromInt(5000), // 10000 - 5000
			checkAccID:          validAccountID,
		},
		{
			name: "Fail - Insufficient Funds",
			payload: map[string]interface{}{
				"source_type": "account",
				"source_id":   validAccountID.String(),
				"amount":      "15000", // Пытаемся снять больше, чем есть
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedBalance:     InitialWithdrawBalance(), // Баланс не изменился
			checkAccID:          validAccountID,
		},
		{
			name: "Fail - Card Not Found in Vault",
			payload: map[string]interface{}{
				"source_type": "card",
				"source_id":   "4276000099998888",
				"amount":      "1000",
			},
			setupMock: func() {
				mockVault.GetTokenByPanFunc = func(ctx context.Context, pan string) (string, error) {
					return "", domain.ErrCardNotFound
				}
			},
			expectedCode:        http.StatusNotFound, // 404
			expectedEventsCount: 0,
			expectedBalance:     InitialWithdrawBalance(),
			checkAccID:          validAccountID,
		},
		{
			name: "Fail - Negative Amount",
			payload: map[string]interface{}{
				"source_type": "account",
				"source_id":   validAccountID.String(),
				"amount":      "-1000",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedBalance:     InitialWithdrawBalance(),
			checkAccID:          validAccountID,
		},
		{
			name: "Fail - Zero Amount",
			payload: map[string]interface{}{
				"source_type": "account",
				"source_id":   validAccountID.String(),
				"amount":      "0",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedBalance:     InitialWithdrawBalance(),
			checkAccID:          validAccountID,
		},
		{
			name: "Fail - Unknown Account",
			payload: map[string]interface{}{
				"source_type": "account",
				"source_id":   uuid.New().String(),
				"amount":      "1000",
			},
			expectedCode:        http.StatusNotFound,
			expectedEventsCount: 0,
			expectedBalance:     InitialWithdrawBalance(),
			checkAccID:          validAccountID, // Основной счет не тронут
		},
		{
			name: "Fail - Withdraw from Closed Account",
			payload: map[string]interface{}{
				"source_type": "account",
				"source_id":   closedAccountID.String(),
				"amount":      "1000",
			},
			expectedCode:        http.StatusBadRequest, // У тебя ErrAccountInactive мапится в 400 Bad Request
			expectedEventsCount: 0,
			expectedBalance:     InitialWithdrawBalance(),
			checkAccID:          closedAccountID,
		},
		{
			name: "Fail - Invalid JSON Amount Format (Letters)",
			payload: map[string]interface{}{
				"source_type": "account",
				"source_id":   validAccountID.String(),
				"amount":      "abc",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedBalance:     InitialWithdrawBalance(),
			checkAccID:          validAccountID,
		},
		{
			name: "Fail - Missing Fields",
			payload: map[string]interface{}{
				"amount": "5000", // Забыли указать откуда снимать
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedBalance:     InitialWithdrawBalance(),
			checkAccID:          validAccountID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPub.Clear()
			if tt.setupMock != nil {
				tt.setupMock()
			}

			// Сбрасываем балансы перед каждым кейсом
			setAccountBalance(t, testDBPool, validAccountID, InitialWithdrawBalance())
			setAccountBalance(t, testDBPool, closedAccountID, InitialWithdrawBalance())

			bodyBytes, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			// Стучимся на новый URL снятия
			req, err := http.NewRequest(http.MethodPost, "/api/withdrawals", bytes.NewReader(bodyBytes))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Idempotency-Key", uuid.New().String())

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "Response body: %s", rr.Body.String())
			require.Len(t, mockPub.Events, tt.expectedEventsCount)

			// Проверяем баланс по явно указанному в тесте счету
			if tt.checkAccID != uuid.Nil {
				actualBalance := getAccountBalance(t, testDBPool, tt.checkAccID)
				require.True(t, tt.expectedBalance.Equal(actualBalance),
					"Balance mismatch: expected %s, got %s", tt.expectedBalance.String(), actualBalance.String())
			}
		})
	}
}
