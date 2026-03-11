package vault

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	pb "github.com/Adopten123/banking-system/service-account/internal/infrastructure/grpc/pb/card_vault/v1"
)

func (c *CardVaultGRPCClient) VerifyCard(ctx context.Context, input domain.VerifyCardInput) (bool, string, error) {
	req := &pb.VerifyCardRequest{
		Pan:         input.PAN,
		Cvv:         input.CVV,
		ExpiryMonth: input.ExpiryMonth,
		ExpiryYear:  input.ExpiryYear,
	}

	resp, err := c.client.VerifyCard(ctx, req)
	if err != nil {
		return false, "", fmt.Errorf("grpc vault VerifyCard failed: %w", err)
	}

	return resp.IsValid, resp.TokenId, nil
}
