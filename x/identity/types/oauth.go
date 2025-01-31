package types

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang-jwt/jwt/v4"
)

// OAuth provider constants
const (
	ProviderGoogle = "google"
	ProviderGithub = "github"
)

// OAuthClaims represents the claims in an OAuth2 JWT token
type OAuthClaims struct {
	jwt.RegisteredClaims
	DID      string `json:"did"`
	Provider string `json:"provider"`
}

// Valid implements jwt.Claims interface
func (c *OAuthClaims) Valid() error {
	// Validate standard claims
	if err := c.RegisteredClaims.Valid(); err != nil {
		return sdkerrors.Wrap(ErrInvalidTokenClaims, err.Error())
	}
	
	// Additional validation for our custom claims
	if c.DID == "" {
		return sdkerrors.Wrap(ErrInvalidTokenClaims, "DID cannot be empty")
	}
	if c.Provider == "" {
		return sdkerrors.Wrap(ErrInvalidTokenClaims, "provider cannot be empty")
	}
	
	return nil
}

// ValidateOAuthProvider performs basic validation of OAuth provider
func ValidateOAuthProvider(p *OAuthProvider) error {
	if p.GetName() == "" {
		return sdkerrors.Wrap(ErrInvalidProvider, "name cannot be empty")
	}
	if p.GetClientId() == "" {
		return sdkerrors.Wrap(ErrInvalidProvider, "client ID cannot be empty")
	}
	if p.GetClientSecret() == "" {
		return sdkerrors.Wrap(ErrInvalidProvider, "client secret cannot be empty")
	}
	if p.GetAuthUrl() == "" {
		return sdkerrors.Wrap(ErrInvalidProvider, "auth URL cannot be empty")
	}
	if p.GetTokenUrl() == "" {
		return sdkerrors.Wrap(ErrInvalidProvider, "token URL cannot be empty")
	}
	if p.GetProfileUrl() == "" {
		return sdkerrors.Wrap(ErrInvalidProvider, "profile URL cannot be empty")
	}
	if len(p.GetScopes()) == 0 {
		return sdkerrors.Wrap(ErrInvalidProvider, "scopes cannot be empty")
	}
	return nil
}

// ValidateOAuthSession performs basic validation of OAuth session
func ValidateOAuthSession(s *OAuthSession) error {
	if s.GetId() == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "session ID cannot be empty")
	}
	if s.GetDid() == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "DID cannot be empty")
	}
	if s.GetProvider() == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "provider cannot be empty")
	}
	if s.GetState() == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "state cannot be empty")
	}
	if s.GetCodeVerifier() == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "code verifier cannot be empty")
	}
	if s.GetCreatedAt() == nil {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "creation time must be set")
	}
	if s.GetExpiresAt() == nil {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "expiration time must be set")
	}
	if s.GetExpiresAt().Before(*s.GetCreatedAt()) {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "expiration time must be after creation time")
	}
	return nil
}

// IsOAuthSessionExpired checks if the OAuth session has expired
func IsOAuthSessionExpired(s *OAuthSession) bool {
	expiresAt := s.GetExpiresAt()
	return expiresAt != nil && time.Now().After(*expiresAt)
}

// IsOAuthSessionActive checks if the OAuth session is active
func IsOAuthSessionActive(s *OAuthSession) bool {
	return !IsOAuthSessionExpired(s)
}
