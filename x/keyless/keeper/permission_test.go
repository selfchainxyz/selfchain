package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	testkeeper "selfchain/testutil/keeper"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
)

func TestPermissionValidation(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	srv := keeper.NewMsgServerImpl(k.Keeper)

	// Create test wallet first
	walletAddr := "cosmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v"
	msg := &types.MsgCreateWallet{
		Creator:       "cosmos1s4ycalgh3gjemd4hmqcvcgmnf647rnd0tpg2w9",
		PubKey:        "pubkey1",
		WalletAddress: walletAddr,
		ChainId:      "test-1",
	}
	_, err := srv.CreateWallet(k.Ctx, msg)
	require.NoError(t, err)

	// Set wallet status to active
	wallet, err := k.GetWallet(k.Ctx, walletAddr)
	require.NoError(t, err)
	wallet.Status = types.WalletStatus_WALLET_STATUS_ACTIVE
	err = k.SaveWallet(k.Ctx, wallet)
	require.NoError(t, err)

	expiresAt := time.Now().Add(24 * time.Hour)

	// Test cases
	tests := []struct {
		name        string
		permission  *types.Permission
		expectError bool
		errorMsg    error
	}{
		{
			name: "valid permission",
			permission: &types.Permission{
				WalletAddress: walletAddr,
				Grantee:      "cosmos1v9jxgu33kewfvynvl5mu8xg3u2m3ugytqqpspa",
				Permissions:  []string{"WALLET_PERMISSION_SIGN"},
				ExpiresAt:    &expiresAt,
			},
			expectError: false,
		},
		{
			name: "empty wallet address",
			permission: &types.Permission{
				WalletAddress: "",
				Grantee:      "cosmos1v9jxgu33kewfvynvl5mu8xg3u2m3ugytqqpspa",
				Permissions:  []string{"WALLET_PERMISSION_SIGN"},
				ExpiresAt:    &expiresAt,
			},
			expectError: true,
			errorMsg:    status.Error(codes.InvalidArgument, "wallet address cannot be empty: invalid permission"),
		},
		{
			name: "empty grantee",
			permission: &types.Permission{
				WalletAddress: walletAddr,
				Grantee:      "",
				Permissions:  []string{"WALLET_PERMISSION_SIGN"},
				ExpiresAt:    &expiresAt,
			},
			expectError: true,
			errorMsg:    status.Error(codes.InvalidArgument, "grantee cannot be empty: invalid permission"),
		},
		{
			name: "no permissions",
			permission: &types.Permission{
				WalletAddress: walletAddr,
				Grantee:      "cosmos1v9jxgu33kewfvynvl5mu8xg3u2m3ugytqqpspa",
				Permissions:  []string{},
				ExpiresAt:    &expiresAt,
			},
			expectError: true,
			errorMsg:    status.Error(codes.InvalidArgument, "permissions cannot be empty: invalid permission"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clear store before each test case
			k.ClearStore()

			// Create wallet again after clearing store
			_, err := srv.CreateWallet(k.Ctx, msg)
			require.NoError(t, err)

			// Set wallet status to active
			wallet, err := k.GetWallet(k.Ctx, walletAddr)
			require.NoError(t, err)
			wallet.Status = types.WalletStatus_WALLET_STATUS_ACTIVE
			err = k.SaveWallet(k.Ctx, wallet)
			require.NoError(t, err)

			err = k.GrantPermission(k.Ctx, tc.permission)
			if tc.expectError {
				require.Error(t, err)
				require.Equal(t, tc.errorMsg.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSinglePermission(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)

	// Create test wallet first
	walletAddr := "cosmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v"
	creator := "cosmos1s4ycalgh3gjemd4hmqcvcgmnf647rnd0tpg2w9"

	srv := keeper.NewMsgServerImpl(k.Keeper)
	msg := &types.MsgCreateWallet{
		Creator:       creator,
		PubKey:        "pubkey1",
		WalletAddress: walletAddr,
		ChainId:       "test-1",
	}
	_, err := srv.CreateWallet(k.Ctx, msg)
	require.NoError(t, err)

	// Set wallet status to active
	wallet, err := k.GetWallet(k.Ctx, walletAddr)
	require.NoError(t, err)
	wallet.Status = types.WalletStatus_WALLET_STATUS_ACTIVE
	err = k.SaveWallet(k.Ctx, wallet)
	require.NoError(t, err)

	expiresAt := time.Now().Add(24 * time.Hour)

	// Grant permission
	permission := &types.Permission{
		WalletAddress: walletAddr,
		Grantee:      "cosmos1grantee",
		Permissions:  []string{"WALLET_PERMISSION_SIGN"},
		ExpiresAt:    &expiresAt,
	}

	err = k.GrantPermission(k.Ctx, permission)
	require.NoError(t, err)

	// Check permission is valid
	isValid := k.Keeper.IsWalletAuthorized(k.Ctx, walletAddr, "cosmos1grantee", types.WalletPermission_WALLET_PERMISSION_SIGN)
	require.True(t, isValid)
}

func TestMultiplePermissions(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)

	// Create test wallet first
	walletAddr := "cosmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v"
	creator := "cosmos1s4ycalgh3gjemd4hmqcvcgmnf647rnd0tpg2w9"

	srv := keeper.NewMsgServerImpl(k.Keeper)
	msg := &types.MsgCreateWallet{
		Creator:       creator,
		PubKey:        "pubkey1",
		WalletAddress: walletAddr,
		ChainId:       "test-1",
	}
	_, err := srv.CreateWallet(k.Ctx, msg)
	require.NoError(t, err)

	// Set wallet status to active
	wallet, err := k.GetWallet(k.Ctx, walletAddr)
	require.NoError(t, err)
	wallet.Status = types.WalletStatus_WALLET_STATUS_ACTIVE
	err = k.SaveWallet(k.Ctx, wallet)
	require.NoError(t, err)

	expiresAt := time.Now().Add(24 * time.Hour)

	// Grant permission
	permission := &types.Permission{
		WalletAddress: walletAddr,
		Grantee:      "cosmos1grantee",
		Permissions:  []string{"WALLET_PERMISSION_SIGN", "WALLET_PERMISSION_ADMIN"},
		ExpiresAt:    &expiresAt,
	}

	err = k.GrantPermission(k.Ctx, permission)
	require.NoError(t, err)

	// Check both permissions are valid
	isValid := k.Keeper.IsWalletAuthorized(k.Ctx, walletAddr, "cosmos1grantee", types.WalletPermission_WALLET_PERMISSION_SIGN)
	require.True(t, isValid)

	isValid = k.Keeper.IsWalletAuthorized(k.Ctx, walletAddr, "cosmos1grantee", types.WalletPermission_WALLET_PERMISSION_ADMIN)
	require.True(t, isValid)
}

func TestPermissionGrantAndRevoke(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	srv := keeper.NewMsgServerImpl(k.Keeper)

	// Create a test wallet first
	walletAddr := "cosmos1test"
	msg := &types.MsgCreateWallet{
		Creator:       "cosmos1owner",
		PubKey:        "pubkey1",
		WalletAddress: walletAddr,
		ChainId:      "test-1",
	}
	_, err := srv.CreateWallet(k.Ctx, msg)
	require.NoError(t, err)

	expiresAt := time.Now().Add(24 * time.Hour)

	// Grant permission
	permission := &types.Permission{
		WalletAddress: walletAddr,
		Grantee:      "cosmos1grantee",
		Permissions:  []string{"WALLET_PERMISSION_SIGN"},
		ExpiresAt:    &expiresAt,
	}

	err = k.Keeper.GrantPermission(k.Ctx, permission)
	require.NoError(t, err)

	// Verify permission exists
	perm, err := k.Keeper.GetPermission(k.Ctx, walletAddr, "cosmos1grantee")
	require.NoError(t, err)
	require.Equal(t, permission.Permissions, perm.Permissions)

	// Revoke permission
	err = k.Keeper.RevokePermission(k.Ctx, walletAddr, "cosmos1grantee")
	require.NoError(t, err)

	// Verify permission is revoked
	perm, err = k.Keeper.GetPermission(k.Ctx, walletAddr, "cosmos1grantee")
	require.NoError(t, err)
	require.True(t, perm.Revoked)
}

func TestGetPermissionsForWallet(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	srv := keeper.NewMsgServerImpl(k.Keeper)

	// Create a test wallet first
	walletAddr := "cosmos1test"
	msg := &types.MsgCreateWallet{
		Creator:       "cosmos1owner",
		PubKey:        "pubkey1",
		WalletAddress: walletAddr,
		ChainId:      "test-1",
	}
	_, err := srv.CreateWallet(k.Ctx, msg)
	require.NoError(t, err)

	expiresAt := time.Now().Add(24 * time.Hour)

	// Grant multiple permissions
	permissions := []*types.Permission{
		{
			WalletAddress: walletAddr,
			Grantee:      "cosmos1grantee1",
			Permissions:  []string{"WALLET_PERMISSION_SIGN"},
			ExpiresAt:    &expiresAt,
		},
		{
			WalletAddress: walletAddr,
			Grantee:      "cosmos1grantee2",
			Permissions:  []string{"WALLET_PERMISSION_RECOVER"},
			ExpiresAt:    &expiresAt,
		},
	}

	for _, perm := range permissions {
		err := k.Keeper.GrantPermission(k.Ctx, perm)
		require.NoError(t, err)
	}

	// Get all permissions for wallet
	perms, err := k.Keeper.GetPermissionsForWallet(k.Ctx, walletAddr)
	require.NoError(t, err)
	require.Len(t, perms, 2)
}

func TestPermissionExpiration(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	srv := keeper.NewMsgServerImpl(k.Keeper)

	// Create a test wallet first
	walletAddr := "cosmos1test"
	msg := &types.MsgCreateWallet{
		Creator:       "cosmos1owner",
		PubKey:        "pubkey1",
		WalletAddress: walletAddr,
		ChainId:      "test-1",
	}
	_, err := srv.CreateWallet(k.Ctx, msg)
	require.NoError(t, err)

	expiresAt := time.Now().Add(24 * time.Hour)

	// Grant permission
	permission := &types.Permission{
		WalletAddress: walletAddr,
		Grantee:      "cosmos1grantee",
		Permissions:  []string{"WALLET_PERMISSION_SIGN"},
		ExpiresAt:    &expiresAt,
	}

	err = k.Keeper.GrantPermission(k.Ctx, permission)
	require.NoError(t, err)

	// Check permission is valid
	isValid := k.Keeper.IsWalletAuthorized(k.Ctx, walletAddr, "cosmos1grantee", types.WalletPermission_WALLET_PERMISSION_SIGN)
	require.True(t, isValid)
}
