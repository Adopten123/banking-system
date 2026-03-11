package vault

import (
	"context"
	"fmt"

	pb "github.com/Adopten123/banking-system/service-account/internal/infrastructure/grpc/pb/card_vault/v1"
)

func (c *CardVaultGRPCClient) VerifyPin(ctx context.Context, tokenID string, pin string) (bool, error) {
	req := &pb.VerifyPinRequest{
		TokenId: tokenID,
		Pin:     pin,
	}

	resp, err := c.client.VerifyPin(ctx, req)
	if err != nil {
		return false, fmt.Errorf("grpc vault VerifyPin failed: %w", err)
	}

	return resp.IsValid, nil
}