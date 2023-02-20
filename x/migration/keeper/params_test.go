package keeper_test

import (
	"testing"

	testkeeper "selfchain/testutil/keeper"
	"selfchain/x/migration/types"

	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.MigrationKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
