package types

import (
	"fmt"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CredentialStatus represents the status of a credential
type CredentialStatus string

const (
	// CredentialStatusActive indicates an active credential
	CredentialStatusActive CredentialStatus = "ACTIVE"
	// CredentialStatusRevoked indicates a revoked credential
	CredentialStatusRevoked CredentialStatus = "REVOKED"
)

// String returns the string representation of CredentialStatus
func (s CredentialStatus) String() string {
	return string(s)
}

// IsValid returns true if the status is valid
func (s CredentialStatus) IsValid() bool {
	switch s {
	case CredentialStatusActive, CredentialStatusRevoked:
		return true
	default:
		return false
	}
}

// CredentialType represents the type of a credential
type CredentialType string

const (
	// CredentialTypeVerifiableCredential is a W3C Verifiable Credential
	CredentialTypeVerifiableCredential CredentialType = "VerifiableCredential"
)

// String returns the string representation of CredentialType
func (t CredentialType) String() string {
	return string(t)
}

// IsValid returns true if the type is valid
func (t CredentialType) IsValid() bool {
	switch t {
	case CredentialTypeVerifiableCredential:
		return true
	default:
		return false
	}
}

// ValidateBasic performs basic validation of the credential proof
func (p *CredentialProof) ValidateBasic() error {
	if p.Type == "" {
		return sdkerrors.Register(ModuleName, 1300, "invalid proof type")
	}
	if p.Created == 0 {
		return sdkerrors.Register(ModuleName, 1301, "invalid proof creation time")
	}
	if p.VerificationMethod == "" {
		return sdkerrors.Register(ModuleName, 1302, "invalid verification method")
	}
	if p.ProofPurpose == "" {
		return sdkerrors.Register(ModuleName, 1303, "invalid proof purpose")
	}
	return nil
}

// ValidateBasic performs basic validation of the credential presentation
func (p *CredentialPresentation) ValidateBasic() error {
	if p.Type == "" {
		return sdkerrors.Register(ModuleName, 1200, "invalid presentation type")
	}
	if p.Created == 0 {
		return sdkerrors.Register(ModuleName, 1201, "invalid presentation creation time")
	}
	if p.VerifiableCredential == "" {
		return sdkerrors.Register(ModuleName, 1202, "invalid verifiable credential")
	}
	if p.Proof == nil {
		return sdkerrors.Register(ModuleName, 1203, "missing presentation proof")
	}
	if err := p.Proof.ValidateBasic(); err != nil {
		return sdkerrors.Wrapf(err, "invalid presentation proof")
	}
	return nil
}

// ValidateBasic performs basic validation of the credential
func (c *Credential) ValidateBasic() error {
	if c == nil {
		return sdkerrors.Register(ModuleName, 1200, "credential cannot be nil")
	}

	if c.Id == "" {
		return sdkerrors.Register(ModuleName, 1201, "credential ID cannot be empty")
	}

	if c.Issuer == "" {
		return sdkerrors.Register(ModuleName, 1202, "issuer cannot be empty")
	}

	if c.Subject == "" {
		return sdkerrors.Register(ModuleName, 1203, "subject cannot be empty")
	}

	if !CredentialStatus(c.Status).IsValid() {
		return sdkerrors.Register(ModuleName, 1204, fmt.Sprintf("invalid credential status: %s", c.Status))
	}

	if c.IssuanceDate == 0 {
		return sdkerrors.Register(ModuleName, 1205, "issuance date must be set")
	}

	// If expiration date is set, ensure it's after issuance date
	if c.ExpirationDate != 0 {
		issuanceTime := time.Unix(c.IssuanceDate, 0)
		expirationTime := time.Unix(c.ExpirationDate, 0)
		if expirationTime.Before(issuanceTime) {
			return sdkerrors.Register(ModuleName, 1206, "expiration date must be after issuance date")
		}
	}

	// Validate proof if present
	if c.Proof != nil {
		if err := c.Proof.ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "invalid credential proof")
		}
	}

	return nil
}

// NewCredential creates a new credential
func NewCredential(
	id string,
	issuer string,
	subject string,
	claims map[string]string,
	proof *CredentialProof,
) *Credential {
	return &Credential{
		Id:           id,
		Issuer:       issuer,
		Subject:      subject,
		Claims:       claims,
		Status:       string(CredentialStatusActive),
		IssuanceDate: time.Now().Unix(),
		Proof:        proof,
	}
}

// IsRevoked checks if the credential is revoked
func (c *Credential) IsRevoked() bool {
	return c.Status == string(CredentialStatusRevoked)
}

// IsExpired checks if the credential has expired
func (c *Credential) IsExpired(now int64) bool {
	return c.ExpirationDate != 0 && c.ExpirationDate < now
}

// IsValid checks if the credential is valid (not revoked and not expired)
func (c *Credential) IsValid(now int64) bool {
	return !c.IsRevoked() && !c.IsExpired(now)
}
