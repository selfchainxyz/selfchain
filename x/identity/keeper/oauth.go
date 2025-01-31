package keeper

import (
	"fmt"

	"selfchain/x/identity/config"
	"selfchain/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// OAuthConfig is the global OAuth configuration
var oauthConfig *config.Config

// InitOAuthConfig initializes the OAuth configuration
func InitOAuthConfig() error {
	var err error
	oauthConfig, err = config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load OAuth configuration: %w", err)
	}
	return nil
}

// LinkSocialIdentity links a social identity to a DID
func (k Keeper) LinkSocialIdentity(ctx sdk.Context, did string, provider string, userInfo *types.UserInfo) error {
	// Check if DID exists
	_, found := k.GetDIDDocument(ctx, did)
	if !found {
		return sdkerrors.Wrapf(types.ErrDIDNotFound, "did %s not found", did)
	}

	// Check if social identity already exists
	key := fmt.Sprintf("%s:%s:%s", did, provider, userInfo.Id)
	if _, found := k.GetSocialIdentity(ctx, did, provider); found {
		return sdkerrors.Wrapf(types.ErrSocialIdentityExists, "social identity %s already exists for did %s", key, did)
	}

	// Get current block time
	now := ctx.BlockTime()

	// Create new social identity
	socialIdentity := types.SocialIdentity{
		Id:         key,
		Did:        did,
		Provider:   provider,
		ProviderId: userInfo.Id,
		Profile: map[string]string{
			"username": userInfo.Username,
			"email":    userInfo.Email,
		},
		CreatedAt:  &now,
		VerifiedAt: &now,
		LastUsed:   &now,
	}

	// Store social identity
	store := k.GetStore(ctx, []byte(types.SocialIdentityPrefix))
	bz := k.cdc.MustMarshal(&socialIdentity)
	store.Set([]byte(key), bz)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOAuthSuccess,
			sdk.NewAttribute(types.AttributeKeyDID, did),
			sdk.NewAttribute(types.AttributeKeyProvider, provider),
			sdk.NewAttribute(types.AttributeKeySocialID, userInfo.Id),
			sdk.NewAttribute(types.AttributeKeyStatus, "linked"),
		),
	)

	return nil
}

// UnlinkSocialIdentity unlinks a social identity from a DID
func (k Keeper) UnlinkSocialIdentity(ctx sdk.Context, did string, provider string, socialId string) error {
	// Check if DID exists
	_, found := k.GetDIDDocument(ctx, did)
	if !found {
		return sdkerrors.Wrapf(types.ErrDIDNotFound, "did %s not found", did)
	}

	// Check if social identity exists
	key := fmt.Sprintf("%s:%s:%s", did, provider, socialId)
	if _, found := k.GetSocialIdentity(ctx, did, provider); !found {
		return sdkerrors.Wrapf(types.ErrSocialIdentityNotFound, "social identity %s not found", key)
	}

	// Delete social identity
	store := k.GetStore(ctx, []byte(types.SocialIdentityPrefix))
	store.Delete([]byte(key))

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOAuthSuccess,
			sdk.NewAttribute(types.AttributeKeyDID, did),
			sdk.NewAttribute(types.AttributeKeyProvider, provider),
			sdk.NewAttribute(types.AttributeKeySocialID, socialId),
			sdk.NewAttribute(types.AttributeKeyStatus, "unlinked"),
		),
	)

	return nil
}

// GetOAuthConfig gets OAuth configuration for a provider
func (k Keeper) GetOAuthConfig(ctx sdk.Context, provider string) (*types.OAuthProvider, error) {
	if oauthConfig == nil {
		if err := InitOAuthConfig(); err != nil {
			return nil, sdkerrors.Wrapf(types.ErrInvalidOAuthConfig, "OAuth configuration not initialized: %s", err.Error())
		}
	}

	switch provider {
	case "google":
		return &types.OAuthProvider{
			Id:           "google",
			Name:         "Google",
			ClientId:     oauthConfig.GoogleOAuth.ClientID,
			ClientSecret: oauthConfig.GoogleOAuth.ClientSecret,
			AuthUrl:      oauthConfig.GoogleOAuth.AuthURL,
			TokenUrl:     oauthConfig.GoogleOAuth.TokenURL,
			ProfileUrl:   oauthConfig.GoogleOAuth.ProfileURL,
			Scopes:       oauthConfig.GoogleOAuth.Scopes,
			Config:       make(map[string]string),
		}, nil
	case "github":
		return &types.OAuthProvider{
			Id:           "github",
			Name:         "GitHub",
			ClientId:     oauthConfig.GithubOAuth.ClientID,
			ClientSecret: oauthConfig.GithubOAuth.ClientSecret,
			AuthUrl:      oauthConfig.GithubOAuth.AuthURL,
			TokenUrl:     oauthConfig.GithubOAuth.TokenURL,
			ProfileUrl:   oauthConfig.GithubOAuth.ProfileURL,
			Scopes:       oauthConfig.GithubOAuth.Scopes,
			Config:       make(map[string]string),
		}, nil
	default:
		return nil, sdkerrors.Wrapf(types.ErrInvalidProvider, "unsupported OAuth provider: %s", provider)
	}
}

// Parameters for OAuth providers
const (
	GoogleOAuthAuthURL    = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleOAuthTokenURL   = "https://oauth2.googleapis.com/token"
	GoogleOAuthProfileURL = "https://www.googleapis.com/oauth2/v3/userinfo"

	GithubOAuthAuthURL    = "https://github.com/login/oauth/authorize"
	GithubOAuthTokenURL   = "https://github.com/login/oauth/access_token"
	GithubOAuthProfileURL = "https://api.github.com/user"
)
