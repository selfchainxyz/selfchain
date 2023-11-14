package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "selfchain/testutil/keeper"
	"selfchain/testutil/nullify"
	"selfchain/x/migration/keeper"
	"selfchain/x/migration/types"
)

func createTestConfig(keeper *keeper.Keeper, ctx sdk.Context) types.Config {
	item := types.Config{}
	keeper.SetConfig(ctx, item)
	return item
}

func TestConfigGet(t *testing.T) {
	keeper, ctx := keepertest.MigrationKeeper(t)
	item := createTestConfig(keeper, ctx)
	rst, found := keeper.GetConfig(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestConfigRemove(t *testing.T) {
	keeper, ctx := keepertest.MigrationKeeper(t)
	createTestConfig(keeper, ctx)
	keeper.RemoveConfig(ctx)
	_, found := keeper.GetConfig(ctx)
	require.False(t, found)
}
