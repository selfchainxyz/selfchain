package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"selfchain/x/keyless/types"
)

type TSSProtocol struct {
	mock.Mock
}

// NewTSSProtocol creates a new mock TSSProtocol
func NewTSSProtocol() *TSSProtocol {
	return &TSSProtocol{}
}

func (m *TSSProtocol) GenerateKey(ctx context.Context, parties []string, threshold int) ([]byte, error) {
	args := m.Called(ctx, parties, threshold)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *TSSProtocol) GenerateKeyShares(ctx context.Context, req *types.KeyGenRequest) (*types.KeyGenResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*types.KeyGenResponse), args.Error(1)
}

func (m *TSSProtocol) Sign(ctx context.Context, sessionID string, message []byte) ([]byte, error) {
	args := m.Called(ctx, sessionID, message)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *TSSProtocol) VerifySignature(pubKey []byte, message []byte, signature []byte) error {
	args := m.Called(pubKey, message, signature)
	return args.Error(0)
}

func (m *TSSProtocol) InitiateSigning(ctx context.Context, message []byte, walletID string) (*types.SigningResponse, error) {
	args := m.Called(ctx, message, walletID)
	return args.Get(0).(*types.SigningResponse), args.Error(1)
}

func (m *TSSProtocol) ProcessKeyGenRound(ctx context.Context, sessionID string, partyData *types.PartyData) error {
	args := m.Called(ctx, sessionID, partyData)
	return args.Error(0)
}
