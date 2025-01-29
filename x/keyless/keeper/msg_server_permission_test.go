package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "selfchain/testutil/keeper"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
)

func TestMsgServerPermission(t *testing.T) {
	k := keepertest.NewKeylessKeeper(t)
	srv := keeper.NewMsgServerImpl(k.Keeper)
	wctx := sdk.WrapSDKContext(k.Ctx)

	// Create test wallet first
	walletAddr := "cosmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v"
	creator := "cosmos1s4ycalgh3gjemd4hmqcvcgmnf647rnd0tpg2w9"
	msg := &types.MsgCreateWallet{
		Creator:       creator,
		PubKey:        "pubkey1",
		WalletAddress: walletAddr,
		ChainId:       "test-1",
	}
	_, err := srv.CreateWallet(wctx, msg)
	require.NoError(t, err)

	// Set wallet status to active
	wallet, err := k.GetWallet(k.Ctx, walletAddr)
	require.NoError(t, err)
	wallet.Status = types.WalletStatus_WALLET_STATUS_ACTIVE
	err = k.SaveWallet(k.Ctx, wallet)
	require.NoError(t, err)

	expiresAt := time.Now().Add(24 * time.Hour)

	t.Run("GrantPermission", func(t *testing.T) {
		tests := []struct {
			name    string
			msg     *types.MsgGrantPermission
			err     error
		}{
			{
				name: "Valid Grant",
				msg: &types.MsgGrantPermission{
					Creator:       creator,
					WalletAddress: walletAddr,
					Grantee:      "cosmos1v9jxgu33kewfvynvl5mu8xg3u2m3ugytqqpspa",
					Permissions:  []types.WalletPermission{types.WalletPermission_WALLET_PERMISSION_SIGN, types.WalletPermission_WALLET_PERMISSION_ADMIN},
					ExpiresAt:    &expiresAt,
				},
			},
			{
				name: "Invalid Creator",
				msg: &types.MsgGrantPermission{
					Creator:       "invalid",
					WalletAddress: walletAddr,
					Grantee:      "cosmos1v9jxgu33kewfvynvl5mu8xg3u2m3ugytqqpspa",
					Permissions:  []types.WalletPermission{types.WalletPermission_WALLET_PERMISSION_SIGN},
					ExpiresAt:    &expiresAt,
				},
				err: status.Error(codes.InvalidArgument, "invalid creator address (decoding bech32 failed: invalid bech32 string length 7): invalid address"),
			},
			{
				name: "Invalid Wallet",
				msg: &types.MsgGrantPermission{
					Creator:       creator,
					WalletAddress: "invalid",
					Grantee:      "cosmos1v9jxgu33kewfvynvl5mu8xg3u2m3ugytqqpspa",
					Permissions:  []types.WalletPermission{types.WalletPermission_WALLET_PERMISSION_SIGN},
					ExpiresAt:    &expiresAt,
				},
				err: status.Error(codes.InvalidArgument, "invalid wallet address (decoding bech32 failed: invalid bech32 string length 7): invalid address"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Clear store before each test case
				k.ClearStore()

				// Create wallet again after clearing store
				_, err := srv.CreateWallet(wctx, msg)
				require.NoError(t, err)

				// Set wallet status to active
				wallet, err := k.GetWallet(k.Ctx, walletAddr)
				require.NoError(t, err)
				wallet.Status = types.WalletStatus_WALLET_STATUS_ACTIVE
				err = k.SaveWallet(k.Ctx, wallet)
				require.NoError(t, err)

				_, err = srv.GrantPermission(wctx, tt.msg)
				if tt.err != nil {
					require.Error(t, err)
					require.Equal(t, tt.err.Error(), err.Error())
					return
				}
				require.NoError(t, err)

				// Verify permission was granted
				perm, err := k.GetPermission(k.Ctx, tt.msg.WalletAddress, tt.msg.Grantee)
				require.NoError(t, err)
				require.NotNil(t, perm)
				require.Equal(t, tt.msg.Permissions, perm.Permissions)
			})
		}
	})

	t.Run("RevokePermission", func(t *testing.T) {
		// Clear store before revoke tests
		k.ClearStore()

		// Create wallet again after clearing store
		_, err := srv.CreateWallet(wctx, msg)
		require.NoError(t, err)

		// Set wallet status to active
		wallet, err := k.GetWallet(k.Ctx, walletAddr)
		require.NoError(t, err)
		wallet.Status = types.WalletStatus_WALLET_STATUS_ACTIVE
		err = k.SaveWallet(k.Ctx, wallet)
		require.NoError(t, err)

		// First grant a permission
		grantMsg := &types.MsgGrantPermission{
			Creator:       creator,
			WalletAddress: walletAddr,
			Grantee:      "cosmos1v9jxgu33kewfvynvl5mu8xg3u2m3ugytqqpspa",
			Permissions:  []types.WalletPermission{types.WalletPermission_WALLET_PERMISSION_SIGN},
			ExpiresAt:    &expiresAt,
		}
		_, err = srv.GrantPermission(wctx, grantMsg)
		require.NoError(t, err)

		tests := []struct {
			name    string
			msg     *types.MsgRevokePermission
			err     error
		}{
			{
				name: "Valid Revoke",
				msg: &types.MsgRevokePermission{
					Creator:       creator,
					WalletAddress: walletAddr,
					Grantee:      "cosmos1v9jxgu33kewfvynvl5mu8xg3u2m3ugytqqpspa",
				},
			},
			{
				name: "Invalid Creator",
				msg: &types.MsgRevokePermission{
					Creator:       "invalid",
					WalletAddress: walletAddr,
					Grantee:      "cosmos1v9jxgu33kewfvynvl5mu8xg3u2m3ugytqqpspa",
				},
				err: status.Error(codes.InvalidArgument, "invalid creator address (decoding bech32 failed: invalid bech32 string length 7): invalid address"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := srv.RevokePermission(wctx, tt.msg)
				if tt.err != nil {
					require.Error(t, err)
					require.Equal(t, tt.err.Error(), err.Error())
					return
				}
				require.NoError(t, err)

				// Verify permission was revoked
				perm, err := k.GetPermission(k.Ctx, tt.msg.WalletAddress, tt.msg.Grantee)
				require.NoError(t, err)
				require.True(t, perm.Revoked)
			})
		}
	})
}
