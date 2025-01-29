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

func TestPermissionQuery(t *testing.T) {
	k := keepertest.NewKeylessKeeper(t)
	srv := keeper.NewMsgServerImpl(k.Keeper)
	wctx := sdk.WrapSDKContext(k.Ctx)

	// Clear store before test
	k.ClearStore()

	// Create test wallet first
	walletAddr := "cosmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v"
	creator := "cosmos1s4ycalgh3gjemd4hmqcvcgmnf647rnd0tpg2w9"
	grantee1 := "cosmos1v9jxgu33kewfvynvl5mu8xg3u2m3ugytqqpspa"
	grantee2 := "cosmos1x2w87cvt5mqjncav4lxy8yfreynn273xn5335v"

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

	// Grant permissions to two different grantees
	perm1 := &types.Permission{
		WalletAddress: walletAddr,
		Grantee:      grantee1,
		Permissions:  []string{"WALLET_PERMISSION_SIGN"},
		ExpiresAt:    &expiresAt,
	}
	err = k.GrantPermission(k.Ctx, perm1)
	require.NoError(t, err)

	perm2 := &types.Permission{
		WalletAddress: walletAddr,
		Grantee:      grantee2,
		Permissions:  []string{"WALLET_PERMISSION_ADMIN"},
		ExpiresAt:    &expiresAt,
	}
	err = k.GrantPermission(k.Ctx, perm2)
	require.NoError(t, err)

	t.Run("Query_Permissions", func(t *testing.T) {
		tests := []struct {
			name    string
			req     *types.QueryPermissionsRequest
			err     error
			perms   []*types.Permission
		}{
			{
				name: "Valid Query",
				req: &types.QueryPermissionsRequest{
					WalletId: walletAddr,
				},
				perms: []*types.Permission{perm1, perm2},
			},
			{
				name: "Invalid Wallet",
				req: &types.QueryPermissionsRequest{
					WalletId: "invalid",
				},
				err: status.Error(codes.InvalidArgument, "invalid wallet address (decoding bech32 failed: invalid bech32 string length 7): invalid address"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := k.Permissions(wctx, tt.req)
				if tt.err != nil {
					require.Error(t, err)
					require.Equal(t, tt.err.Error(), err.Error())
					return
				}
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Len(t, resp.Permissions, len(tt.perms))
				for i, perm := range resp.Permissions {
					require.Equal(t, tt.perms[i].WalletAddress, perm.WalletAddress)
					require.Equal(t, tt.perms[i].Grantee, perm.Grantee)
					require.Equal(t, tt.perms[i].Permissions, perm.Permissions)
				}
			})
		}
	})
}
