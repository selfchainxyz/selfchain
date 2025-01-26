package mocks

import (
	"context"

	"selfchain/x/keyless/types"
)

// MockTSSProtocol is a mock implementation of the TSSProtocol interface for testing
type MockTSSProtocol struct {
	InitiateSigningFn    func(ctx context.Context, msg []byte, walletID string) (*types.SigningResponse, error)
	GenerateKeySharesFn  func(ctx context.Context, req *types.KeyGenRequest) (*types.KeyGenResponse, error)
	ProcessKeyGenRoundFn func(ctx context.Context, sessionID string, partyData *types.PartyData) error
}

// InitiateSigning implements the TSSProtocol interface
func (m *MockTSSProtocol) InitiateSigning(ctx context.Context, msg []byte, walletID string) (*types.SigningResponse, error) {
	if m.InitiateSigningFn != nil {
		return m.InitiateSigningFn(ctx, msg, walletID)
	}
	return &types.SigningResponse{
		Signature: []byte("mock_signature"),
	}, nil
}

// GenerateKeyShares implements the TSSProtocol interface
func (m *MockTSSProtocol) GenerateKeyShares(ctx context.Context, req *types.KeyGenRequest) (*types.KeyGenResponse, error) {
	if m.GenerateKeySharesFn != nil {
		return m.GenerateKeySharesFn(ctx, req)
	}
	return nil, nil
}

// ProcessKeyGenRound implements the TSSProtocol interface
func (m *MockTSSProtocol) ProcessKeyGenRound(ctx context.Context, sessionID string, partyData *types.PartyData) error {
	if m.ProcessKeyGenRoundFn != nil {
		return m.ProcessKeyGenRoundFn(ctx, sessionID, partyData)
	}
	return nil
}
