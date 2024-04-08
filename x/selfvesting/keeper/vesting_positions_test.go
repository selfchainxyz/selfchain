package keeper_test

import (
	"strconv"
	"testing"

	keepertest "selfchain/testutil/keeper"
	"selfchain/testutil/nullify"
	"selfchain/x/selfvesting/keeper"
	"selfchain/x/selfvesting/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNVestingPositions(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.VestingPositions {
	items := make([]types.VestingPositions, n)
	for i := range items {
		items[i].Beneficiary = strconv.Itoa(i)

		keeper.SetVestingPositions(ctx, items[i])
	}
	return items
}

func TestVestingPositionsGet(t *testing.T) {
	keeper, ctx := keepertest.SelfvestingKeeper(t)
	items := createNVestingPositions(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetVestingPositions(ctx,
			item.Beneficiary,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestVestingPositionsGetAll(t *testing.T) {
	keeper, ctx := keepertest.SelfvestingKeeper(t)
	items := createNVestingPositions(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllVestingPositions(ctx)),
	)
}
