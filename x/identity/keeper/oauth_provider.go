package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"selfchain/x/identity/types"
)

// StoreOAuthProvider stores an OAuth provider
func (k Keeper) StoreOAuthProvider(ctx sdk.Context, provider *types.OAuthProvider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OAuthProviderPrefix))
	bz := k.cdc.MustMarshal(provider)
	store.Set([]byte(provider.Id), bz)
	return nil
}

// GetAllOAuthProviders returns all OAuth providers
func (k Keeper) GetAllOAuthProviders(ctx sdk.Context) []*types.OAuthProvider {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OAuthProviderPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var providers []*types.OAuthProvider
	for ; iterator.Valid(); iterator.Next() {
		var provider types.OAuthProvider
		k.cdc.MustUnmarshal(iterator.Value(), &provider)
		providers = append(providers, &provider)
	}
	return providers
}
