package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

// MockKeylessKeeper is a mock implementation of KeylessKeeper for testing
type MockKeylessKeeper struct{}

func NewMockKeylessKeeper() *MockKeylessKeeper {
	return &MockKeylessKeeper{}
}

func (m *MockKeylessKeeper) ReconstructWallet(ctx sdk.Context, didDoc types.DIDDocument) ([]byte, error) {
	// Mock implementation
	return []byte("mock_reconstructed_wallet"), nil
}

func (m *MockKeylessKeeper) StoreKeyShare(ctx sdk.Context, did string, keyShare []byte) error {
	// Mock implementation
	return nil
}

func (m *MockKeylessKeeper) GetKeyShare(ctx sdk.Context, did string) ([]byte, bool) {
	// Mock implementation
	return []byte("mock_key_share"), true
}

func (m *MockKeylessKeeper) InitiateRecovery(ctx sdk.Context, did string, recoveryToken string, recoveryAddress string) error {
	// Mock implementation
	return nil
}
