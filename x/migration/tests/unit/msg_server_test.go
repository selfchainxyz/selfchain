package test

import (
	"context"
	"testing"

	keepertest "frontier/testutil/keeper"
	"frontier/x/migration"
	"frontier/x/migration/keeper"
	test "frontier/x/migration/tests"
	mocktest "frontier/x/migration/tests/mock"
	"frontier/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func setup(t testing.TB) (types.MsgServer, context.Context, keeper.Keeper, *gomock.Controller, *mocktest.MockBankKeeper) {	
	ctrl := gomock.NewController(t)
	bankMock := mocktest.NewMockBankKeeper(ctrl)
	k, ctx := keepertest.MigrationKeeperWithMocks(t, bankMock)
	
	// setup genesis params for this module
	genesis := *types.DefaultGenesis()
	genesis.MigratorList = []types.Migrator {
		{
			Migrator: test.Migrator_1,
		},
		{
			Migrator: test.Migrator_2,
		},
	}

	migration.InitGenesis(ctx, *k, genesis)
	
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)

	return server, context, *k, ctrl, bankMock
}

func TestShouldFailIfInvalidMigrator(t *testing.T) {
	// create a couple of migrators
	server, ctx, _, ctrl, _ := setup(t)
	defer ctrl.Finish()

	_, err := server.Migrate(ctx, &types.MsgMigrate{
		Creator: test.Alice,
		TxHash:  "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount: "1000000000000000000000000", // 1 Milion
		Token: 0,
	})

	require.ErrorIs(t, err, types.ErrUnknownMigrator)
}

func TestShouldMintAmount(t *testing.T) {
	test.InitSDKConfig()

	// create a couple of migrators
	server, ctx, _, ctrl, mock := setup(t)
	defer ctrl.Finish()
	
	_, err := server.Migrate(ctx, &types.MsgMigrate{
		Creator: test.Migrator_1,
		TxHash:  "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount: "1000000000000000000000000", // 1 Milion
		Token: 0,
	})

	_ = err
	
	mock.ExpectMint(ctx, test.Alice, 1000000000000)
}
