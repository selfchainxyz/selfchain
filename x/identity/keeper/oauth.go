package keeper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"selfchain/x/identity/types"
)

// Parameters for OAuth providers
const (
	GoogleOAuthAuthURL    = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleOAuthTokenURL   = "https://oauth2.googleapis.com/token"
	GoogleOAuthProfileURL = "https://www.googleapis.com/oauth2/v3/userinfo"

	GithubOAuthAuthURL    = "https://github.com/login/oauth/authorize"
	GithubOAuthTokenURL   = "https://github.com/login/oauth/access_token"
	GithubOAuthProfileURL = "https://api.github.com/user"
)

// SaveOAuthProvider saves OAuth provider configuration
func (k Keeper) SaveOAuthProvider(ctx sdk.Context, provider *types.OAuthProvider) error {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyPrefix(types.OAuthProviderPrefix)
	key = append(key, []byte(provider.Id)...)

	bz := k.cdc.MustMarshal(provider)
	store.Set(key, bz)
	return nil
}

// GetOAuthSession gets an OAuth session
func (k Keeper) GetOAuthSession(ctx sdk.Context, sessionID string) (*types.OAuthSession, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyPrefix(types.OAuthSessionPrefix)
	key = append(key, []byte(sessionID)...)

	if !store.Has(key) {
		return nil, false
	}

	var session types.OAuthSession
	k.cdc.MustUnmarshal(store.Get(key), &session)
	return &session, true
}

// SaveOAuthSession saves an OAuth session
func (k Keeper) SaveOAuthSession(ctx sdk.Context, session *types.OAuthSession) error {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyPrefix(types.OAuthSessionPrefix)
	key = append(key, []byte(session.Id)...)

	bz := k.cdc.MustMarshal(session)
	store.Set(key, bz)
	return nil
}

// DeleteOAuthSession deletes an OAuth session
func (k Keeper) DeleteOAuthSession(ctx sdk.Context, sessionID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyPrefix(types.OAuthSessionPrefix)
	key = append(key, []byte(sessionID)...)
	store.Delete(key)
}

// GetSocialIdentityBySocialID gets a social identity by social ID
func (k Keeper) GetSocialIdentityBySocialID(ctx sdk.Context, provider string, socialID string) (*types.SocialIdentity, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSocialIdentityBySocialIDKey(provider, socialID)

	if !store.Has(key) {
		return nil, false
	}

	var identity types.SocialIdentity
	k.cdc.MustUnmarshal(store.Get(key), &identity)
	return &identity, true
}

// GetUserInfo gets user info from an OAuth provider using a token
func (k Keeper) GetUserInfo(ctx sdk.Context, provider string, token string) (*types.UserInfo, error) {
	// Get OAuth provider configuration
	oauthProvider, err := k.GetOAuthProvider(ctx, provider)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrOAuthProviderNotFound, "provider %s not found", provider)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request to profile URL
	req, err := http.NewRequest("GET", oauthProvider.ProfileUrl, nil)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOAuthConfig, "failed to create request: %s", err.Error())
	}

	// Add authorization header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrOAuthVerification, "failed to make request: %s", err.Error())
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, sdkerrors.Wrapf(types.ErrInvalidToken, "invalid token, status: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidTokenResponse, "failed to read response: %s", err.Error())
	}

	// Parse user info
	var userInfo types.UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidTokenResponse, "failed to parse user info: %s", err.Error())
	}

	// Verify user info has required fields
	if userInfo.Id == "" {
		return nil, sdkerrors.Wrap(types.ErrInvalidTokenResponse, "user info missing id")
	}

	return &userInfo, nil
}
