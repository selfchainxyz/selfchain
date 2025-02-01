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
	walletAddr := "cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzf439"
	creator := "cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzf439"
	grantee := "cosmos1nxv42u3lv642q0fuzu2qmrku27zgut3n3z7lll"
	nonExistentWallet := "cosmos1nxv42u3lv642q0fuzu2qmrku27zgut3n3z7lll"

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

	// Test granting permission with empty permissions
	expiresAt := time.Now().Add(24 * time.Hour)
	permission := &types.Permission{
		WalletAddress: walletAddr,
		Grantee:      grantee,
		Permissions:  []string{},
		ExpiresAt:    &expiresAt,
	}
	err = k.GrantPermission(k.Ctx, permission)
	require.Error(t, err)
	require.Equal(t, types.ErrEmptyPermissions, err)

	// Test granting permission with empty wallet address
	permission.Permissions = []string{types.WalletPermission_WALLET_PERMISSION_SIGN.String()}
	permission.WalletAddress = ""
	err = k.GrantPermission(k.Ctx, permission)
	require.Error(t, err)
	require.Contains(t, err.Error(), "wallet address cannot be empty")

	// Test granting permission with empty grantee
	permission.WalletAddress = walletAddr
	permission.Grantee = ""
	err = k.GrantPermission(k.Ctx, permission)
	require.Error(t, err)
	require.Contains(t, err.Error(), "grantee cannot be empty")

	// Test granting permission with invalid wallet address
	permission.WalletAddress = "invalid"
	permission.Grantee = grantee
	err = k.GrantPermission(k.Ctx, permission)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid bech32 address")

	// Test granting permission with invalid grantee
	permission.WalletAddress = walletAddr
	permission.Grantee = "invalid"
	err = k.GrantPermission(k.Ctx, permission)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid bech32 address")

	// Test granting permission with non-existent wallet
	permission.WalletAddress = nonExistentWallet
	permission.Grantee = grantee
	err = k.GrantPermission(k.Ctx, permission)
	require.Error(t, err)
	require.Contains(t, err.Error(), "wallet not found")

	// Test successful permission grant
	permission.WalletAddress = walletAddr
	permission.Grantee = grantee
	err = k.GrantPermission(k.Ctx, permission)
	require.NoError(t, err)

	// Verify permission exists
	perms, err := k.GetPermissionsForWallet(k.Ctx, walletAddr)
	require.NoError(t, err)
	require.Len(t, perms, 1)
	require.Equal(t, grantee, perms[0].Grantee)
	require.Equal(t, []string{types.WalletPermission_WALLET_PERMISSION_SIGN.String()}, perms[0].Permissions)

	// Test revoking permission with empty wallet address
	err = k.RevokePermission(k.Ctx, "", grantee)
	require.Error(t, err)
	require.Contains(t, err.Error(), "wallet address cannot be empty")

	// Test revoking permission with empty grantee
	err = k.RevokePermission(k.Ctx, walletAddr, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "grantee cannot be empty")

	// Test revoking permission with invalid wallet address
	err = k.RevokePermission(k.Ctx, "invalid", grantee)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid bech32 address")

	// Test revoking permission with invalid grantee
	err = k.RevokePermission(k.Ctx, walletAddr, "invalid")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid bech32 address")

	// Test revoking permission with non-existent wallet
	err = k.RevokePermission(k.Ctx, nonExistentWallet, grantee)
	require.Error(t, err)
	require.Contains(t, err.Error(), "wallet not found")

	// Test successful permission revocation
	err = k.RevokePermission(k.Ctx, walletAddr, grantee)
	require.NoError(t, err)

	// Verify permission is revoked
	perms, err = k.GetPermissionsForWallet(k.Ctx, walletAddr)
	require.NoError(t, err)
	require.Empty(t, perms)
}
