package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewPermission creates a new Permission instance
func NewPermission(walletAddress string, grantee string, permissions []string, expiresAt *time.Time) *Permission {
	now := time.Now()
	return &Permission{
		WalletAddress: walletAddress,
		Grantee:       grantee,
		Permissions:   permissions,
		GrantedAt:     &now,
		ExpiresAt:     expiresAt,
		Revoked:       false,
		RevokedAt:     nil,
	}
}

// Validate validates the permission
func (p Permission) Validate() error {
	if p.WalletAddress == "" {
		return ErrInvalidPermission.Wrap("wallet address cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(p.WalletAddress); err != nil {
		return ErrInvalidPermission.Wrapf("invalid wallet address: %s", err)
	}

	if p.Grantee == "" {
		return ErrInvalidPermission.Wrap("grantee cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(p.Grantee); err != nil {
		return ErrInvalidPermission.Wrapf("invalid grantee address: %s", err)
	}

	if len(p.Permissions) == 0 {
		return ErrInvalidPermission.Wrap("permissions cannot be empty")
	}
	if p.ExpiresAt != nil && p.ExpiresAt.IsZero() {
		return ErrInvalidPermission.Wrap("expiry time cannot be zero")
	}
	return nil
}

// IsExpired checks if the permission has expired
func (p Permission) IsExpired() bool {
	return p.ExpiresAt != nil && !p.ExpiresAt.IsZero() && p.ExpiresAt.Before(time.Now())
}

// IsRevoked checks if the permission has been revoked
func (p Permission) IsRevoked() bool {
	return p.Revoked
}

// Revoke marks the permission as revoked
func (p *Permission) Revoke() {
	p.Revoked = true
	now := time.Now()
	p.RevokedAt = &now
}

// HasPermission checks if the permission contains a specific permission type
func (p Permission) HasPermission(permType string) bool {
	if p.IsExpired() || p.IsRevoked() {
		return false
	}

	for _, perm := range p.Permissions {
		if perm == permType {
			return true
		}
	}
	return false
}
