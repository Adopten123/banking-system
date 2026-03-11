package vault

import (
	"context"
	"fmt"

	pb "github.com/Adopten123/banking-system/service-account/internal/infrastructure/grpc/pb/card_vault/v1"
)

func (c *CardVaultGRPCClient) DeleteCardData(ctx context.Context, tokenID string) error {
	req := &pb.DeleteCardDataRequest{
		TokenId: tokenID,
	}

	resp, err := c.client.DeleteCardData(ctx, req)
	if err != nil {
		return fmt.Errorf("grpc vault DeleteCardData failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("vault refused to delete card data")
	}

	return nil
}