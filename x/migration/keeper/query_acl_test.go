package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "selfchain/testutil/keeper"
	"selfchain/testutil/nullify"
	"selfchain/x/migration/types"
)

func TestAclQuery(t *testing.T) {
	keeper, ctx := keepertest.MigrationKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	item := createTestAcl(keeper, ctx)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetAclRequest
		response *types.QueryGetAclResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetAclRequest{},
			response: &types.QueryGetAclResponse{Acl: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Acl(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
