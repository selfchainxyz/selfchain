package mocks

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"selfchain/x/keyless/types"
)

// TSSProtocol is a mock implementation of the TSSProtocol interface
type TSSProtocol struct {
	mock.Mock
}

// NewTSSProtocol creates a new mock TSSProtocol
func NewTSSProtocol() *TSSProtocol {
	return &TSSProtocol{}
}

// GenerateKeyShares implements the TSSProtocol interface
func (m *TSSProtocol) GenerateKeyShares(ctx sdk.Context, walletAddress string, threshold uint32, securityLevel types.SecurityLevel) (*types.KeyGenResponse, error) {
	args := m.Called(ctx, walletAddress, threshold, securityLevel)
	if resp, ok := args.Get(0).(*types.KeyGenResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

// ReconstructKey implements the TSSProtocol interface
func (m *TSSProtocol) ReconstructKey(ctx sdk.Context, shares [][]byte) ([]byte, error) {
	args := m.Called(ctx, shares)
	return args.Get(0).([]byte), args.Error(1)
}

// VerifyShare implements the TSSProtocol interface
func (m *TSSProtocol) VerifyShare(ctx sdk.Context, share []byte, publicKey []byte) error {
	args := m.Called(ctx, share, publicKey)
	return args.Error(0)
}

// SignMessage implements the TSSProtocol interface
func (m *TSSProtocol) SignMessage(ctx sdk.Context, message []byte, shares [][]byte) ([]byte, error) {
	args := m.Called(ctx, message, shares)
	return args.Get(0).([]byte), args.Error(1)
}

// VerifySignature implements the TSSProtocol interface
func (m *TSSProtocol) VerifySignature(ctx sdk.Context, message []byte, signature []byte, publicKey []byte) error {
	args := m.Called(ctx, message, signature, publicKey)
	return args.Error(0)
}

// GetPartyData implements the TSSProtocol interface
func (m *TSSProtocol) GetPartyData(ctx sdk.Context, partyID string) (*types.PartyData, error) {
	args := m.Called(ctx, partyID)
	if resp, ok := args.Get(0).(*types.PartyData); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

// SetPartyData implements the TSSProtocol interface
func (m *TSSProtocol) SetPartyData(ctx sdk.Context, data *types.PartyData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}
