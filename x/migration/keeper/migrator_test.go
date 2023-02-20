package keeper_test

import (
	"strconv"
	"testing"

	keepertest "selfchain/testutil/keeper"
	"selfchain/testutil/nullify"
	"selfchain/x/migration/keeper"
	"selfchain/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNMigrator(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Migrator {
	items := make([]types.Migrator, n)
	for i := range items {
		items[i].Migrator = strconv.Itoa(i)

		keeper.SetMigrator(ctx, items[i])
	}
	return items
}

func TestMigratorGet(t *testing.T) {
	keeper, ctx := keepertest.MigrationKeeper(t)
	items := createNMigrator(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetMigrator(ctx,
			item.Migrator,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestMigratorRemove(t *testing.T) {
	keeper, ctx := keepertest.MigrationKeeper(t)
	items := createNMigrator(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveMigrator(ctx,
			item.Migrator,
		)
		_, found := keeper.GetMigrator(ctx,
			item.Migrator,
		)
		require.False(t, found)
	}
}

func TestMigratorGetAll(t *testing.T) {
	keeper, ctx := keepertest.MigrationKeeper(t)
	items := createNMigrator(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllMigrator(ctx)),
	)
}
