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
	"github.com/stretchr/testify/require"
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

	vestingPositions_1, _ := keeper.GetVestingPositions(ctx_1, test.Alice)
	vestingInfo_1 := vestingPositions_1.VestingInfos[0]

	require.Equal(t, vestingPositions_1.Beneficiary, test.Alice)
	require.Equal(t,len(vestingPositions_1.VestingInfos), 2)
	require.Equal(t, vestingInfo_1.StartTime, uint64(0))
	require.Equal(t, vestingInfo_1.Cliff, uint64(0 + migrationTypes.VESTING_CLIFF))
	require.Equal(t, vestingInfo_1.Duration, uint64(migrationTypes.VESTING_DURATION))
	require.Equal(t, vestingInfo_1.Amount, "100000000000")
	require.Equal(t, vestingInfo_1.TotalClaimed, "50000000000")
	require.Equal(t, vestingInfo_1.PeriodClaimed, uint64(moveTo.Unix()))

	// move 31 days (i.e. past the 30 days vesting duration) into the future
	moveTo_2 := time.Unix(0, int64(vestingInfo.StartTime) + (NANO_SECONDS_IN_SECONDS * migrationTypes.VESTING_DURATION))
	ctx_2 := sdkContext.WithBlockTime(moveTo_2)

	bankMock.ExpectReceiveCoins(ctx_2, test.Alice, 50000000000)
	server.Release(ctx_2, &types.MsgRelease {
		Creator: test.Alice,
    PosIndex:    0,
	})

	vestingPositions_2, _ := keeper.GetVestingPositions(ctx_2, test.Alice)
	vestingInfo_2 := vestingPositions_2.VestingInfos[0]

	require.Equal(t, vestingPositions_2.Beneficiary, test.Alice)
	require.Equal(t,len(vestingPositions_2.VestingInfos), 2)
	require.Equal(t, vestingInfo_2.StartTime, uint64(0))
	require.Equal(t, vestingInfo_2.Cliff, uint64(0 + migrationTypes.VESTING_CLIFF))
	require.Equal(t, vestingInfo_2.Duration, uint64(migrationTypes.VESTING_DURATION))
	require.Equal(t, vestingInfo_2.Amount, "100000000000")
	require.Equal(t, vestingInfo_2.TotalClaimed, "100000000000")
	require.Equal(t, vestingInfo_2.PeriodClaimed, uint64(migrationTypes.VESTING_DURATION))

	// Additional calls will cause a failure since the full amount was released
	_, releaseError := server.Release(ctx_2, &types.MsgRelease {
		Creator: test.Alice,
    PosIndex:    0,
	})

	require.ErrorIs(t, releaseError, types.ErrPositionFullyClaimed)

	// At this point Alice can release the seconds position as well (in full since 30 days are passed)
	ctx_3 := sdkContext.WithBlockTime(moveTo_2)
	bankMock.ExpectReceiveCoins(ctx_3, test.Alice, 200000000000)
	server.Release(ctx_3, &types.MsgRelease {
		Creator: test.Alice,
    PosIndex:    1,
	})

	vestingPositions_3, _ := keeper.GetVestingPositions(ctx_3, test.Alice)
	vestingInfo_3 := vestingPositions_3.VestingInfos[1]

	require.Equal(t, vestingPositions_3.Beneficiary, test.Alice)
	require.Equal(t,len(vestingPositions_3.VestingInfos), 2)
	require.Equal(t, vestingInfo_3.StartTime, uint64(0))
	require.Equal(t, vestingInfo_3.Cliff, uint64(0 + migrationTypes.VESTING_CLIFF))
	require.Equal(t, vestingInfo_3.Duration, uint64(migrationTypes.VESTING_DURATION))
	require.Equal(t, vestingInfo_3.Amount, "200000000000")
	require.Equal(t, vestingInfo_3.TotalClaimed, "200000000000")
	require.Equal(t, vestingInfo_3.PeriodClaimed, uint64(migrationTypes.VESTING_DURATION))
}
