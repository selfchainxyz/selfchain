package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

const (
	// SocialIdentityPrefix is the prefix for storing social identities
	SocialIdentityPrefix = "social_identity/"
	// SocialIdentityByDIDPrefix is the prefix for storing social identities by DID
	SocialIdentityByDIDPrefix = "social_identity_by_did/"
	// SocialIdentityBySocialIDPrefix is the prefix for storing social identities by social ID
	SocialIdentityBySocialIDPrefix = "social_identity_by_social_id/"
)

// GetSocialIdentityByDIDAndProvider returns a social identity by DID and provider
func (k Keeper) GetSocialIdentityByDIDAndProvider(ctx sdk.Context, did string, provider string) (types.SocialIdentity, bool) {
	return k.GetSocialIdentityByDID(ctx, did, provider)
}

// SetSocialIdentity stores a social identity
func (k Keeper) SetSocialIdentity(ctx sdk.Context, identity types.SocialIdentity) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(SocialIdentityPrefix))
	key := []byte(fmt.Sprintf("%s/%s", identity.Provider, identity.Id))
	value := k.cdc.MustMarshal(&identity)
	store.Set(key, value)

	// Store by DID and provider
	byDIDStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(SocialIdentityByDIDPrefix))
	byDIDKey := []byte(fmt.Sprintf("%s/%s", identity.Did, identity.Provider))
	byDIDStore.Set(byDIDKey, value)

	// Store by provider and social ID
	bySocialIDStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(SocialIdentityBySocialIDPrefix))
	bySocialIDKey := []byte(fmt.Sprintf("%s/%s", identity.Provider, identity.Id))
	bySocialIDStore.Set(bySocialIDKey, value)
}

// DeleteSocialIdentity deletes a social identity
func (k Keeper) DeleteSocialIdentity(ctx sdk.Context, identity types.SocialIdentity) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(SocialIdentityPrefix))
	key := []byte(fmt.Sprintf("%s/%s", identity.Provider, identity.Id))
	store.Delete(key)

	// Delete by DID and provider
	byDIDStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(SocialIdentityByDIDPrefix))
	byDIDKey := []byte(fmt.Sprintf("%s/%s", identity.Did, identity.Provider))
	byDIDStore.Delete(byDIDKey)

	// Delete by provider and social ID
	bySocialIDStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(SocialIdentityBySocialIDPrefix))
	bySocialIDKey := []byte(fmt.Sprintf("%s/%s", identity.Provider, identity.Id))
	bySocialIDStore.Delete(bySocialIDKey)
}
