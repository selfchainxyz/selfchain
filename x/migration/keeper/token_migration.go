package keeper

import (
	"selfchain/x/migration/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetTokenMigration set a specific tokenMigration in the store from its index
func (k Keeper) SetTokenMigration(ctx sdk.Context, tokenMigration types.TokenMigration) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenMigrationKeyPrefix))
	b := k.cdc.MustMarshal(&tokenMigration)
	store.Set(types.TokenMigrationKey(
		tokenMigration.MsgHash,
	), b)
}

// GetTokenMigration returns a tokenMigration from its index
func (k Keeper) GetTokenMigration(
	ctx sdk.Context,
	msgHash string,

) (val types.TokenMigration, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenMigrationKeyPrefix))

	b := store.Get(types.TokenMigrationKey(
		msgHash,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllTokenMigration returns all tokenMigration
func (k Keeper) GetAllTokenMigration(ctx sdk.Context) (list []types.TokenMigration) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TokenMigrationKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TokenMigration
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
