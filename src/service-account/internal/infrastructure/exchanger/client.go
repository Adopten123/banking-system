package exchanger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/shopspring/decimal"
)

type exchangerResponse struct {
	BaseCurrency   string `json:"base_currency"`
	TargetCurrency string `json:"target_currency"`
	Rate           string `json:"rate"`
	Timestamp      string `json:"timestamp"`
}

type HTTPClient struct {
	endpointURL string
	client      *http.Client
}

func NewHTTPClient(endpointURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		endpointURL: endpointURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *HTTPClient) GetRate(ctx context.Context, baseCurrency, targetCurrency string) (decimal.Decimal, error) {
	reqURL, err := url.Parse(c.endpointURL)
	if err != nil {
		return decimal.Zero, fmt.Errorf("invalid endpoint url: %w", err)
	}

	q := reqURL.Query()
	q.Set("base", baseCurrency)
	q.Set("target", targetCurrency)
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return decimal.Zero, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("unexpected status code from exchange service: %d", resp.StatusCode)
	}

	var rateResp exchangerResponse
	if err := json.NewDecoder(resp.Body).Decode(&rateResp); err != nil {
		return decimal.Zero, fmt.Errorf("failed to decode json response: %w", err)
	}

	rate, err := decimal.NewFromString(rateResp.Rate)
	if err != nil {
		return decimal.Zero, fmt.Errorf("invalid rate format from exchange service: %w", err)
	}

	if !rate.IsPositive() {
		return decimal.Zero, fmt.Errorf("exchange rate must be positive, got: %s", rate.String())
	}

	return rate, nil
}
