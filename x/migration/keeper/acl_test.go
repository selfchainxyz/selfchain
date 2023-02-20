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

func createTestAcl(keeper *keeper.Keeper, ctx sdk.Context) types.Acl {
	item := types.Acl{}
	keeper.SetAcl(ctx, item)
	return item
}

func TestAclGet(t *testing.T) {
	keeper, ctx := keepertest.MigrationKeeper(t)
	item := createTestAcl(keeper, ctx)
	rst, found := keeper.GetAcl(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}
