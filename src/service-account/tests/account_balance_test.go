package tests

import (
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

func TestGetAccountBalance_TableDriven(t *testing.T) {
	router, _ := setupTestEnv(testDBPool)

	validUserID := createTestUser(t, testDBPool)
	validAccountID := createTestAccount(t, testDBPool, validUserID, "RUB")

	expectedBalance := decimal.NewFromFloat(99950)
	setAccountBalance(t, testDBPool, validAccountID, expectedBalance)

	tests := []struct {
		name          string
		accountPathID string
		expectedCode  int
	}{
		{
			name:          "Success - Valid Account Balance",
			accountPathID: validAccountID.String(),
			expectedCode:  http.StatusOK,
		},
		{
			name:          "Fail - Unknown Account",
			accountPathID: uuid.New().String(),
			expectedCode:  http.StatusNotFound,
		},
		{
			name:          "Fail - Invalid UUID Format",
			accountPathID: "not-a-valid-uuid",
			expectedCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/accounts/%s/balance", tt.accountPathID)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "Response body: %s", rr.Body.String())

			if tt.expectedCode == http.StatusOK {
				var resp domain.AccountBalanceResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err, "failed to unmarshal JSON response")

				require.Equal(t, validAccountID, resp.AccountID, "returned account ID mismatch")

				require.True(t, expectedBalance.Equal(resp.Balance),
					"Balance mismatch: expected %s, got %s", expectedBalance.String(), resp.Balance.String())
			}
		})
	}
}
