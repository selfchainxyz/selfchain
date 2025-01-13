package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "selfchain/testutil/keeper"
	"selfchain/x/keyless/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.KeylessKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
