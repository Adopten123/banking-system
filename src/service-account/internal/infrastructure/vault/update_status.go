package vault

import (
	"context"
	"fmt"

	pb "github.com/Adopten123/banking-system/service-account/internal/infrastructure/grpc/pb/card_vault/v1"
)

func (c *CardVaultGRPCClient) UpdateCardStatus(ctx context.Context, tokenID string, status string) error {
	req := &pb.UpdateCardStatusRequest{
		TokenId: tokenID,
		Status:  status,
	}

	resp, err := c.client.UpdateCardStatus(ctx, req)
	if err != nil {
		return fmt.Errorf("grpc vault UpdateCardStatus failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("vault refused to update card status")
	}

	return nil
}