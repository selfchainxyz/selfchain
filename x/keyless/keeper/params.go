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

// GetMaxWalletsPerDID returns the maximum number of wallets allowed per DID
func (k Keeper) GetMaxWalletsPerDID(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).MaxWalletsPerDid
}

// GetMaxSharesPerWallet returns the maximum number of shares allowed per wallet
func (k Keeper) GetMaxSharesPerWallet(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).MaxSharesPerWallet
}

// GetMinRecoveryThreshold returns the minimum number of shares required for recovery
func (k Keeper) GetMinRecoveryThreshold(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).MinRecoveryThreshold
}

// GetMaxRecoveryThreshold returns the maximum number of shares allowed for recovery
func (k Keeper) GetMaxRecoveryThreshold(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).MaxRecoveryThreshold
}

// GetRecoveryWindowSeconds returns the time window in seconds for recovery
func (k Keeper) GetRecoveryWindowSeconds(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).RecoveryWindowSeconds
}

// GetMaxSigningAttempts returns the maximum number of signing attempts allowed
func (k Keeper) GetMaxSigningAttempts(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).MaxSigningAttempts
}
