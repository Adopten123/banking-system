package vault

import (
	"context"
	"fmt"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	pb "github.com/Adopten123/banking-system/service-account/internal/infrastructure/grpc/pb/card_vault/v1"
)

func (c *CardVaultGRPCClient) IssueCard(ctx context.Context, params domain.IssueCardParams) (domain.IssuedCardData, error) {
	req := &pb.IssueCardRequest{
		PaymentSystem: params.PaymentSystem,
		IsVirtual:     params.IsVirtual,
	}

	resp, err := c.client.IssueCard(ctx, req)
	if err != nil {
		return domain.IssuedCardData{}, fmt.Errorf("failed to issue card via grpc vault: %w", err)
	}

	return domain.IssuedCardData{
		TokenID:     resp.TokenId,
		PANMask:     resp.PanMask,
		ExpiryMonth: resp.ExpiryMonth,
		ExpiryYear:  resp.ExpiryYear,
	}, nil
}