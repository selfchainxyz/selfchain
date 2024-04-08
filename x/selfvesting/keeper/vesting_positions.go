package keeper

import (
	"selfchain/x/selfvesting/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
