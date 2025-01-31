package keygen

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"selfchain/x/keyless/testutil/mocks"
	"selfchain/x/keyless/types"
)

func setupTestManager(t *testing.T) *KeyGenManager {
	mockStorage := mocks.SetupMockStorage()
	mgr := NewKeyGenManager(mockStorage)
	return mgr
}

func TestNewKeyGenManager(t *testing.T) {
	mgr := setupTestManager(t)
	require.NotNil(t, mgr)
}

func TestKeyGenManager_GenerateKeyShares(t *testing.T) {
	mgr := setupTestManager(t)
	require.NotNil(t, mgr)

	// Test successful key generation
	request := &types.KeyGenRequest{
		WalletAddress: "test_wallet",
		ChainId:      "test_chain",
		SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_STANDARD,
	}

	// First call should succeed
	shares, err := mgr.GenerateKeyShares(context.Background(), request)
	require.NoError(t, err)
	require.NotNil(t, shares)
	assert.Equal(t, request.WalletAddress, shares.WalletAddress)
	assert.NotNil(t, shares.PublicKey)
	assert.NotNil(t, shares.Metadata)
	assert.Equal(t, request.SecurityLevel, shares.Metadata.SecurityLevel)

	// Second call should fail due to existing key shares
	_, err = mgr.GenerateKeyShares(context.Background(), request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key shares already exist")
}

func TestKeyGenManager_GenerateKeyShares_InvalidRequest(t *testing.T) {
	mgr := setupTestManager(t)
	require.NotNil(t, mgr)

	testCases := []struct {
		name    string
		request *types.KeyGenRequest
		errMsg  string
	}{
		{
			name: "empty wallet address",
			request: &types.KeyGenRequest{
				WalletAddress: "",
				ChainId:      "test_chain",
				SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_STANDARD,
			},
			errMsg: "wallet address cannot be empty",
		},
		{
			name: "empty chain ID",
			request: &types.KeyGenRequest{
				WalletAddress: "test_wallet",
				ChainId:      "",
				SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_STANDARD,
			},
			errMsg: "chain id cannot be empty",
		},
		{
			name:    "nil request",
			request: nil,
			errMsg:  "nil request",
		},
		{
			name: "unspecified security level",
			request: &types.KeyGenRequest{
				WalletAddress: "test_wallet",
				ChainId:      "test_chain",
				SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_UNSPECIFIED,
			},
			errMsg: "security level must be specified",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := mgr.GenerateKeyShares(context.Background(), tc.request)
			assert.Error(t, err)
			if tc.request != nil {
				assert.Contains(t, err.Error(), tc.errMsg)
			}
		})
	}
}
