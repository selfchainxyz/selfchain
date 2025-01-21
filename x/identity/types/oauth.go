package types

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidOAuthProvider   = sdkerrors.Register(ModuleName, 1300, "invalid OAuth provider")
	ErrInvalidOAuthSession    = sdkerrors.Register(ModuleName, 1301, "invalid OAuth session")
	ErrOAuthSessionExpired    = sdkerrors.Register(ModuleName, 1303, "OAuth session expired")
	ErrOAuthProviderNotFound  = sdkerrors.Register(ModuleName, 1304, "OAuth provider not found")
	ErrOAuthSessionNotFound   = sdkerrors.Register(ModuleName, 1305, "OAuth session not found")
	ErrSocialIdentityNotFound = sdkerrors.Register(ModuleName, 1306, "social identity not found")
)

// ValidateBasic performs basic validation of social identity
func (s *SocialIdentity) ValidateBasic() error {
	if s.Id == "" {
		return sdkerrors.Wrap(ErrInvalidSocialIdentity, "ID cannot be empty")
	}
	if s.Did == "" {
		return sdkerrors.Wrap(ErrInvalidSocialIdentity, "DID cannot be empty")
	}
	if s.Provider == "" {
		return sdkerrors.Wrap(ErrInvalidSocialIdentity, "provider cannot be empty")
	}
	if s.ProviderId == "" {
		return sdkerrors.Wrap(ErrInvalidSocialIdentity, "provider ID cannot be empty")
	}
	if s.CreatedAt == nil {
		return sdkerrors.Wrap(ErrInvalidSocialIdentity, "creation time must be set")
	}
	return nil
}

// IsVerified checks if the social identity is verified
func (s *SocialIdentity) IsVerified() bool {
	return s.VerifiedAt != nil
}

// ValidateBasic performs basic validation of OAuth provider
func (p *OAuthProvider) ValidateBasic() error {
	if p.Id == "" {
		return sdkerrors.Wrap(ErrInvalidOAuthProvider, "ID cannot be empty")
	}
	if p.Name == "" {
		return sdkerrors.Wrap(ErrInvalidOAuthProvider, "name cannot be empty")
	}
	if p.ClientId == "" {
		return sdkerrors.Wrap(ErrInvalidOAuthProvider, "client ID cannot be empty")
	}
	if p.ClientSecret == "" {
		return sdkerrors.Wrap(ErrInvalidOAuthProvider, "client secret cannot be empty")
	}
	if p.AuthUrl == "" {
		return sdkerrors.Wrap(ErrInvalidOAuthProvider, "auth URL cannot be empty")
	}
	if p.TokenUrl == "" {
		return sdkerrors.Wrap(ErrInvalidOAuthProvider, "token URL cannot be empty")
	}
	if p.ProfileUrl == "" {
		return sdkerrors.Wrap(ErrInvalidOAuthProvider, "profile URL cannot be empty")
	}
	return nil
}

// ValidateBasic performs basic validation of OAuth session
func (s *OAuthSession) ValidateBasic() error {
	if s.Id == "" {
		return sdkerrors.Wrap(ErrInvalidOAuthSession, "ID cannot be empty")
	}
	if s.Did == "" {
		return sdkerrors.Wrap(ErrInvalidOAuthSession, "DID cannot be empty")
	}
	if s.Provider == "" {
		return sdkerrors.Wrap(ErrInvalidOAuthSession, "provider cannot be empty")
	}
	if s.State == "" {
		return sdkerrors.Wrap(ErrInvalidOAuthSession, "state cannot be empty")
	}
	if s.CodeVerifier == "" {
		return sdkerrors.Wrap(ErrInvalidOAuthSession, "code verifier cannot be empty")
	}
	if s.CreatedAt == nil || s.ExpiresAt == nil {
		return sdkerrors.Wrap(ErrInvalidOAuthSession, "timestamps must be set")
	}
	now := time.Date(2025, 1, 12, 15, 5, 4, 0, time.FixedZone("IST", 5*60*60+30*60))
	if s.ExpiresAt.Before(now) {
		return ErrOAuthSessionExpired
	}
	return nil
}

// IsExpired checks if the OAuth session has expired
func (s *OAuthSession) IsExpired() bool {
	if s.ExpiresAt == nil {
		return true
	}
	now := time.Date(2025, 1, 12, 15, 5, 4, 0, time.FixedZone("IST", 5*60*60+30*60))
	return s.ExpiresAt.Before(now)
}
