package types

// Validate validates a KeyGenRequest
func (r *KeyGenRequest) Validate() error {
	if r.WalletId == "" {
		return ErrInvalidRequest.Wrap("wallet ID is required")
	}
	if r.SecurityLevel == "" {
		return ErrInvalidRequest.Wrap("security level is required")
	}
	return nil
}

// ValidateBasic performs basic validation of an EncryptedShare
func (s *EncryptedShare) ValidateBasic() error {
	if s.Data == nil {
		return ErrInvalidShare.Wrap("encrypted data is required")
	}
	if s.PublicKey == nil {
		return ErrInvalidShare.Wrap("public key is required")
	}
	if s.Nonce == nil {
		return ErrInvalidShare.Wrap("nonce is required")
	}
	return nil
}

// ValidateBasic performs basic validation of a KeyGenResponse
func (r *KeyGenResponse) ValidateBasic() error {
	if r.WalletId == "" {
		return ErrInvalidResponse.Wrap("wallet ID is required")
	}
	if r.PublicKey == nil {
		return ErrInvalidResponse.Wrap("public key is required")
	}
	if r.Metadata == nil {
		return ErrInvalidResponse.Wrap("metadata is required")
	}
	return nil
}
