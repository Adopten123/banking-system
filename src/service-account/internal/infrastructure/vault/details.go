package vault

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	pb "github.com/Adopten123/banking-system/service-account/internal/infrastructure/grpc/pb/card_vault/v1"
)

func (c *CardVaultGRPCClient) GetCardDetails(ctx context.Context, tokenID string) (*domain.CardDetails, error) {
	req := &pb.GetCardDetailsRequest{
		TokenId: tokenID,
	}

	resp, err := c.client.GetCardDetails(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc vault GetCardDetails failed: %w", err)
	}

	return &domain.CardDetails{
		PAN:         resp.Pan,
		CVV:         resp.Cvv,
		ExpiryMonth: resp.ExpiryMonth,
		ExpiryYear:  resp.ExpiryYear,
	}, nil
}
