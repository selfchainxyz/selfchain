package keeper

import (
	"fmt"

	"selfchain/x/identity/types"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
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
	// set KeyTable if it has not already been set
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

// GetParams returns the total set of identity parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the identity parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// VerifyOAuthToken verifies an OAuth token with the provider
func (k Keeper) VerifyOAuthToken(ctx sdk.Context, provider string, token string) (string, error) {
	// TODO: Implement OAuth token verification
	return "", fmt.Errorf("not implemented")
}

// GetSocialIdentitiesByDID returns all social identities for a DID
func (k Keeper) GetSocialIdentitiesByDID(ctx sdk.Context, did string) []types.SocialIdentity {
	// TODO: Implement getting social identities by DID
	return []types.SocialIdentity{}
}

// HasSocialIdentity checks if a social identity exists
func (k Keeper) HasSocialIdentity(ctx sdk.Context, did string, provider string, userInfo string) bool {
	// TODO: Implement checking if social identity exists
	return false
}

// GetSocialIdentity returns a social identity by DID and provider
func (k Keeper) GetSocialIdentity(ctx sdk.Context, did string, provider string) (*types.SocialIdentity, bool) {
	// TODO: Implement getting social identity
	return nil, false
}

// GetSocialIdentityBySocialID returns a social identity by provider and social ID
func (k Keeper) GetSocialIdentityBySocialID(ctx sdk.Context, provider string, socialID string) (*types.SocialIdentity, bool) {
	store := k.GetStore(ctx, []byte(types.SocialIdentityByIDPrefix))
	prefixKey := append([]byte(provider), []byte(socialID)...)
	bz := store.Get(prefixKey)
	if bz == nil {
		return nil, false
	}

	var identity types.SocialIdentity
	k.cdc.MustUnmarshal(bz, &identity)
	return &identity, true
}

// IsAuthorized checks if an address is authorized for a DID
func (k Keeper) IsAuthorized(ctx sdk.Context, did string, address string) bool {
	doc, found := k.GetDIDDocument(ctx, did)
	if !found {
		return false
	}

	// Check if the address is the controller
	for _, controller := range doc.Controller {
		if controller == address {
			return true
		}
	}

	return false
}

// GetDIDFromMsg extracts DID from various message types
func (k Keeper) GetDIDFromMsg(msg sdk.Msg) string {
	switch v := msg.(type) {
	case *types.MsgCreateDID:
		return v.Id
	case *types.MsgUpdateDID:
		return v.Id
	case *types.MsgDeleteDID:
		return v.Id
	case *types.MsgLinkSocialIdentity:
		return v.Creator
	case *types.MsgUnlinkSocialIdentity:
		return v.Creator
	default:
		return ""
	}
}

// GetStore returns a store for a given prefix
func (k Keeper) GetStore(ctx sdk.Context, storePrefix []byte) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), storePrefix)
}
