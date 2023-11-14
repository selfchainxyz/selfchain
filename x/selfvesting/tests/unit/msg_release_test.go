package test

import (
	"context"
	"testing"
	"time"

	keepertest "selfchain/testutil/keeper"
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

const SECONDS_IN_DAY = 60 * 60 * 24

func setup_positions(t testing.TB, ctx context.Context, keeper keeper.Keeper) {
	ctx_1 := sdk.UnwrapSDKContext(ctx)
	// Start from 0, 0 time for simple math calculation
	sdkContext := ctx_1.WithBlockTime(time.Unix(0, 0))

	// Add two position for Alice and one for Bob
	keeper.AddBeneficiary(sdkContext, types.AddBeneficiaryRequest{
		Beneficiary: test.Alice,
		Cliff:       604800,
		Duration:    2592000,
		Amount:      "100000000000",
	})

	keeper.AddBeneficiary(sdkContext, types.AddBeneficiaryRequest{
		Beneficiary: test.Alice,
		Cliff:       604800,
		Duration:    2592000,
		Amount:      "200000000000",
	})

	keeper.AddBeneficiary(sdkContext, types.AddBeneficiaryRequest{
		Beneficiary: test.Bob,
		Cliff:       604800,
		Duration:    2592000,
		Amount:      "500000000000",
	})
}

func TestShouldReleaseLinearly(t *testing.T) {
	server, ctx, keeper, ctrl, bankMock := setup_release(t)
	setup_positions(t, ctx, keeper)
	defer ctrl.Finish()

	sdkContext := sdk.UnwrapSDKContext(ctx)
	vestingPositions, _ := keeper.GetVestingPositions(sdkContext, test.Alice)
	vestingInfo := vestingPositions.VestingInfos[0]

	// move 15 days (i.e. half of the vesting duration) into the future
	moveTo := time.Unix(0, int64(vestingInfo.StartTime)+(NANO_SECONDS_IN_SECONDS*SECONDS_IN_DAY*15))
	ctx_1 := sdkContext.WithBlockTime(moveTo)

	bankMock.ExpectReceiveCoins(ctx_1, test.Alice, 50000000000)

	server.Release(ctx_1, &types.MsgRelease{
		Creator:  test.Alice,
		PosIndex: 0,
	})

	vestingPositions_1, _ := keeper.GetVestingPositions(ctx_1, test.Alice)
	vestingInfo_1 := vestingPositions_1.VestingInfos[0]

	require.Equal(t, vestingPositions_1.Beneficiary, test.Alice)
	require.Equal(t, len(vestingPositions_1.VestingInfos), 2)
	require.Equal(t, vestingInfo_1.StartTime, uint64(0))
	require.Equal(t, vestingInfo_1.Cliff, uint64(0+604800))
	require.Equal(t, vestingInfo_1.Duration, uint64(2592000))
	require.Equal(t, vestingInfo_1.Amount, "100000000000")
	require.Equal(t, vestingInfo_1.TotalClaimed, "50000000000")
	require.Equal(t, vestingInfo_1.PeriodClaimed, uint64(moveTo.Unix()))

	// move 31 days (i.e. past the 30 days vesting duration) into the future
	moveTo_2 := time.Unix(0, int64(vestingInfo.StartTime)+(NANO_SECONDS_IN_SECONDS*2592000))
	ctx_2 := sdkContext.WithBlockTime(moveTo_2)

	bankMock.ExpectReceiveCoins(ctx_2, test.Alice, 50000000000)
	server.Release(ctx_2, &types.MsgRelease{
		Creator:  test.Alice,
		PosIndex: 0,
	})

	vestingPositions_2, _ := keeper.GetVestingPositions(ctx_2, test.Alice)
	vestingInfo_2 := vestingPositions_2.VestingInfos[0]

	require.Equal(t, vestingPositions_2.Beneficiary, test.Alice)
	require.Equal(t, len(vestingPositions_2.VestingInfos), 2)
	require.Equal(t, vestingInfo_2.StartTime, uint64(0))
	require.Equal(t, vestingInfo_2.Cliff, uint64(0+604800))
	require.Equal(t, vestingInfo_2.Duration, uint64(2592000))
	require.Equal(t, vestingInfo_2.Amount, "100000000000")
	require.Equal(t, vestingInfo_2.TotalClaimed, "100000000000")
	require.Equal(t, vestingInfo_2.PeriodClaimed, uint64(2592000))

	// Additional calls will cause a failure since the full amount was released
	_, releaseError := server.Release(ctx_2, &types.MsgRelease{
		Creator:  test.Alice,
		PosIndex: 0,
	})

	require.ErrorIs(t, releaseError, types.ErrPositionFullyClaimed)

	// Alice should be able to release the second position in parallel. Here we move it to 1/4 of the total duration i.e 7.5 days
	moveTo_3 := time.Unix(0, int64(vestingInfo.StartTime)+(NANO_SECONDS_IN_SECONDS*SECONDS_IN_DAY*7.5))
	ctx_3 := sdkContext.WithBlockTime(moveTo_3)
	bankMock.ExpectReceiveCoins(ctx_3, test.Alice, 50000000000)
	server.Release(ctx_3, &types.MsgRelease{
		Creator:  test.Alice,
		PosIndex: 1,
	})

	vestingPositions_3, _ := keeper.GetVestingPositions(ctx_3, test.Alice)
	vestingInfo_3 := vestingPositions_3.VestingInfos[1]

	require.Equal(t, vestingPositions_3.Beneficiary, test.Alice)
	require.Equal(t, len(vestingPositions_3.VestingInfos), 2)
	require.Equal(t, vestingInfo_3.StartTime, uint64(0))
	require.Equal(t, vestingInfo_3.Cliff, uint64(0+604800))
	require.Equal(t, vestingInfo_3.Duration, uint64(2592000))
	require.Equal(t, vestingInfo_3.Amount, "200000000000")
	require.Equal(t, vestingInfo_3.TotalClaimed, "50000000000")
	require.Equal(t, vestingInfo_3.PeriodClaimed, uint64(moveTo_3.Unix()))

	// Alice can release the entire amount after the end of the vesting period
	ctx_4 := sdkContext.WithBlockTime(moveTo_2)
	bankMock.ExpectReceiveCoins(ctx_4, test.Alice, 150000000000)
	server.Release(ctx_4, &types.MsgRelease{
		Creator:  test.Alice,
		PosIndex: 1,
	})

	vestingPositions_4, _ := keeper.GetVestingPositions(ctx_4, test.Alice)
	vestingInfo_4 := vestingPositions_4.VestingInfos[1]

	require.Equal(t, vestingPositions_4.Beneficiary, test.Alice)
	require.Equal(t, len(vestingPositions_4.VestingInfos), 2)
	require.Equal(t, vestingInfo_4.StartTime, uint64(0))
	require.Equal(t, vestingInfo_4.Cliff, uint64(0+604800))
	require.Equal(t, vestingInfo_4.Duration, uint64(2592000))
	require.Equal(t, vestingInfo_4.Amount, "200000000000")
	require.Equal(t, vestingInfo_4.TotalClaimed, "200000000000")
	require.Equal(t, vestingInfo_4.PeriodClaimed, uint64(2592000))

	// Bob should also be able to release in parallel with Alice. Here he starts at 1/4 of the total vesting duration
	ctx_5 := sdkContext.WithBlockTime(moveTo_3)
	bankMock.ExpectReceiveCoins(ctx_5, test.Bob, 125000000000)
	server.Release(ctx_5, &types.MsgRelease{
		Creator:  test.Bob,
		PosIndex: 0,
	})

	vestingPositions_5, _ := keeper.GetVestingPositions(ctx_5, test.Bob)
	vestingInfo_5 := vestingPositions_5.VestingInfos[0]

	require.Equal(t, vestingPositions_5.Beneficiary, test.Bob)
	require.Equal(t, len(vestingPositions_5.VestingInfos), 1)
	require.Equal(t, vestingInfo_5.StartTime, uint64(0))
	require.Equal(t, vestingInfo_5.Cliff, uint64(0+604800))
	require.Equal(t, vestingInfo_5.Duration, uint64(2592000))
	require.Equal(t, vestingInfo_5.Amount, "500000000000")
	require.Equal(t, vestingInfo_5.TotalClaimed, "125000000000")
	require.Equal(t, vestingInfo_5.PeriodClaimed, uint64(moveTo_3.Unix()))

	// Finally Bob can claim the full amount after the end of the vesting period
	ctx_6 := sdkContext.WithBlockTime(moveTo_2)
	bankMock.ExpectReceiveCoins(ctx_6, test.Bob, 375000000000)
	server.Release(ctx_6, &types.MsgRelease{
		Creator:  test.Bob,
		PosIndex: 0,
	})

	vestingPositions_6, _ := keeper.GetVestingPositions(ctx_6, test.Bob)
	vestingInfo_6 := vestingPositions_6.VestingInfos[0]

	require.Equal(t, vestingPositions_6.Beneficiary, test.Bob)
	require.Equal(t, len(vestingPositions_6.VestingInfos), 1)
	require.Equal(t, vestingInfo_6.StartTime, uint64(0))
	require.Equal(t, vestingInfo_6.Cliff, uint64(0+604800))
	require.Equal(t, vestingInfo_6.Duration, uint64(2592000))
	require.Equal(t, vestingInfo_6.Amount, "500000000000")
	require.Equal(t, vestingInfo_6.TotalClaimed, "500000000000")
	require.Equal(t, vestingInfo_6.PeriodClaimed, uint64(2592000))
}

func TestShouldFailIfNoVestingPositionExistForAccount(t *testing.T) {
	server, ctx, keeper, ctrl, _ := setup_release(t)
	setup_positions(t, ctx, keeper)
	defer ctrl.Finish()

	_, releaseError := server.Release(ctx, &types.MsgRelease{
		Creator:  test.Carol,
		PosIndex: 0,
	})

	require.ErrorIs(t, releaseError, types.ErrNoVestingPositions)
}

func TestShouldFailIfPosIndexOutofBounds(t *testing.T) {
	server, ctx, keeper, ctrl, _ := setup_release(t)
	setup_positions(t, ctx, keeper)
	defer ctrl.Finish()

	_, releaseError := server.Release(ctx, &types.MsgRelease{
		Creator:  test.Alice,
		PosIndex: 2,
	})

	require.ErrorIs(t, releaseError, types.ErrPositionIndexOutOfBounds)
}

func TestShouldFailIfCliffNotReached(t *testing.T) {
	server, ctx, keeper, ctrl, _ := setup_release(t)
	setup_positions(t, ctx, keeper)

	ctx_1 := sdk.UnwrapSDKContext(ctx)
	sdkContext := ctx_1.WithBlockTime(time.Unix(604800-1, 0))

	defer ctrl.Finish()

	_, releaseError := server.Release(sdkContext, &types.MsgRelease{
		Creator:  test.Alice,
		PosIndex: 1,
	})

	require.ErrorIs(t, releaseError, types.ErrCliffViolation)
}
