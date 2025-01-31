package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "selfchain/testutil/keeper"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
)

func TestWalletManagement(t *testing.T) {
	k := keepertest.NewKeylessKeeper(t)
	
	tests := []struct {
		name        string
		wallet      *types.Wallet
		expectError bool
	}{
		{
			name: "valid wallet creation",
			wallet: &types.Wallet{
				Creator:       "cosmos1w3jhxap3ta047h6lta047h6lta047h6lx84s66",
				WalletAddress: "cosmos1w3jhxap3ta047h6lta047h6lta047h6lx84s66",
				ChainId:       "self-1",
				Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
				KeyVersion:    1,
			},
			expectError: false,
		},
		{
			name: "duplicate wallet address",
			wallet: &types.Wallet{
				Creator:       "cosmos1w3jhxap3ta047h6lta047h6lta047h6lx84s66",
				WalletAddress: "cosmos1w3jhxap3ta047h6lta047h6lta047h6lx84s66",
				ChainId:       "self-1",
				Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
				KeyVersion:    1,
			},
			expectError: true,
		},
		{
			name: "invalid chain ID",
			wallet: &types.Wallet{
				Creator:       "cosmos1w3jhxap3ta047h6lta047h6lta047h6lx84s66",
				WalletAddress: "cosmos1w3jhxapjta047h6lta047h6lta047h6lwuy8a3",
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
				wallet, found := k.GetWallet(k.Ctx, tc.wallet.WalletAddress)
				require.True(t, found, "wallet should be found")
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
	k := keepertest.NewKeylessKeeper(t)
	srv := keeper.NewMsgServerImpl(k.Keeper)
	wctx := sdk.WrapSDKContext(k.Ctx)

	// Clear store before test
	k.ClearStore()

	// Create test wallet first
	walletAddr := "cosmos1w3jhxap3ta047h6lta047h6lta047h6lx84s66"
	creator := "cosmos1w3jhxap3ta047h6lta047h6lta047h6lx84s66"
	grantee := "cosmos1w3jhxapjta047h6lta047h6lta047h6lwuy8a3"

	msg := &types.MsgCreateWallet{
		Creator:       creator,
		PubKey:        "pubkey1",
		WalletAddress: walletAddr,
		ChainId:       "test-1",
	}
	_, err := srv.CreateWallet(wctx, msg)
	require.NoError(t, err)

	// Set wallet status to active
	err = k.SetWalletStatus(k.Ctx, walletAddr, types.WalletStatus_WALLET_STATUS_ACTIVE)
	require.NoError(t, err)

	// Test granting permission
	expiresAt := time.Now().Add(24 * time.Hour)
	permission := &types.Permission{
		WalletAddress: walletAddr,
		Grantee:      grantee,
		Permissions:  []string{types.WalletPermission_WALLET_PERMISSION_SIGN.String()},
		ExpiresAt:    &expiresAt,
	}

	err = k.GrantPermission(k.Ctx, permission)
	require.NoError(t, err)

	// Test checking permission
	hasPermission, err := k.HasPermission(k.Ctx, walletAddr, grantee, types.WalletPermission_WALLET_PERMISSION_SIGN.String())
	require.NoError(t, err)
	require.True(t, hasPermission)

	// Test checking non-existent permission
	hasPermission, err = k.HasPermission(k.Ctx, walletAddr, grantee, types.WalletPermission_WALLET_PERMISSION_RECOVER.String())
	require.NoError(t, err)
	require.False(t, hasPermission)

	// Test checking permission for non-existent grantee
	hasPermission, err = k.HasPermission(k.Ctx, walletAddr, "cosmos1v9jxgu33kewfvynvl5mu8xg3u2m3ugytkut3hz", types.WalletPermission_WALLET_PERMISSION_SIGN.String())
	require.NoError(t, err)
	require.False(t, hasPermission)

	// Test checking permission for non-existent wallet
	hasPermission, err = k.HasPermission(k.Ctx, "cosmos1v9jxgu33kewfvynvl5mu8xg3u2m3ugytkut3hz", grantee, types.WalletPermission_WALLET_PERMISSION_SIGN.String())
	require.NoError(t, err)
	require.False(t, hasPermission)
}
