package test

import (
	"context"
	"testing"

	keepertest "selfchain/testutil/keeper"
	"selfchain/x/selfvesting"
	"selfchain/x/selfvesting/keeper"
	test "selfchain/x/selfvesting/tests"
	mocktest "selfchain/x/selfvesting/tests/mock"
	"selfchain/x/selfvesting/types"
	"selfchain/x/selfvesting/utils"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func setup(t testing.TB) (types.MsgServer, context.Context, keeper.Keeper, *gomock.Controller, *mocktest.MockBankKeeper) {
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

func TestShouldFailIfInvalidBeneficiaryAddr(t *testing.T) {
	test.InitSDKConfig()

	_, ctx, keeper, ctrl, _ := setup(t)
	defer ctrl.Finish()

	_, err := keeper.AddBeneficiary(sdk.UnwrapSDKContext(ctx), types.AddBeneficiaryRequest{
		Beneficiary: "Invalid Address",
		Cliff:       0,
		Duration:    0,
		Amount:      "",
	})

	require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
}

func TestShouldCreateNewVestingPosition(t *testing.T) {
	_, ctx, keeper, ctrl, _ := setup(t)
	defer ctrl.Finish()

	sdkContext := sdk.UnwrapSDKContext(ctx)
	startTime := utils.BlockTime(sdkContext)

	addBeneficiaryRequest := types.AddBeneficiaryRequest{
		Beneficiary: test.Alice,
		Cliff:       604800,
		Duration:    2592000,
		Amount:      "100000000000",
	}

	_, err := keeper.AddBeneficiary(sdkContext, addBeneficiaryRequest)

	vestingPositions, _ := keeper.GetVestingPositions(sdkContext, test.Alice)
	vestingInfo := vestingPositions.VestingInfos[0]

	require.Equal(t, vestingPositions.Beneficiary, test.Alice)
	require.Equal(t, len(vestingPositions.VestingInfos), 1)
	require.Equal(t, vestingInfo.StartTime, startTime)
	require.Equal(t, vestingInfo.Cliff, startTime+604800)
	require.Equal(t, vestingInfo.Duration, uint64(2592000))
	require.Equal(t, vestingInfo.Amount, "100000000000")
	require.Equal(t, vestingInfo.TotalClaimed, "0")
	require.Equal(t, vestingInfo.PeriodClaimed, uint64(0))

	_ = err

	// Add one more position
	addBeneficiaryRequest2 := types.AddBeneficiaryRequest{
		Beneficiary: test.Alice,
		Cliff:       604800,
		Duration:    2592000,
		Amount:      "200000000000",
	}

	keeper.AddBeneficiary(sdkContext, addBeneficiaryRequest2)

	vestingPositions_1, _ := keeper.GetVestingPositions(sdkContext, test.Alice)
	vestingInfo_1 := vestingPositions_1.VestingInfos[0]

	// The first position remains intact
	require.Equal(t, vestingPositions_1.Beneficiary, test.Alice)
	require.Equal(t, len(vestingPositions_1.VestingInfos), 2)
	require.Equal(t, vestingInfo_1.StartTime, startTime)
	require.Equal(t, vestingInfo_1.Cliff, startTime+604800)
	require.Equal(t, vestingInfo_1.Duration, uint64(2592000))
	require.Equal(t, vestingInfo_1.Amount, "100000000000")
	require.Equal(t, vestingInfo_1.TotalClaimed, "0")
	require.Equal(t, vestingInfo_1.PeriodClaimed, uint64(0))

	// second position should be stored
	vestingPositions_2, _ := keeper.GetVestingPositions(sdkContext, test.Alice)
	vestingInfo_2 := vestingPositions_2.VestingInfos[1]

	require.Equal(t, vestingPositions_2.Beneficiary, test.Alice)
	require.Equal(t, len(vestingPositions_2.VestingInfos), 2)
	require.Equal(t, vestingInfo_2.StartTime, startTime)
	require.Equal(t, vestingInfo_2.Cliff, startTime+604800)
	require.Equal(t, vestingInfo_2.Duration, uint64(2592000))
	require.Equal(t, vestingInfo_2.Amount, "200000000000")
	require.Equal(t, vestingInfo_2.TotalClaimed, "0")
	require.Equal(t, vestingInfo_2.PeriodClaimed, uint64(0))
}
