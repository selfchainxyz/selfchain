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

func TestPermissionQuery(t *testing.T) {
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

	// Test get permission
	perm, found := k.GetPermission(k.Ctx, walletAddr, grantee)
	require.True(t, found)
	require.NotNil(t, perm)
	require.Equal(t, grantee, perm.Grantee)
	require.Equal(t, []string{types.WalletPermission_WALLET_PERMISSION_SIGN.String()}, perm.Permissions)

	// Test get permission for non-existent grantee
	perm, found = k.GetPermission(k.Ctx, walletAddr, "cosmos1w3jhxapnta047h6lta047h6lta047h6l34t280")
	require.False(t, found)
	require.Nil(t, perm)

	// Test querying all permissions for wallet
	perms, err := k.GetPermissionsForWallet(k.Ctx, walletAddr)
	require.NoError(t, err)
	require.Len(t, perms, 1)
	require.Equal(t, grantee, perms[0].Grantee)
	require.Equal(t, []string{types.WalletPermission_WALLET_PERMISSION_SIGN.String()}, perms[0].Permissions)

	// Test querying all permissions for non-existent wallet
	perms, err = k.GetPermissionsForWallet(k.Ctx, "cosmos1w3jhxapnta047h6lta047h6lta047h6l34t280")
	require.NoError(t, err)
	require.Len(t, perms, 0)
}
