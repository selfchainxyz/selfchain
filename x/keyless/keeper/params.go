package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.ParamsKey))
	if bz == nil {
		return params
	}
	
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set([]byte(types.ParamsKey), bz)
}

// GetMaxParties returns the maximum number of parties allowed
func (k Keeper) GetMaxParties(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).MaxParties
}

// GetMaxThreshold returns the maximum threshold allowed
func (k Keeper) GetMaxThreshold(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).MaxThreshold
}

// GetMaxSecurityLevel returns the maximum security level allowed
func (k Keeper) GetMaxSecurityLevel(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).MaxSecurityLevel
}

// GetMaxBatchSize returns the maximum batch size allowed
func (k Keeper) GetMaxBatchSize(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).MaxBatchSize
}

// GetMaxMetadataSize returns the maximum metadata size allowed
func (k Keeper) GetMaxMetadataSize(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).MaxMetadataSize
}
