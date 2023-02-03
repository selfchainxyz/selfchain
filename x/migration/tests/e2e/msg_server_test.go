package test

import (
	test "frontier/x/migration/tests"
	"frontier/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// func setup(t testing.TB) (types.MsgServer, context.Context, keeper.Keeper) {
// 	k, ctx := keepertest.MigrationKeeper(t)

// 	// setup genesis params for this module
// 	genesis := *types.DefaultGenesis()
// 	genesis.MigratorList = []types.Migrator {
// 		{
// 			Migrator: Migrator_1,
// 		},
// 		{
// 			Migrator: Migrator_2,
// 		},
// 	}

// 	migration.InitGenesis(ctx, *k, genesis)

// 	server := keeper.NewMsgServerImpl(*k)
// 	context := sdk.WrapSDKContext(ctx)

// 	return server, context, *k
// }

func (suite *IntegrationTestSuite) TestShouldFailIfInvalidMigrator() {
	suite.setupSuiteWithBalances()
	ctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.msgServer.Migrate(ctx, &types.MsgMigrate{
		Creator: test.Alice,
		TxHash:  "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount: "1000000000000000000000000", // 1 Milion
		Token: 0,
	})

	suite.Require().ErrorIs(err, types.ErrIntOverflowAcl)
}

