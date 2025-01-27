package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	testkeeper "selfchain/testutil/keeper"
	"selfchain/x/keyless/types"
)

func TestPermissionsQuery(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	wctx := sdk.WrapSDKContext(k.Ctx)
	testCases := []struct {
		desc     string
		request  *types.QueryPermissionsRequest
		response *types.QueryPermissionsResponse
		err      error
	}{
		{
			desc:    "First",
			request: &types.QueryPermissionsRequest{},
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
		{
			desc: "NoPermissions",
			request: &types.QueryPermissionsRequest{
				WalletId: "test-wallet",
			},
			response: &types.QueryPermissionsResponse{
				Permissions: []*types.Permission{},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.Permissions(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestPermissionQuery(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	wctx := sdk.WrapSDKContext(k.Ctx)
	testCases := []struct {
		desc     string
		request  *types.QueryPermissionRequest
		response *types.QueryPermissionResponse
		err      error
	}{
		{
			desc:    "InvalidRequest",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
		{
			desc: "NotFound",
			request: &types.QueryPermissionRequest{
				WalletId: "test-wallet",
				Grantee:  "test-grantee",
			},
			err: status.Error(codes.NotFound, "permission not found"),
		},
		{
			desc: "NoWalletId",
			request: &types.QueryPermissionRequest{
				Grantee: "test-grantee",
			},
			err: status.Error(codes.InvalidArgument, "wallet ID cannot be empty"),
		},
		{
			desc: "NoGrantee",
			request: &types.QueryPermissionRequest{
				WalletId: "test-wallet",
			},
			err: status.Error(codes.InvalidArgument, "grantee cannot be empty"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.Permission(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestPermissionsQueryPaginated(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	wctx := sdk.WrapSDKContext(k.Ctx)

	// Create test permissions
	testWalletId := "test-wallet"
	for i := 0; i < 5; i++ {
		perm := types.Permission{
			WalletId: testWalletId,
			Grantee:  "grantee-" + string(rune(i+'0')),
			Permissions: []string{
				types.WalletPermission_WALLET_PERMISSION_SIGN.String(),
			},
		}
		k.SetPermission(k.Ctx, &perm)
	}

	request := &types.QueryPermissionsRequest{
		WalletId: testWalletId,
		Pagination: &query.PageRequest{
			Limit: 3,
		},
	}

	response, err := k.Permissions(wctx, request)
	require.NoError(t, err)
	require.NotNil(t, response.Pagination)
	require.Len(t, response.Permissions, 3)
}
