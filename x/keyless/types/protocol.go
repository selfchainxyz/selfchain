package types

import (
	"context"
)

// TSSProtocol defines the interface for TSS protocol implementations
type TSSProtocol interface {
	// Key Generation
	GenerateKeyShares(ctx context.Context, req *KeyGenRequest) (*KeyGenResponse, error)
	ProcessKeyGenRound(ctx context.Context, sessionID string, partyData *PartyData) error

	// Signing
	InitiateSigning(ctx context.Context, msg []byte, walletID string) (*SigningResponse, error)
}
