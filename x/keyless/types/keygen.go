package types

// Validate validates a KeyGenRequest
func (r *KeyGenRequest) Validate() error {
	if r.WalletId == "" {
		return ErrInvalidRequest.Wrap("wallet ID is required")
	}
	if r.ChainId == "" {
		return ErrInvalidRequest.Wrap("chain ID is required")
	}
	if !r.SecurityLevel.IsValid() {
		return ErrInvalidRequest.Wrap("invalid security level")
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

// ValidateBasic performs basic validation of a KeyGenResponse
func (r *KeyGenResponse) ValidateBasic() error {
	if r.WalletId == "" {
		return ErrInvalidResponse.Wrap("wallet ID is required")
	}
	if len(r.PublicKey) == 0 {
		return ErrInvalidResponse.Wrap("public key is required")
	}
	if r.Metadata == nil {
		return ErrInvalidResponse.Wrap("metadata is required")
	}
	return nil
}
