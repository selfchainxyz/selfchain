package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/selfvesting/types"
)

// SetVestingPositions set a specific vestingPositions in the store from its index
func (k Keeper) SetVestingPositions(ctx sdk.Context, vestingPositions types.VestingPositions) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VestingPositionsKeyPrefix))
	b := k.cdc.MustMarshal(&vestingPositions)
	store.Set(types.VestingPositionsKey(
		vestingPositions.Beneficiary,
	), b)
}

// GetVestingPositions returns a vestingPositions from its index
func (k Keeper) GetVestingPositions(
	ctx sdk.Context,
	beneficiary string,

) (val types.VestingPositions, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VestingPositionsKeyPrefix))

	b := store.Get(types.VestingPositionsKey(
		beneficiary,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveVestingPositions removes a vestingPositions from the store
func (k Keeper) RemoveVestingPositions(
	ctx sdk.Context,
	beneficiary string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VestingPositionsKeyPrefix))
	store.Delete(types.VestingPositionsKey(
		beneficiary,
	))
}

// GetAllVestingPositions returns all vestingPositions
func (k Keeper) GetAllVestingPositions(ctx sdk.Context) (list []types.VestingPositions) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VestingPositionsKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.VestingPositions
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
