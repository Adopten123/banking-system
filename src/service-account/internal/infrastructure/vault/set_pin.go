package vault

import (
	"context"
	"fmt"

	pb "github.com/Adopten123/banking-system/service-account/internal/infrastructure/grpc/pb/card_vault/v1"
)

func (c *CardVaultGRPCClient) SetPin(ctx context.Context, tokenID string, pin string) error {
	req := &pb.SetPinRequest{
		TokenId: tokenID,
		Pin:     pin,
	}

	resp, err := c.client.SetPin(ctx, req)
	if err != nil {
		return fmt.Errorf("grpc vault SetPin failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("vault refused to set PIN")
	}

	return nil
}
