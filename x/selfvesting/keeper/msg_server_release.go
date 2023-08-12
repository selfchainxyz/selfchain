package keeper

import (
	"context"

	"selfchain/x/selfvesting/types"
	"selfchain/x/selfvesting/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func getTokenReleaseInfo(
	k msgServer,
	ctx sdk.Context,
	beneficiary string,
	posIndex uint64,
) (uint64, uint64, error) {
	vestingPositions, positionsExist := k.GetVestingPositions(ctx, beneficiary)

	if !positionsExist {
		return 0, 0, types.ErrNoVestingPositions
	}

	// Check that position at the given index exist
	if int(posIndex) > len(vestingPositions.VestingInfos) {
    return 0, 0, types.ErrPositionIndexOutOfBounds
	}

	vestingInfo := vestingPositions.VestingInfos[posIndex]

	// check if fully claimed and thus no more tokens exist to be released
	if vestingInfo.TotalClaimed >= vestingInfo.Amount {
		return 0, 0, types.ErrPositionFullyClaimed
	}

	now := utils.BlockTime(ctx)
	// For vesting created with a future start date, that hasn't been reached, return 0, 0
	if now < vestingInfo.Cliff {
		return 0, 0, nil
	}

	elapsedPeriod := now - vestingInfo.StartTime
	periodToVest := elapsedPeriod - vestingInfo.PeriodClaimed

	if elapsedPeriod >= vestingInfo.Duration {
		amountToVest := vestingInfo.Amount - vestingInfo.TotalClaimed
		return periodToVest, amountToVest, nil
	} else {
		amountToVest := (periodToVest * vestingInfo.Amount) / vestingInfo.Duration
		return periodToVest, amountToVest, nil
	}
}

func (k msgServer) Release(goCtx context.Context, msg *types.MsgRelease) (*types.MsgReleaseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgReleaseResponse{}, nil
}
