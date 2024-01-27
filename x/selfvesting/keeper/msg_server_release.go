package keeper

import (
	"context"

	"selfchain/x/selfvesting/types"
	"selfchain/x/selfvesting/utils"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func getTokenReleaseInfo(
	k msgServer,
	ctx sdk.Context,
	beneficiary string,
	posIndex uint64,
) (*types.VestingInfo, uint64, sdkmath.Uint, error) {
	vestingPositions, positionsExist := k.GetVestingPositions(ctx, beneficiary)

	if !positionsExist {
		return &types.VestingInfo{}, 0, sdkmath.Uint{}, types.ErrNoVestingPositions
	}

	// Check that position at the given index exist
	if int(posIndex) >= len(vestingPositions.VestingInfos) {
		return &types.VestingInfo{}, 0, sdkmath.Uint{}, types.ErrPositionIndexOutOfBounds
	}

	vestingInfo := vestingPositions.VestingInfos[posIndex]

	// convert string values to uint 256
	amount := sdkmath.NewUintFromString(vestingInfo.Amount)
	totalClaimed := sdkmath.NewUintFromString(vestingInfo.TotalClaimed)

	// check if fully claimed and thus no more tokens exist to be released
	if totalClaimed.GTE(amount) {
		return &types.VestingInfo{}, 0, sdkmath.Uint{}, types.ErrPositionFullyClaimed
	}

	now := utils.BlockTime(ctx)
	// For vesting created with a future start date, that hasn't been reached, return 0, 0
	if now < vestingInfo.Cliff {
		return vestingInfo, 0, sdkmath.Uint{}, types.ErrCliffViolation
	}

	elapsedPeriod := now - vestingInfo.StartTime
	periodToVest := elapsedPeriod - vestingInfo.PeriodClaimed

	if elapsedPeriod >= vestingInfo.Duration {
		amountToVest := amount.Sub(totalClaimed)
		return vestingInfo, periodToVest, amountToVest, nil
	} else {
		amountToVest := amount.MulUint64(periodToVest).QuoUint64(vestingInfo.Duration)
		return vestingInfo, periodToVest, amountToVest, nil
	}
}

func (k msgServer) Release(goCtx context.Context, msg *types.MsgRelease) (*types.MsgReleaseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	vestingInfo, periodToVest, amountToVest, calcError := getTokenReleaseInfo(
		k,
		ctx,
		msg.Creator,
		msg.PosIndex,
	)

	if calcError != nil {
		return nil, calcError
	}

	if amountToVest.GT(sdkmath.NewUint(0)) {
		totalClaimed := sdkmath.NewUintFromString(vestingInfo.TotalClaimed)
		vestingInfo.PeriodClaimed += periodToVest
		vestingInfo.TotalClaimed = totalClaimed.Add(amountToVest).String()

		// store state changes
		vestingPositions, _ := k.GetVestingPositions(ctx, msg.Creator)
		vestingPositions.VestingInfos[msg.PosIndex] = vestingInfo
		k.SetVestingPositions(ctx, types.VestingPositions{
			Beneficiary:  msg.Creator,
			VestingInfos: vestingPositions.VestingInfos,
		})

		// transfer amountToVest to the beneficiary
		beneficiary, _ := sdk.AccAddressFromBech32(msg.Creator)
		vestedCoins := sdk.NewCoins(sdk.NewCoin(
			types.DENOM,
			sdkmath.NewIntFromBigInt(amountToVest.BigInt()),
		))

		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, beneficiary, vestedCoins)

		return &types.MsgReleaseResponse{
			PeriodToVest: periodToVest,
			AmountToVest: amountToVest.String(),
		}, nil
	}

	return &types.MsgReleaseResponse{
		PeriodToVest: 0,
		AmountToVest: "0",
	}, nil
}
