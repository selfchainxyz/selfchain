package keeper_test

import (
	"context"
	"testing"

	commontest "frontier/testutil"
	keepertest "frontier/testutil/keeper"
	"frontier/x/migration"
	"frontier/x/migration/keeper"
	"frontier/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setup(t testing.TB) (types.MsgServer, context.Context, keeper.Keeper) {
	commontest.InitSDKConfig()

	k, ctx := keepertest.MigrationKeeper(t)
	
	// setup genesis params for this module
	genesis := *types.DefaultGenesis()
	genesis.MigratorList = []types.Migrator {
		{
			Migrator: commontest.Migrator_1,
		},
		{
			Migrator: commontest.Migrator_2,
		},
	}

	migration.InitGenesis(ctx, *k, genesis)
	
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)

	return server, context, *k
}

func TestShouldFailIfInvalidMigrator(t *testing.T) {
	// create a couple of migrators
	server, ctx, _ := setup(t)

	_, err := server.Migrate(ctx, &types.MsgMigrate{
		Creator: commontest.Alice,
		TxHash:  "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: commontest.Alice,
		Amount: "1000000000000000000000000", // 1 Milion
		Token: 0,
	})

	require.ErrorIs(t, err, types.ErrUnknownMigrator)
}
