package keeper

import (
	"selfchain/x/migration/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetMigrator set a specific migrator in the store from its index
func (k Keeper) SetMigrator(ctx sdk.Context, migrator types.Migrator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MigratorKeyPrefix))
	b := k.cdc.MustMarshal(&migrator)
	store.Set(types.MigratorKey(
		migrator.Migrator,
	), b)
}

// GetMigrator returns a migrator from its index
func (k Keeper) GetMigrator(
	ctx sdk.Context,
	migrator string,

) (val types.Migrator, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MigratorKeyPrefix))

	b := store.Get(types.MigratorKey(
		migrator,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveMigrator removes a migrator from the store
func (k Keeper) RemoveMigrator(
	ctx sdk.Context,
	migrator string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MigratorKeyPrefix))
	store.Delete(types.MigratorKey(
		migrator,
	))
}

// GetAllMigrator returns all migrator
func (k Keeper) GetAllMigrator(ctx sdk.Context) (list []types.Migrator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MigratorKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Migrator
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
