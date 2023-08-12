package keeper

import (
	"context"
	"selfchain/x/selfvesting/types"
	"selfchain/x/selfvesting/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) AddBeneficiary(goCtx context.Context, req types.AddBeneficiaryRequest) (*types.VestingInfo, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check the benficiary address is a valid bech32 address
	_, err := sdk.AccAddressFromBech32(req.Beneficiary)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid beneficiary address (%s)", err)
	}

	// create and add a new vesting position under the beneficiary key
	vestingInfos := []*types.VestingInfo{}
	vestingPositions, positionsExist := k.GetVestingPositions(ctx, req.Beneficiary)

	if positionsExist {
		vestingInfos = vestingPositions.VestingInfos
	}

	// startTime := uint64(ctx.BlockHeader().Time.Unix())
	startTime := utils.BlockTime(ctx)
	newPosition := &types.VestingInfo{
		StartTime:     startTime,
		Duration:      req.Duration,
		Cliff:         startTime + req.Cliff,
		Amount:        req.Amount,
		TotalClaimed:  0,
		PeriodClaimed: 0,
	}

	vestingPositions.VestingInfos = append(vestingInfos, newPosition)

	return newPosition, nil
}
