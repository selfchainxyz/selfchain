package v002

import (
	types "frontier/x/migration/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func addTimeField(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.TokenMigrationKeyPrefix))

	defer iterator.Close()
  for ; iterator.Valid(); iterator.Next() {
		var acc types.TokenMigration

		if err := cdc.UnmarshalInterface(iterator.Value(), &acc); err != nil {
			return err
		}
	}
	return nil
}

// MigrateStore performs in-place store migrations from v0.01 to v0.02. The
// migration includes:
// - Add a new additional field to TokenMigration
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	return addTimeField(ctx, storeKey, cdc)
}
