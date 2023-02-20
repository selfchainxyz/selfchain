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

func createNTokenMigration(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.TokenMigration {
	items := make([]types.TokenMigration, n)
	for i := range items {
		items[i].MsgHash = strconv.Itoa(i)

		keeper.SetTokenMigration(ctx, items[i])
	}
	return items
}

func TestTokenMigrationGet(t *testing.T) {
	keeper, ctx := keepertest.MigrationKeeper(t)
	items := createNTokenMigration(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetTokenMigration(ctx,
			item.MsgHash,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestTokenMigrationGetAll(t *testing.T) {
	keeper, ctx := keepertest.MigrationKeeper(t)
	items := createNTokenMigration(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllTokenMigration(ctx)),
	)
}
