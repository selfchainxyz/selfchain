package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "selfchain/testutil/keeper"
	"selfchain/x/keyless/types"
)

func TestWalletManagement(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	
	tests := []struct {
		name        string
		wallet      *types.Wallet
		expectError bool
	}{
		{
			name: "valid wallet creation",
			wallet: &types.Wallet{
				Creator:       "self1creator",
				WalletAddress: "self1wallet",
				ChainId:       "self-1",
				Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
				KeyVersion:    1,
			},
			expectError: false,
		},
		{
			name: "duplicate wallet address",
			wallet: &types.Wallet{
				Creator:       "self1creator",
				WalletAddress: "self1wallet",
				ChainId:       "self-1",
				Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
				KeyVersion:    1,
			},
			expectError: true,
		},
		{
			name: "invalid chain ID",
			wallet: &types.Wallet{
				Creator:       "self1creator",
				WalletAddress: "self1wallet3",
				ChainId:       "",
				Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
				KeyVersion:    1,
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := k.SaveWallet(k.Ctx, tc.wallet)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				
				// Verify wallet was created correctly
				wallet, err := k.GetWallet(k.Ctx, tc.wallet.WalletAddress)
				require.NoError(t, err)
				require.Equal(t, tc.wallet.Creator, wallet.Creator)
				require.Equal(t, tc.wallet.ChainId, wallet.ChainId)
				require.Equal(t, tc.wallet.Status, wallet.Status)
			}
		})
	}

	// Test getting all wallets
	wallets, err := k.GetAllWalletsFromStore(k.Ctx)
	require.NoError(t, err)
	require.NotEmpty(t, wallets)
}

func TestWalletAccess(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	
	// Create a test wallet
	wallet := &types.Wallet{
		Creator:       "self1creator",
		WalletAddress: "self1wallet",
		ChainId:       "self-1",
		Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
		KeyVersion:    1,
	}

	err := k.SaveWallet(k.Ctx, wallet)
	require.NoError(t, err)

	// Test wallet access validation
	err = k.ValidateWalletAccess(k.Ctx, wallet.WalletAddress, "sign")
	require.NoError(t, err)

	// Test wallet authorization
	authorized, err := k.IsWalletAuthorized(k.Ctx, wallet.Creator, wallet.WalletAddress)
	require.NoError(t, err)
	require.True(t, authorized)

	// Test unauthorized access
	authorized, err = k.IsWalletAuthorized(k.Ctx, "unauthorized_creator", wallet.WalletAddress)
	require.NoError(t, err)
	require.False(t, authorized)
}
