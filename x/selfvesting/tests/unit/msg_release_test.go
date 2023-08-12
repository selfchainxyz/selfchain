package test

import (
	"context"
	"testing"
	"time"

	keepertest "selfchain/testutil/keeper"
	migrationTypes "selfchain/x/migration/types"
	"selfchain/x/selfvesting"
	"selfchain/x/selfvesting/keeper"
	test "selfchain/x/selfvesting/tests"
	mocktest "selfchain/x/selfvesting/tests/mock"
	"selfchain/x/selfvesting/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
)

const (
	NANO_SECONDS_IN_SECONDS = 1000000000
)

func setup_release(t testing.TB) (types.MsgServer, context.Context, keeper.Keeper, *gomock.Controller, *mocktest.MockBankKeeper) {
	ctrl := gomock.NewController(t)
	bankMock := mocktest.NewMockBankKeeper(ctrl)
	k, ctx := keepertest.SelfvestingKeeperWithMocks(t, bankMock)

	// setup genesis params for this module
	genesis := *types.DefaultGenesis()
	genesis.VestingPositionsList = []types.VestingPositions{}

	selfvesting.InitGenesis(ctx, *k, genesis)

	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)

	return server, context, *k, ctrl, bankMock
}

func setup_positions(t testing.TB, ctx context.Context, keeper keeper.Keeper) {
	ctx_1 := sdk.UnwrapSDKContext(ctx)
	// Start from 0, 0 time for simple math calculation
	sdkContext := ctx_1.WithBlockTime(time.Unix(0, 0))

	// Add two position for Alice and one for Bob
	keeper.AddBeneficiary(sdkContext, types.AddBeneficiaryRequest{
		Beneficiary: test.Alice,
		Cliff:       migrationTypes.VESTING_CLIFF,
		Duration:    migrationTypes.VESTING_DURATION,
		Amount:      "100000000000",
	})

	keeper.AddBeneficiary(sdkContext, types.AddBeneficiaryRequest{
		Beneficiary: test.Alice,
		Cliff:       migrationTypes.VESTING_CLIFF,
		Duration:    migrationTypes.VESTING_DURATION,
		Amount:      "200000000000",
	})

	keeper.AddBeneficiary(sdkContext, types.AddBeneficiaryRequest{
		Beneficiary: test.Bob,
		Cliff:       migrationTypes.VESTING_CLIFF,
		Duration:    migrationTypes.VESTING_DURATION,
		Amount:      "500000000000",
	})
}


func TestShouldReleaseLinearly(t *testing.T) {
	// test.InitSDKConfig()

	server, ctx, keeper, ctrl, bankMock := setup_release(t)
	setup_positions(t, ctx, keeper)
	defer ctrl.Finish()

	sdkContext := sdk.UnwrapSDKContext(ctx)
	vestingPositions, _ := keeper.GetVestingPositions(sdkContext, test.Alice)
	vestingInfo := vestingPositions.VestingInfos[0]

	// move 15 days (i.e. half of the vesting duration) into the future
	moveTo := time.Unix(0, int64(vestingInfo.StartTime) + (NANO_SECONDS_IN_SECONDS * migrationTypes.SECONDS_IN_DAY * 15))
	ctx_1 := sdkContext.WithBlockTime(moveTo)

	bankMock.ExpectReceiveCoins(ctx_1, test.Alice, 50000000000)

	server.Release(ctx_1, &types.MsgRelease {
		Creator: test.Alice,
    PosIndex:    0,
	})
}
