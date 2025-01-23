package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

// RecoveryStatus represents the status of a recovery session
type RecoveryStatus int

const (
	RecoveryStatusPending RecoveryStatus = iota
	RecoveryStatusVerified
	RecoveryStatusExpired
)

// ValidateBasic performs basic validation of recovery session
func (s *RecoverySession) ValidateBasic() error {
	if s.Id == "" {
		return errors.Wrap(ErrInvalidRecoverySession, "session ID cannot be empty")
	}
	if s.Did == "" {
		return errors.Wrap(ErrInvalidRecoverySession, "DID cannot be empty")
	}
	if s.SocialProvider == "" {
		return errors.Wrap(ErrInvalidRecoverySession, "social provider cannot be empty")
	}
	if s.SocialId == "" {
		return errors.Wrap(ErrInvalidRecoverySession, "social ID cannot be empty")
	}
	if s.CreatedAt.IsZero() {
		return errors.Wrap(ErrInvalidRecoverySession, "creation time must be set")
	}
	if s.ExpiresAt.IsZero() {
		return errors.Wrap(ErrInvalidRecoverySession, "expiry time must be set")
	}
	if s.ExpiresAt.Before(s.CreatedAt) {
		return errors.Wrap(ErrInvalidRecoverySession, "expiry time cannot be before creation time")
	}
	return nil
}

// IsExpired checks if the recovery session has expired
func (s *RecoverySession) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}
