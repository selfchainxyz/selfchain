package keeper

import (
	"selfchain/x/migration/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetAcl set acl in the store
func (k Keeper) SetAcl(ctx sdk.Context, acl types.Acl) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AclKey))
	b := k.cdc.MustMarshal(&acl)
	store.Set([]byte{0}, b)
}

// GetAcl returns acl
func (k Keeper) GetAcl(ctx sdk.Context) (val types.Acl, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AclKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
