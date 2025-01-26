package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	testkeeper "selfchain/testutil/keeper"
	"selfchain/x/keyless/types"
)

func TestParamsQuery(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	wctx := sdk.WrapSDKContext(k.Ctx)
	params := types.DefaultParams()
	k.SetParams(k.Ctx, params)

	response, err := k.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
