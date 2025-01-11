package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

const (
	VerificationStatusPrefix = "verification_status/"
)

// SetVerificationStatus stores a verification status
func (k Keeper) SetVerificationStatus(ctx sdk.Context, did string, status types.VerificationStatus) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(VerificationStatusPrefix))
	statusBytes := k.cdc.MustMarshal(&status)
	store.Set([]byte(did), statusBytes)
	return nil
}

// GetVerificationStatus returns a verification status by DID
func (k Keeper) GetVerificationStatus(ctx sdk.Context, did string) (types.VerificationStatus, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(VerificationStatusPrefix))
	statusBytes := store.Get([]byte(did))
	if statusBytes == nil {
		return types.VerificationStatus{}, false
	}

	var status types.VerificationStatus
	k.cdc.MustUnmarshal(statusBytes, &status)
	return status, true
}

// DeleteExpiredVerifications deletes all expired verification statuses
func (k Keeper) DeleteExpiredVerifications(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(VerificationStatusPrefix))
	now := ctx.BlockTime()

	var toDelete [][]byte
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var status types.VerificationStatus
		k.cdc.MustUnmarshal(iterator.Value(), &status)

		if status.ExpiresAt != nil && status.ExpiresAt.Before(now) {
			toDelete = append(toDelete, iterator.Key())
		}
	}

	for _, key := range toDelete {
		store.Delete(key)
	}
}
