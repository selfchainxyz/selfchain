package keeper

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

// StoreSocialIdentity stores a social identity
func (k Keeper) StoreSocialIdentity(ctx sdk.Context, identity types.SocialIdentity) error {
	// Validate required fields
	if identity.Did == "" {
		return fmt.Errorf("DID cannot be empty")
	}
	if identity.Provider == "" {
		return fmt.Errorf("provider cannot be empty")
	}
	if identity.Id == "" {
		return fmt.Errorf("social ID cannot be empty")
	}

	// Check if DID exists
	if !k.HasDIDDocument(ctx, identity.Did) {
		return fmt.Errorf("DID not found: %s", identity.Did)
	}

	// Store by DID and provider
	storeDID := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SocialIdentityKey+"did/"))
	keyDID := []byte(identity.Did + "/" + identity.Provider)
	bz := k.cdc.MustMarshal(&identity)
	storeDID.Set(keyDID, bz)

	// Store by social ID and provider
	storeSocial := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SocialIdentityKey+"social/"))
	keySocial := []byte(identity.Provider + "/" + identity.Id)
	storeSocial.Set(keySocial, bz)

	return nil
}

// GetSocialIdentityByDID returns a social identity by DID and provider
func (k Keeper) GetSocialIdentityByDID(ctx sdk.Context, did string, provider string) (types.SocialIdentity, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SocialIdentityKey+"did/"))
	key := []byte(did + "/" + provider)
	b := store.Get(key)
	if b == nil {
		return types.SocialIdentity{}, false
	}

	var identity types.SocialIdentity
	k.cdc.MustUnmarshal(b, &identity)
	return identity, true
}

// GetSocialIdentityBySocialID returns a social identity by social ID and provider
func (k Keeper) GetSocialIdentityBySocialID(ctx sdk.Context, provider string, socialId string) (types.SocialIdentity, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SocialIdentityKey+"social/"))
	key := []byte(provider + "/" + socialId)
	b := store.Get(key)
	if b == nil {
		return types.SocialIdentity{}, false
	}

	var identity types.SocialIdentity
	k.cdc.MustUnmarshal(b, &identity)
	return identity, true
}

// VerifyOAuthToken verifies an OAuth token and returns the social ID
func (k Keeper) VerifyOAuthToken(ctx sdk.Context, provider string, token string) (string, error) {
	// Check rate limit
	rateLimitKey := fmt.Sprintf("oauth:%s:%s", provider, token[:8]) // Use first 8 chars of token
	if err := k.CheckRateLimit(ctx, rateLimitKey); err != nil {
		return "", fmt.Errorf("rate limit exceeded: %w", err)
	}

	provider = strings.ToLower(provider)
	switch provider {
	case "google":
		return k.verifyGoogleToken(ctx, token)
	case "github":
		return k.verifyGithubToken(ctx, token)
	case "apple":
		return k.verifyAppleToken(ctx, token)
	default:
		return "", fmt.Errorf("unsupported OAuth provider: %s", provider)
	}
}
