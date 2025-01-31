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

func TestPermissionManagement(t *testing.T) {
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

	// Verify permission exists
	perms, err := k.GetPermissionsForWallet(k.Ctx, walletAddr)
	require.NoError(t, err)
	require.Len(t, perms, 1)
	require.Equal(t, grantee, perms[0].Grantee)
	require.Equal(t, []string{types.WalletPermission_WALLET_PERMISSION_SIGN.String()}, perms[0].Permissions)

	// Test revoking permission
	err = k.RevokePermission(k.Ctx, walletAddr, grantee)
	require.NoError(t, err)

	// Verify permission is revoked
	perms, err = k.GetPermissionsForWallet(k.Ctx, walletAddr)
	require.NoError(t, err)
	require.Empty(t, perms)

	// Test granting permission with invalid grantee
	permission.Grantee = "invalid"
	err = k.GrantPermission(k.Ctx, permission)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid bech32 address")

	// Test revoking permission with invalid grantee
	err = k.RevokePermission(k.Ctx, walletAddr, "invalid")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid bech32 address")

	// Test granting permission with invalid wallet address
	permission.WalletAddress = "invalid"
	err = k.GrantPermission(k.Ctx, permission)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid bech32 address")

	// Test revoking permission with invalid wallet address
	err = k.RevokePermission(k.Ctx, "invalid", grantee)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid bech32 address")

	// Test granting permission with empty permissions
	permission.WalletAddress = walletAddr
	permission.Grantee = grantee
	permission.Permissions = []string{}
	err = k.GrantPermission(k.Ctx, permission)
	require.Error(t, err)
	require.Contains(t, err.Error(), "permissions cannot be empty")
}
