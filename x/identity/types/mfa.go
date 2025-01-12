package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

// MFAMethodStatus values
const (
	MFAMethodStatusUnspecified = MFAMethodStatus_MFA_METHOD_STATUS_UNSPECIFIED
	MFAMethodStatusActive     = MFAMethodStatus_MFA_METHOD_STATUS_ACTIVE
	MFAMethodStatusPending    = MFAMethodStatus_MFA_METHOD_STATUS_PENDING
	MFAMethodStatusDisabled   = MFAMethodStatus_MFA_METHOD_STATUS_DISABLED
)

// ValidateBasic performs basic validation of MFA configuration
func (m *MFAConfig) ValidateBasic() error {
	if m.Did == "" {
		return errors.Wrap(ErrInvalidMFAMethod, "DID cannot be empty")
	}
	if len(m.Methods) == 0 {
		return errors.Wrap(ErrInvalidMFAMethod, "at least one MFA method is required")
	}
	for _, method := range m.Methods {
		if err := method.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

// ValidateBasic performs basic validation of MFA method
func (m *MFAMethod) ValidateBasic() error {
	if m.Type == "" {
		return ErrInvalidMFAMethod
	}
	if m.Secret == "" {
		return ErrInvalidMFAMethod
	}
	if m.CreatedAt == nil {
		return errors.Wrap(ErrInvalidMFAMethod, "creation time must be set")
	}
	return nil
}

// IsActive checks if the MFA method is active
func (m *MFAMethod) IsActive() bool {
	return m.Status == MFAMethodStatus_MFA_METHOD_STATUS_ACTIVE
}

// IsPending checks if the MFA method is pending
func (m *MFAMethod) IsPending() bool {
	return m.Status == MFAMethodStatus_MFA_METHOD_STATUS_PENDING
}

// IsDisabled checks if the MFA method is disabled
func (m *MFAMethod) IsDisabled() bool {
	return m.Status == MFAMethodStatus_MFA_METHOD_STATUS_DISABLED
}

// ValidateBasic performs basic validation of MFA challenge
func (m *MFAChallenge) ValidateBasic() error {
	if m.Id == "" {
		return errors.Wrap(ErrInvalidMFAChallenge, "ID cannot be empty")
	}
	if m.Did == "" {
		return errors.Wrap(ErrInvalidMFAChallenge, "DID cannot be empty")
	}
	if m.Method == "" {
		return errors.Wrap(ErrInvalidMFAChallenge, "method cannot be empty")
	}
	if m.CreatedAt == nil {
		return errors.Wrap(ErrInvalidMFAChallenge, "creation time must be set")
	}
	if m.ExpiresAt == nil {
		return errors.Wrap(ErrInvalidMFAChallenge, "expiry time must be set")
	}
	if m.ExpiresAt.Before(*m.CreatedAt) {
		return errors.Wrap(ErrInvalidMFAChallenge, "expiry time must be after creation time")
	}
	return nil
}

// IsExpired checks if the MFA challenge has expired
func (m *MFAChallenge) IsExpired() bool {
	if m.ExpiresAt == nil {
		return true
	}
	return time.Now().After(*m.ExpiresAt)
}
