package keeper

import (
	"selfchain/x/selfvesting/types"
	"selfchain/x/selfvesting/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) AddBeneficiary(ctx sdk.Context, req types.AddBeneficiaryRequest) (*types.VestingInfo, error) {
	// check the benficiary address is a valid bech32 address
	_, err := sdk.AccAddressFromBech32(req.Beneficiary)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid beneficiary address (%s)", err)
	}

	// create and add a new vesting position under the beneficiary key
	_, positionsExist := k.GetVestingPositions(ctx, req.Beneficiary)

	// if this is the first vesting position then create a default entry
	if !positionsExist {
		k.SetVestingPositions(ctx, types.VestingPositions{
			Beneficiary: req.Beneficiary,
		})
	}

	vestingPositions, _ := k.GetVestingPositions(ctx, req.Beneficiary)

	// startTime := uint64(ctx.BlockHeader().Time.Unix())
	startTime := utils.BlockTime(ctx)
	newPosition := &types.VestingInfo{
		StartTime:     startTime,
		Duration:      req.Duration,
		Cliff:         startTime + req.Cliff,
		Amount:        req.Amount,
		TotalClaimed:  "0",
		PeriodClaimed: 0,
	}

	k.SetVestingPositions(ctx, types.VestingPositions{
		Beneficiary:  req.Beneficiary,
		VestingInfos: append(vestingPositions.VestingInfos, newPosition),
	})

	return newPosition, nil
}
