package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/testutil"
	"selfchain/x/keyless/types"
)

func TestMsgServerGrantPermission(t *testing.T) {
	k, ctx := testutil.NewKeeper(t)
	server := keeper.NewMsgServerImpl(*k)

	// Create test wallet
	wallet := &types.Wallet{
		Id:            "test_wallet_1",
		Creator:       "creator1",
		WalletAddress: "wallet1",
		ChainId:      "test-chain-1",
		Status:       types.WalletStatus_WALLET_STATUS_ACTIVE,
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
					types.WalletPermission_WALLET_PERMISSION_RECOVER,
				},
			},
			expectError: false,
		},
		{
			name: "unauthorized creator",
			msg: &types.MsgGrantPermission{
				Creator:  "creator2",
				WalletId: "test_wallet_1",
				Grantee:  "grantee1",
				Permissions: []types.WalletPermission{
					types.WalletPermission_WALLET_PERMISSION_SIGN,
				},
			},
			expectError: true,
		},
		{
			name: "wallet not found",
			msg: &types.MsgGrantPermission{
				Creator:  "creator1",
				WalletId: "nonexistent_wallet",
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
			_, err := server.GrantPermission(sdk.WrapSDKContext(ctx), tc.msg)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgServerRevokePermission(t *testing.T) {
	k, ctx := testutil.NewKeeper(t)
	server := keeper.NewMsgServerImpl(*k)

	// Create test wallet
	wallet := &types.Wallet{
		Id:            "test_wallet_1",
		Creator:       "creator1",
		WalletAddress: "wallet1",
		ChainId:      "test-chain-1",
		Status:       types.WalletStatus_WALLET_STATUS_ACTIVE,
	}
	err := k.SaveWallet(ctx, wallet)
	require.NoError(t, err)

	// Grant initial permission
	grant := &types.MsgGrantPermission{
		Creator:  "creator1",
		WalletId: "test_wallet_1",
		Grantee:  "grantee1",
		Permissions: []types.WalletPermission{
			types.WalletPermission_WALLET_PERMISSION_SIGN,
			types.WalletPermission_WALLET_PERMISSION_RECOVER,
		},
	}
	_, err = server.GrantPermission(sdk.WrapSDKContext(ctx), grant)
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
			name: "unauthorized creator",
			msg: &types.MsgRevokePermission{
				Creator:  "creator2",
				WalletId: "test_wallet_1",
				Grantee:  "grantee1",
				Permissions: []types.WalletPermission{
					types.WalletPermission_WALLET_PERMISSION_SIGN,
				},
			},
			expectError: true,
		},
		{
			name: "permission not found",
			msg: &types.MsgRevokePermission{
				Creator:  "creator1",
				WalletId: "test_wallet_1",
				Grantee:  "grantee2",
				Permissions: []types.WalletPermission{
					types.WalletPermission_WALLET_PERMISSION_SIGN,
				},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := server.RevokePermission(sdk.WrapSDKContext(ctx), tc.msg)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
