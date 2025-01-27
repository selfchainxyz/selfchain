package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	testkeeper "selfchain/testutil/keeper"
	"selfchain/x/keyless/types"
)

func TestPermissionGrant(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	ctx := k.Ctx

	// Create test wallet
	wallet := &types.Wallet{
		Id:            "test_wallet_1",
		Creator:       "creator1",
		WalletAddress: "wallet1",
		ChainId:       "test-chain-1",
		Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
	}
	err := k.SaveWallet(ctx, wallet)
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name        string
		msg         *types.MsgGrantPermission
		expectError bool
	}{
		{
			name: "valid permission grant",
			msg: &types.MsgGrantPermission{
				Creator:  "creator1",
				WalletId: "test_wallet_1",
				Grantee:  "grantee1",
				Permissions: []types.WalletPermission{
					types.WalletPermission_WALLET_PERMISSION_SIGN,
				},
			},
			expectError: false,
		},
		{
			name: "invalid wallet",
			msg: &types.MsgGrantPermission{
				Creator:  "creator1",
				WalletId: "invalid_wallet",
				Grantee:  "grantee1",
				Permissions: []types.WalletPermission{
					types.WalletPermission_WALLET_PERMISSION_SIGN,
				},
			},
			expectError: true,
		},
		{
			name: "unauthorized creator",
			msg: &types.MsgGrantPermission{
				Creator:  "unauthorized",
				WalletId: "test_wallet_1",
				Grantee:  "grantee1",
				Permissions: []types.WalletPermission{
					types.WalletPermission_WALLET_PERMISSION_SIGN,
				},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			perm, err := k.GrantPermission(ctx, tc.msg)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, perm)
				// Verify permission was granted
				storedPerm, err := k.GetPermission(ctx, tc.msg.WalletId, tc.msg.Grantee)
				require.NoError(t, err)
				require.Equal(t, types.WalletPermission_name[int32(tc.msg.Permissions[0])], storedPerm.Permissions[0])
			}
		})
	}
}

func TestPermissionRevoke(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	ctx := k.Ctx

	// Create test wallet and grant initial permission
	wallet := &types.Wallet{
		Id:            "test_wallet_1",
		Creator:       "creator1",
		WalletAddress: "wallet1",
		ChainId:       "test-chain-1",
		Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
	}
	err := k.SaveWallet(ctx, wallet)
	require.NoError(t, err)

	// Grant initial permission
	grantMsg := &types.MsgGrantPermission{
		Creator:  "creator1",
		WalletId: "test_wallet_1",
		Grantee:  "grantee1",
		Permissions: []types.WalletPermission{
			types.WalletPermission_WALLET_PERMISSION_SIGN,
		},
	}
	_, err = k.GrantPermission(ctx, grantMsg)
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name        string
		msg         *types.MsgRevokePermission
		expectError bool
	}{
		{
			name: "valid permission revoke",
			msg: &types.MsgRevokePermission{
				Creator:  "creator1",
				WalletId: "test_wallet_1",
				Grantee:  "grantee1",
				Permissions: []types.WalletPermission{
					types.WalletPermission_WALLET_PERMISSION_SIGN,
				},
			},
			expectError: false,
		},
		{
			name: "invalid wallet",
			msg: &types.MsgRevokePermission{
				Creator:  "creator1",
				WalletId: "invalid_wallet",
				Grantee:  "grantee1",
				Permissions: []types.WalletPermission{
					types.WalletPermission_WALLET_PERMISSION_SIGN,
				},
			},
			expectError: true,
		},
		{
			name: "unauthorized creator",
			msg: &types.MsgRevokePermission{
				Creator:  "unauthorized",
				WalletId: "test_wallet_1",
				Grantee:  "grantee1",
				Permissions: []types.WalletPermission{
					types.WalletPermission_WALLET_PERMISSION_SIGN,
				},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := k.RevokePermission(ctx, tc.msg)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Verify permission was revoked
				perm, err := k.GetPermission(ctx, tc.msg.WalletId, tc.msg.Grantee)
				require.NoError(t, err)
				require.True(t, perm.Revoked)
				require.NotNil(t, perm.RevokedAt)
			}
		})
	}
}

func TestPermissionValidation(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)

	// Test cases
	tests := []struct {
		name        string
		permissions []types.WalletPermission
		expectError bool
	}{
		{
			name: "valid permissions",
			permissions: []types.WalletPermission{
				types.WalletPermission_WALLET_PERMISSION_SIGN,
				types.WalletPermission_WALLET_PERMISSION_RECOVER,
			},
			expectError: false,
		},
		{
			name: "invalid permission - unspecified",
			permissions: []types.WalletPermission{
				types.WalletPermission_WALLET_PERMISSION_UNSPECIFIED,
			},
			expectError: true,
		},
		{
			name:        "empty permissions",
			permissions: []types.WalletPermission{},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := k.ValidatePermissions(tc.permissions)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPermissionExpiry(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	ctx := k.Ctx

	// Create test wallet
	wallet := &types.Wallet{
		Id:            "test_wallet_1",
		Creator:       "creator1",
		WalletAddress: "wallet1",
		ChainId:       "test-chain-1",
		Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
	}
	err := k.SaveWallet(ctx, wallet)
	require.NoError(t, err)

	// Set expiry time
	expiryTime := time.Now().Add(time.Hour * 24)

	// Grant permission with expiry
	grantMsg := &types.MsgGrantPermission{
		Creator:  "creator1",
		WalletId: "test_wallet_1",
		Grantee:  "grantee1",
		Permissions: []types.WalletPermission{
			types.WalletPermission_WALLET_PERMISSION_SIGN,
		},
		ExpiresAt: &expiryTime,
	}
	_, err = k.GrantPermission(ctx, grantMsg)
	require.NoError(t, err)

	// Test permission before expiry
	hasPermission := k.HasPermission(ctx, "test_wallet_1", "grantee1", types.WalletPermission_WALLET_PERMISSION_SIGN)
	require.True(t, hasPermission)

	// Test permission after expiry
	expiredTime := expiryTime.Add(time.Hour * 48)
	ctx = ctx.WithBlockTime(expiredTime)
	hasPermission = k.HasPermission(ctx, "test_wallet_1", "grantee1", types.WalletPermission_WALLET_PERMISSION_SIGN)
	require.False(t, hasPermission)
}
