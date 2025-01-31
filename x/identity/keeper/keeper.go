package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"selfchain/x/identity/types"
)

type (
	Keeper struct {
		cdc         codec.BinaryCodec
		storeKey    storetypes.StoreKey
		memKey      storetypes.StoreKey
		paramstore  paramtypes.Subspace
		keyless     types.KeylessKeeper
		rateLimiter *RateLimiter
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	keyless types.KeylessKeeper,
) *Keeper {
	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:         cdc,
		storeKey:    storeKey,
		memKey:      memKey,
		paramstore:  ps,
		keyless:     keyless,
		rateLimiter: NewRateLimiter(100, 100), // 100 requests per second with burst of 100
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetStore returns a store for a specific prefix
func (k Keeper) GetStore(ctx sdk.Context, storePrefix []byte) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), storePrefix)
}

// GetSocialIdentitiesByDID returns all social identities for a DID
func (k Keeper) GetSocialIdentitiesByDID(ctx sdk.Context, did string) []types.SocialIdentity {
	store := k.GetStore(ctx, []byte(types.SocialIdentityPrefix))
	iterator := store.Iterator([]byte(did), nil)
	defer iterator.Close()

	var identities []types.SocialIdentity
	for ; iterator.Valid(); iterator.Next() {
		var identity types.SocialIdentity
		k.cdc.MustUnmarshal(iterator.Value(), &identity)
		identities = append(identities, identity)
	}

	return identities
}

// GetSocialIdentity gets a social identity for a DID and provider
func (k Keeper) GetSocialIdentity(ctx sdk.Context, did string, provider string) (*types.SocialIdentity, bool) {
	store := k.GetStore(ctx, []byte(types.SocialIdentityPrefix))
	key := append([]byte(did), []byte(provider)...)

	if !store.Has(key) {
		return nil, false
	}

	var identity types.SocialIdentity
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &identity)

	return &identity, true
}

// HasSocialIdentity checks if a social identity exists for a DID and provider
func (k Keeper) HasSocialIdentity(ctx sdk.Context, did string, provider string) bool {
	store := k.GetStore(ctx, []byte(types.SocialIdentityPrefix))
	key := append([]byte(did), []byte(provider)...)
	return store.Has(key)
}

// GetOAuthProvider gets an OAuth provider configuration
func (k Keeper) GetOAuthProvider(ctx sdk.Context, provider string) (*types.OAuthProvider, error) {
	store := k.GetStore(ctx, []byte(types.OAuthProviderPrefix))
	key := []byte(provider)

	if !store.Has(key) {
		return nil, sdkerrors.Wrapf(types.ErrOAuthProviderNotFound, "provider %s not found", provider)
	}

	var oauthProvider types.OAuthProvider
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &oauthProvider)

	return &oauthProvider, nil
}

// IsAuthorized checks if an address is authorized for a DID
func (k Keeper) IsAuthorized(ctx sdk.Context, did string, address string) bool {
	doc, found := k.GetDIDDocument(ctx, did)
	if !found {
		return false
	}

	// Check if the address is in the controllers list
	for _, controller := range doc.Controller {
		if controller == address {
			return true
		}
	}

	return false
}

// GetParams returns the total set of identity parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the identity parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// ReconstructWalletFromDID reconstructs a wallet from a DID document
func (k Keeper) ReconstructWalletFromDID(ctx sdk.Context, doc types.DIDDocument) ([]byte, error) {
	if doc.Id == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidDIDDocument, "DID document ID is empty")
	}

	// Get the key share from the keyless module
	keyShare, found := k.keyless.GetKeyShare(ctx, doc.Id)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrKeyShareNotFound, "key share not found for DID")
	}

	// Get the wallet from the keyless module
	wallet, err := k.keyless.ReconstructWallet(ctx, doc)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to reconstruct wallet")
	}

	// Combine the key share with the wallet
	combinedWallet := append(wallet, keyShare...)
	return combinedWallet, nil
}

// SaveDIDDocument saves a DID document to the store
func (k Keeper) SaveDIDDocument(ctx sdk.Context, doc types.DIDDocument) error {
	if doc.Id == "" {
		return sdkerrors.Wrap(types.ErrInvalidDIDDocument, "DID document ID is empty")
	}

	store := k.GetStore(ctx, types.DIDDocumentKey)
	bz := k.cdc.MustMarshal(&doc)
	store.Set([]byte(doc.Id), bz)

	return nil
}
