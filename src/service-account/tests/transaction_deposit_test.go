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

func InitialBalance() decimal.Decimal {
	return decimal.NewFromInt(10000)
}

func TestDeposit_TableDriven(t *testing.T) {
	// Добавляем mockVault, чтобы тестировать пополнение по номеру карты
	router, mockPub, mockVault := setupTestEnv(testDBPool)

	validUserID := createTestUser(t, testDBPool)
	validAccountID := createTestAccount(t, testDBPool, validUserID, "RUB")

	const closedStatus = 4
	closedAccountID := createTestAccount(t, testDBPool, validUserID, "RUB")
	updateAccountStatus(t, testDBPool, closedAccountID, closedStatus)

	// --- СОЗДАЕМ ТЕСТОВУЮ КАРТУ ДЛЯ ПОЛУЧАТЕЛЯ ---
	destCardUUID := uuid.New()
	destPAN := "4276000011112222"
	// Привязываем карту к нашему validAccountID
	_, err := testDBPool.Exec(context.Background(),
		"INSERT INTO cards (id, account_id, pan_mask, status) VALUES ($1, (SELECT id FROM accounts WHERE public_id = $2), $3, 'active')",
		destCardUUID, validAccountID, "4276 **** **** 2222")
	require.NoError(t, err)

	tests := []struct {
		name                string
		payload             map[string]interface{}
		setupMock           func() // Настройка Сейфа
		expectedCode        int
		expectedEventsCount int
		expectedBalance     decimal.Decimal
		checkAccID          uuid.UUID // ID счета, на котором будем проверять баланс
	}{
		{
			name: "Success - Valid Deposit to Account",
			payload: map[string]interface{}{
				"destination_type": "account",
				"destination_id":   validAccountID.String(),
				"amount":           "5000",
			},
			expectedCode:        http.StatusCreated,
			expectedEventsCount: 1,
			expectedBalance:     decimal.NewFromInt(15000),
			checkAccID:          validAccountID,
		},
		{
			name: "Success - Valid Deposit to Card (PAN to UUID)",
			payload: map[string]interface{}{
				"destination_type": "card",
				"destination_id":   destPAN, // Передаем номер карты!
				"amount":           "5000",
			},
			setupMock: func() {
				// Учим Сейф распознавать PAN и возвращать UUID карты
				mockVault.GetTokenByPanFunc = func(ctx context.Context, pan string) (string, error) {
					if pan == destPAN {
						return destCardUUID.String(), nil
					}
					return "", fmt.Errorf("card not found")
				}
			},
			expectedCode:        http.StatusCreated,
			expectedEventsCount: 1,
			expectedBalance:     decimal.NewFromInt(15000), // Баланс изменится на привязанном счете!
			checkAccID:          validAccountID,
		},
		{
			name: "Fail - Card Not Found in Vault",
			payload: map[string]interface{}{
				"destination_type": "card",
				"destination_id":   "4276000099998888", // Левый номер
				"amount":           "5000",
			},
			setupMock: func() {
				mockVault.GetTokenByPanFunc = func(ctx context.Context, pan string) (string, error) {
					return "", domain.ErrCardNotFound // Сейф отклоняет
				}
			},
			expectedCode:        http.StatusNotFound, // Ожидаем 404
			expectedEventsCount: 0,
			expectedBalance:     InitialBalance(),
			checkAccID:          validAccountID,
		},
		{
			name: "Fail - Negative Amount",
			payload: map[string]interface{}{
				"destination_type": "account",
				"destination_id":   validAccountID.String(),
				"amount":           "-1000",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedBalance:     InitialBalance(),
			checkAccID:          validAccountID,
		},
		{
			name: "Fail - Zero Amount",
			payload: map[string]interface{}{
				"destination_type": "account",
				"destination_id":   validAccountID.String(),
				"amount":           "0",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedBalance:     InitialBalance(),
			checkAccID:          validAccountID,
		},
		{
			name: "Fail - Unknown Account",
			payload: map[string]interface{}{
				"destination_type": "account",
				"destination_id":   uuid.New().String(),
				"amount":           "5000",
			},
			expectedCode:        http.StatusNotFound,
			expectedEventsCount: 0,
			expectedBalance:     InitialBalance(),
			checkAccID:          validAccountID, // Просто проверим, что основной счет не изменился
		},
		{
			name: "Success - Minimum Valid Amount (1 kopeck)",
			payload: map[string]interface{}{
				"destination_type": "account",
				"destination_id":   validAccountID.String(),
				"amount":           "1",
			},
			expectedCode:        http.StatusCreated,
			expectedEventsCount: 1,
			expectedBalance:     InitialBalance().Add(decimal.NewFromInt(1)),
			checkAccID:          validAccountID,
		},
		{
			name: "Success - Gigantic Amount",
			payload: map[string]interface{}{
				"destination_type": "account",
				"destination_id":   validAccountID.String(),
				"amount":           "9999999999999999999999",
			},
			expectedCode:        http.StatusCreated,
			expectedEventsCount: 1,
			expectedBalance:     InitialBalance().Add(decimal.RequireFromString("9999999999999999999999")),
			checkAccID:          validAccountID,
		},
		{
			name: "Fail - Deposit to Closed Account",
			payload: map[string]interface{}{
				"destination_type": "account",
				"destination_id":   closedAccountID.String(),
				"amount":           "5000",
			},
			expectedCode:        http.StatusForbidden,
			expectedEventsCount: 0,
			expectedBalance:     InitialBalance(), // Зададим начальный баланс и для closedAccountID
			checkAccID:          closedAccountID,
		},
		{
			name: "Fail - Invalid JSON Amount Format (Letters)",
			payload: map[string]interface{}{
				"destination_type": "account",
				"destination_id":   validAccountID.String(),
				"amount":           "abc",
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedBalance:     InitialBalance(),
			checkAccID:          validAccountID,
		},
		{
			name: "Fail - Missing Fields",
			payload: map[string]interface{}{
				"amount": "5000", // Забыли указать куда пополнять
			},
			expectedCode:        http.StatusBadRequest,
			expectedEventsCount: 0,
			expectedBalance:     InitialBalance(),
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
			setAccountBalance(t, testDBPool, validAccountID, InitialBalance())
			setAccountBalance(t, testDBPool, closedAccountID, InitialBalance())

			bodyBytes, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			// URL теперь статичный
			req, err := http.NewRequest(http.MethodPost, "/api/deposits", bytes.NewReader(bodyBytes))
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
