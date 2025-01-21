package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

// StoreSocialIdentity stores a social identity
func (k Keeper) StoreSocialIdentity(ctx sdk.Context, identity types.SocialIdentity) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(SocialIdentityPrefix))
	bz := k.cdc.MustMarshal(&identity)
	store.Set([]byte(identity.Id), bz)
	return nil
}

// GetSocialIdentityById returns a social identity by ID
func (k Keeper) GetSocialIdentityById(ctx sdk.Context, id string) (*types.SocialIdentity, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(SocialIdentityPrefix))
	bz := store.Get([]byte(id))
	if bz == nil {
		return nil, sdkerrors.Wrap(types.ErrSocialIdentityNotFound, id)
	}

	var identity types.SocialIdentity
	k.cdc.MustUnmarshal(bz, &identity)
	return &identity, nil
}

// DeleteSocialIdentity deletes a social identity
func (k Keeper) DeleteSocialIdentity(ctx sdk.Context, id string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(SocialIdentityPrefix))
	if !store.Has([]byte(id)) {
		return sdkerrors.Wrap(types.ErrSocialIdentityNotFound, id)
	}
	store.Delete([]byte(id))
	return nil
}

// GetAllSocialIdentities returns all social identities
func (k Keeper) GetAllSocialIdentities(ctx sdk.Context) []types.SocialIdentity {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(SocialIdentityPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var identities []types.SocialIdentity
	for ; iterator.Valid(); iterator.Next() {
		var identity types.SocialIdentity
		k.cdc.MustUnmarshal(iterator.Value(), &identity)
		identities = append(identities, identity)
	}
	return identities
}

// ValidateSocialIdentity validates a social identity
func (k Keeper) ValidateSocialIdentity(ctx sdk.Context, identity *types.SocialIdentity) error {
	if identity.Id == "" {
		return sdkerrors.Wrap(types.ErrInvalidSocialIdentity, "missing ID")
	}
	if identity.Did == "" {
		return sdkerrors.Wrap(types.ErrInvalidSocialIdentity, "missing DID")
	}
	if identity.Provider == "" {
		return sdkerrors.Wrap(types.ErrInvalidSocialIdentity, "missing provider")
	}
	if identity.ProviderId == "" {
		return sdkerrors.Wrap(types.ErrInvalidSocialIdentity, "missing social ID")
	}

	// Check if provider is allowed
	if !k.IsOAuthProviderAllowed(ctx, identity.Provider) {
		return sdkerrors.Wrap(types.ErrInvalidSocialIdentity, "provider not allowed")
	}

	return nil
}

// GetSocialIdentities returns all social identities for a DID
func (k Keeper) GetSocialIdentities(ctx sdk.Context, did string) ([]*types.SocialIdentity, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.SocialIdentityPrefix+did))
	defer iterator.Close()

	var identities []*types.SocialIdentity
	for ; iterator.Valid(); iterator.Next() {
		var identity types.SocialIdentity
		k.cdc.MustUnmarshal(iterator.Value(), &identity)
		identities = append(identities, &identity)
	}

	return identities, nil
}

// GetLinkedDID returns the DID linked to a social identity
func (k Keeper) GetLinkedDID(ctx sdk.Context, provider string, socialId string) (string, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(SocialIdentityBySocialIDPrefix))
	key := []byte(provider + ":" + socialId)
	bz := store.Get(key)
	if bz == nil {
		return "", false
	}

	var identity types.SocialIdentity
	k.cdc.MustUnmarshal(bz, &identity)
	return identity.Did, true
}
