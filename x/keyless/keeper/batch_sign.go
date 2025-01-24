package keeper

import (
	"selfchain/x/keyless/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetBatchSignStatus retrieves the batch sign status for a given wallet and batch ID
func (k Keeper) GetBatchSignStatus(ctx sdk.Context, walletId string, batchId string) (status types.BatchSignStatus, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.BatchSignStatusKey(walletId, batchId)
	
	bz := store.Get(key)
	if bz == nil {
		return types.BatchSignStatus_BATCH_SIGN_STATUS_UNSPECIFIED, false
	}

	val, ok := types.BatchSignStatus_value[string(bz)]
	if !ok {
		return types.BatchSignStatus_BATCH_SIGN_STATUS_UNSPECIFIED, false
	}
	return types.BatchSignStatus(val), true
}

// SetBatchSignStatus stores the batch sign status for a given wallet and batch ID
func (k Keeper) SetBatchSignStatus(ctx sdk.Context, walletId string, batchId string, status types.BatchSignStatus) {
	store := ctx.KVStore(k.storeKey)
	key := types.BatchSignStatusKey(walletId, batchId)
	store.Set(key, []byte(status.String()))
}
