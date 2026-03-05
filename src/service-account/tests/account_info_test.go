package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetAccount_TableDriven(t *testing.T) {
	router, _ := setupTestEnv(testDBPool)

	validUserID := createTestUser(t, testDBPool)
	validAccountID := createTestAccount(t, testDBPool, validUserID, "RUB")

	tests := []struct {
		name          string
		accountPathID string
		expectedCode  int
	}{
		{
			name:          "Success - Valid Account",
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
			url := fmt.Sprintf("/api/accounts/%s", tt.accountPathID)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code, "Response body: %s", rr.Body.String())

			if tt.expectedCode == http.StatusOK {
				var resp map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err)

				returnedID, exists := resp["id"]
				require.True(t, exists, "JSON response must contain 'id' field")
				require.Equal(t, tt.accountPathID, returnedID, "Returned ID does not match requested ID")

				require.Equal(t, "RUB", resp["currency_code"])
			}
		})
	}
}
