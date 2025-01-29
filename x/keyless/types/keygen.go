package types

import (
	"fmt"
)

// Validate validates a KeyGenRequest
func (r *KeyGenRequest) Validate() error {
	if r.WalletAddress == "" {
		return fmt.Errorf("wallet address cannot be empty")
	}
	if r.ChainId == "" {
		return fmt.Errorf("chain id cannot be empty")
	}
	if r.SecurityLevel == SecurityLevel_SECURITY_LEVEL_UNSPECIFIED {
		return fmt.Errorf("security level must be specified")
	}
	return nil
}

// ValidateBasic performs basic validation of an EncryptedShare
func (s *EncryptedShare) ValidateBasic() error {
	if s.EncryptedData == "" {
		return ErrInvalidShare.Wrap("encrypted data is required")
	}
	if s.KeyId == "" {
		return ErrInvalidShare.Wrap("key ID is required")
	}
	if s.Version == 0 {
		return ErrInvalidShare.Wrap("version is required")
	}
	if s.CreatedAt.IsZero() {
		return ErrInvalidShare.Wrap("creation time is required")
	}
	return nil
}

// Validate validates a KeyGenResponse
func (r *KeyGenResponse) Validate() error {
	if r.WalletAddress == "" {
		return fmt.Errorf("wallet address cannot be empty")
	}
	if r.PublicKey == nil {
		return fmt.Errorf("public key cannot be nil")
	}
	if r.Metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}
	return nil
}
