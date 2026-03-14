package vault

import (
	"context"
	"fmt"

	pb "github.com/Adopten123/banking-system/service-account/internal/infrastructure/grpc/pb/card_vault/v1"
)

func (c *CardVaultGRPCClient) GetTokenByPan(ctx context.Context, pan string) (string, error) {
	req := &pb.GetTokenByPanRequest{
		Pan: pan,
	}

	resp, err := c.client.GetTokenByPan(ctx, req)
	if err != nil {
		return "", fmt.Errorf("grpc vault GetTokenByPan failed: %w", err)
	}

	return resp.TokenId, nil
}
