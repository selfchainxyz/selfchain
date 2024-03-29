package test

import (
	"context"
	"testing"

	keepertest "selfchain/testutil/keeper"
	"selfchain/x/migration"
	"selfchain/x/migration/keeper"
	test "selfchain/x/migration/tests"
	mocktest "selfchain/x/migration/tests/mock"
	"selfchain/x/migration/types"
	selfvestingTypes "selfchain/x/selfvesting/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func setup(t testing.TB) (types.MsgServer, context.Context, keeper.Keeper, *gomock.Controller, *mocktest.MockSelfvestingKeeper, *mocktest.MockBankKeeper) {
	ctrl := gomock.NewController(t)
	bankMock := mocktest.NewMockBankKeeper(ctrl)
	selfVestingMock := mocktest.NewMockSelfvestingKeeper(ctrl)
	k, ctx := keepertest.MigrationKeeperWithMocks(t, selfVestingMock, bankMock)

	// setup genesis params for this module
	genesis := *types.DefaultGenesis()
	genesis.MigratorList = []types.Migrator{
		{
			Migrator: test.Migrator_1,
			Exists:   true,
		},
		{
			Migrator: test.Migrator_2,
			Exists:   true,
		},
	}

	genesis.Config = &types.Config{
		VestingDuration:    2592000, // 1 month in seconds
		VestingCliff:       604800, // // 1 week in seconds
		MinMigrationAmount: 1000000000000000000,
	}

	migration.InitGenesis(ctx, *k, genesis)

	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)

	return server, context, *k, ctrl, selfVestingMock, bankMock
}

func TestShouldFailIfInvalidMigrator(t *testing.T) {
	// create a couple of migrators
	server, ctx, _, ctrl, _, _ := setup(t)
	defer ctrl.Finish()

	_, err := server.Migrate(ctx, &types.MsgMigrate{
		Creator:     test.Alice,
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount:      "1000000000000000000000000", // 1 Milion
		Token:       0,
		LogIndex:    0,
	})

	require.ErrorIs(t, err, types.ErrUnknownMigrator)
}

func TestShouldMintAmountAndAddBeneficiary(t *testing.T) {
	// create a couple of migrators
	server, ctx, _, ctrl, selfVestingMock, bankMock := setup(t)
	defer ctrl.Finish()

	addBeneficiaryRequest := selfvestingTypes.AddBeneficiaryRequest{
		Beneficiary: test.Alice,
		Cliff:       604800,
		Duration:    2592000,
		Amount:      "999999000000",
	}
	
	bankMock.ExpectMintToModule(ctx, 1000000000000)
	// 1 SLF is instanly resealed
	bankMock.ExpectReceiveCoins(ctx, selfvestingTypes.ModuleName, test.Alice, 1000000)
	selfVestingMock.ExpectAddBeneficiary(ctx, addBeneficiaryRequest)

	_, err := server.Migrate(ctx, &types.MsgMigrate{
		Creator:     test.Migrator_1,
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount:      "1000000000000000000000000", // 1 Milion
		Token:       0,
		LogIndex:    0,
	})

	_ = err
}

func TestShouldMintAmountButNotAddBeneficiary(t *testing.T) {
	// create a couple of migrators
	server, ctx, _, ctrl, _, bankMock := setup(t)
	defer ctrl.Finish()

	bankMock.ExpectMintToModule(ctx, 1000000)
	bankMock.ExpectReceiveCoins(ctx, selfvestingTypes.ModuleName, test.Alice, 1000000)

	_, err := server.Migrate(ctx, &types.MsgMigrate{
		Creator:     test.Migrator_1,
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount:      "1000000000000000000", // 1 FRONT token
		Token:       0,
		LogIndex:    0,
	})

	_ = err
}
