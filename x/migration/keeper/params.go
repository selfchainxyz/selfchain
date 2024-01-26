package keeper

import (
	"selfchain/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HotcrossRatio(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.HotcrossRatio, &res)
	return
}

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.HotcrossRatio(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
