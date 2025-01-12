package keeper

import (
	"fmt"

	"selfchain/x/identity/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// LinkSocialIdentity links a social identity to a DID
func (k Keeper) LinkSocialIdentity(ctx sdk.Context, did string, provider string, token string) error {
	// Verify DID exists
	if !k.HasDIDDocument(ctx, did) {
		return sdkerrors.Wrap(types.ErrDIDNotFound, "DID not found")
	}

	// Verify token with provider
	socialID, err := k.verifySocialToken(ctx, provider, token)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidToken, err.Error())
	}

	// Create social identity
	blockTime := ctx.BlockTime()
	identity := types.SocialIdentity{
		Did:        did,
		Provider:   provider,
		ProviderId: socialID,
		CreatedAt:  &blockTime,
		LastUsed:   &blockTime,
	}

	// Store social identity
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SocialIdentityPrefix))
	key := []byte(fmt.Sprintf("%s-%s", did, provider))
	bz, err := k.cdc.Marshal(&identity)
	if err != nil {
		return err
	}
	store.Set(key, bz)

	return nil
}

// UnlinkSocialIdentity unlinks a social identity from a DID
func (k Keeper) UnlinkSocialIdentity(ctx sdk.Context, did string, socialID string) error {
	// Find and delete the social identity
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SocialIdentityPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte(did))
	defer iterator.Close()

	found := false
	for ; iterator.Valid(); iterator.Next() {
		var identity types.SocialIdentity
		if err := k.cdc.Unmarshal(iterator.Value(), &identity); err != nil {
			return err
		}
		if identity.ProviderId == socialID {
			store.Delete(iterator.Key())
			found = true
			break
		}
	}

	if !found {
		return sdkerrors.Wrap(types.ErrSocialIdentityNotFound, "social identity not found")
	}

	return nil
}

// verifySocialToken verifies a social identity token with the provider
func (k Keeper) verifySocialToken(ctx sdk.Context, provider string, token string) (string, error) {
	// TODO: Implement actual token verification with providers
	// For now, just return a dummy social ID
	return fmt.Sprintf("%s-user-123", provider), nil
}
