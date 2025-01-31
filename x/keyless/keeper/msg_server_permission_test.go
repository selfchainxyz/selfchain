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

func TestMsgServerPermission(t *testing.T) {
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
	grantMsg := &types.MsgGrantPermission{
		Creator:       creator,
		WalletAddress: walletAddr,
		Grantee:      grantee,
		Permissions:  []types.WalletPermission{types.WalletPermission_WALLET_PERMISSION_SIGN},
		ExpiresAt:    &expiresAt,
	}

	_, err = srv.GrantPermission(wctx, grantMsg)
	require.NoError(t, err)

	// Verify permission exists
	perms, err := k.GetPermissionsForWallet(k.Ctx, walletAddr)
	require.NoError(t, err)
	require.Len(t, perms, 1)
	require.Equal(t, grantee, perms[0].Grantee)
	require.Equal(t, []string{types.WalletPermission_WALLET_PERMISSION_SIGN.String()}, perms[0].Permissions)

	// Test revoking permission
	revokeMsg := &types.MsgRevokePermission{
		Creator:       creator,
		WalletAddress: walletAddr,
		Grantee:      grantee,
		Permissions:  []types.WalletPermission{types.WalletPermission_WALLET_PERMISSION_SIGN},
	}

	_, err = srv.RevokePermission(wctx, revokeMsg)
	require.NoError(t, err)

	// Verify permission is revoked
	perms, err = k.GetPermissionsForWallet(k.Ctx, walletAddr)
	require.NoError(t, err)
	require.Empty(t, perms)

	// Test granting permission with invalid grantee
	grantMsg.Grantee = "invalid"
	_, err = srv.GrantPermission(wctx, grantMsg)
	require.Error(t, err)

	// Test revoking permission with invalid grantee
	revokeMsg.Grantee = "invalid"
	_, err = srv.RevokePermission(wctx, revokeMsg)
	require.Error(t, err)
}
