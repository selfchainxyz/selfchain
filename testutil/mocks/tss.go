package mocks

import (
	"context"
	"time"

	"selfchain/x/keyless/types"
)

// MockTSSProtocol is a mock implementation of the TSS protocol for testing
type MockTSSProtocol struct {
	InitiateSigningFn   func(ctx context.Context, msg []byte, walletID string) (*types.SigningResponse, error)
	GenerateKeySharesFn func(ctx context.Context, req *types.KeyGenRequest) (*types.KeyGenResponse, error)
	ProcessKeyGenRoundFn func(ctx context.Context, sessionID string, partyData *types.PartyData) error
}

// InitiateSigning implements the TSS protocol for testing
func (m *MockTSSProtocol) InitiateSigning(ctx context.Context, msg []byte, walletID string) (*types.SigningResponse, error) {
	if m.InitiateSigningFn != nil {
		return m.InitiateSigningFn(ctx, msg, walletID)
	}

	// Default mock implementation
	now := time.Now().UTC()
	return &types.SigningResponse{
		WalletId:  walletID,
		Signature: []byte("test_signature"),
		Metadata: &types.SignatureMetadata{
			Timestamp: &now,
			ChainId:   "test-chain",
			SignType:  types.SignatureType_SIGNATURE_TYPE_ECDSA,
		},
	}, nil
}

// GenerateKeyShares implements the TSS protocol for testing
func (m *MockTSSProtocol) GenerateKeyShares(ctx context.Context, req *types.KeyGenRequest) (*types.KeyGenResponse, error) {
	if m.GenerateKeySharesFn != nil {
		return m.GenerateKeySharesFn(ctx, req)
	}

	// Default mock implementation
	now := time.Now().UTC()
	return &types.KeyGenResponse{
		WalletId:  req.WalletId,
		PublicKey: []byte("test_public_key"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     now,
			LastRotated:   now,
			LastUsed:      now,
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: req.SecurityLevel,
		},
	}, nil
}

// ProcessKeyGenRound implements the TSS protocol for testing
func (m *MockTSSProtocol) ProcessKeyGenRound(ctx context.Context, sessionID string, partyData *types.PartyData) error {
	if m.ProcessKeyGenRoundFn != nil {
		return m.ProcessKeyGenRoundFn(ctx, sessionID, partyData)
	}

	// Default mock implementation returns nil (success)
	return nil
}
