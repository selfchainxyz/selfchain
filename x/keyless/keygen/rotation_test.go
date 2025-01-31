package keygen

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"selfchain/x/keyless/testutil/mocks"
	"selfchain/x/keyless/types"
)

func setupTestRotationManager(t *testing.T) (*KeyRotationManager, *mocks.MockStorage) {
	mockStorage := mocks.SetupMockStorage()
	keyGen := NewKeyGenManager(mockStorage)
	rotationManager := NewKeyRotationManager(keyGen, mockStorage)
	return rotationManager, mockStorage
}

func TestNewKeyRotationManager(t *testing.T) {
	mgr, _ := setupTestRotationManager(t)
	require.NotNil(t, mgr)
}

func TestKeyRotationManager_RotateKeyShares(t *testing.T) {
	mgr, mockStorage := setupTestRotationManager(t)
	require.NotNil(t, mgr)

	// First generate initial key shares
	walletAddress := "test_wallet"
	chainId := "test_chain"
	securityLevel := types.SecurityLevel_SECURITY_LEVEL_STANDARD

	request := &types.KeyGenRequest{
		WalletAddress: walletAddress,
		ChainId:      chainId,
		SecurityLevel: securityLevel,
	}

	// Pre-populate metadata
	err := mocks.PrePopulateMetadata(mockStorage, walletAddress, chainId, securityLevel)
	require.NoError(t, err)

	// Generate initial shares
	_, err = mgr.keyGen.GenerateKeyShares(context.Background(), request)
	require.NoError(t, err)

	// Test successful key rotation
	err = mgr.RotateKeyShares(context.Background(), walletAddress)
	require.NoError(t, err)

	// Test rotation with non-existent wallet
	err = mgr.RotateKeyShares(context.Background(), "non_existent_wallet")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "metadata not found")
}

func TestKeyRotationManager_ShouldRotate(t *testing.T) {
	mgr, _ := setupTestRotationManager(t)
	require.NotNil(t, mgr)

	now := time.Now()
	testCases := []struct {
		name     string
		metadata *types.KeyMetadata
		want     bool
	}{
		{
			name:     "nil metadata",
			metadata: nil,
			want:     false,
		},
		{
			name: "high security level - should rotate",
			metadata: &types.KeyMetadata{
				SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
				CreatedAt:    now.Add(-31 * 24 * time.Hour),
				LastRotated:  now.Add(-31 * 24 * time.Hour),
			},
			want: true,
		},
		{
			name: "high security level - should not rotate",
			metadata: &types.KeyMetadata{
				SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
				CreatedAt:    now.Add(-29 * 24 * time.Hour),
				LastRotated:  now.Add(-29 * 24 * time.Hour),
			},
			want: false,
		},
		{
			name: "standard security level - should rotate",
			metadata: &types.KeyMetadata{
				SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_STANDARD,
				CreatedAt:    now.Add(-91 * 24 * time.Hour),
				LastRotated:  now.Add(-91 * 24 * time.Hour),
			},
			want: true,
		},
		{
			name: "standard security level - should not rotate",
			metadata: &types.KeyMetadata{
				SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_STANDARD,
				CreatedAt:    now.Add(-89 * 24 * time.Hour),
				LastRotated:  now.Add(-89 * 24 * time.Hour),
			},
			want: false,
		},
		{
			name: "enterprise security level - should rotate",
			metadata: &types.KeyMetadata{
				SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_ENTERPRISE,
				CreatedAt:    now.Add(-181 * 24 * time.Hour),
				LastRotated:  now.Add(-181 * 24 * time.Hour),
			},
			want: true,
		},
		{
			name: "enterprise security level - should not rotate",
			metadata: &types.KeyMetadata{
				SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_ENTERPRISE,
				CreatedAt:    now.Add(-179 * 24 * time.Hour),
				LastRotated:  now.Add(-179 * 24 * time.Hour),
			},
			want: false,
		},
		{
			name: "unknown security level",
			metadata: &types.KeyMetadata{
				SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_UNSPECIFIED,
				CreatedAt:    now.Add(-91 * 24 * time.Hour),
				LastRotated:  now.Add(-91 * 24 * time.Hour),
			},
			want: false,
		},
		{
			name: "zero last rotation time",
			metadata: &types.KeyMetadata{
				SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_STANDARD,
				CreatedAt:    now.Add(-91 * 24 * time.Hour),
				LastRotated:  time.Time{},
			},
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := mgr.ShouldRotate(tc.metadata)
			assert.Equal(t, tc.want, got)
		})
	}
}
