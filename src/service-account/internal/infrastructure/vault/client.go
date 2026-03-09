package vault

import (
	pb "github.com/Adopten123/banking-system/service-account/internal/infrastructure/grpc/pb/card_vault/v1"
	"google.golang.org/grpc"
)

type CardVaultGRPCClient struct {
	client pb.CardVaultServiceClient
}

func NewCardVaultClient(conn grpc.ClientConnInterface) *CardVaultGRPCClient {
	return &CardVaultGRPCClient{
		client: pb.NewCardVaultServiceClient(conn),
	}
}
