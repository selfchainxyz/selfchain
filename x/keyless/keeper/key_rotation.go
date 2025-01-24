package keeper

import (
	"selfchain/x/keyless/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetKeyRotationStatus retrieves the key rotation status for a given wallet
func (k Keeper) GetKeyRotationStatus(ctx sdk.Context, walletId string) (status types.KeyRotationStatus, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyRotationStatusKey(walletId)
	
	bz := store.Get(key)
	if bz == nil {
		return types.KeyRotationStatus_KEY_ROTATION_STATUS_UNSPECIFIED, false
	}

	val, ok := types.KeyRotationStatus_value[string(bz)]
	if !ok {
		return types.KeyRotationStatus_KEY_ROTATION_STATUS_UNSPECIFIED, false
	}
	return types.KeyRotationStatus(val), true
}

// SetKeyRotationStatus stores the key rotation status for a given wallet
func (k Keeper) SetKeyRotationStatus(ctx sdk.Context, walletId string, status types.KeyRotationStatus) {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyRotationStatusKey(walletId)
	store.Set(key, []byte(status.String()))
}
