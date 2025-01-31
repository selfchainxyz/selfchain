package types

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Recovery status constants
const (
	RecoveryStatusPending  = "pending"
	RecoveryStatusActive   = "active"
	RecoveryStatusComplete = "complete"
	RecoveryStatusExpired  = "expired"
	RecoveryStatusFailed   = "failed"
)

// IsExpired checks if the recovery session has expired
func (s *RecoverySession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsActive checks if the recovery session is active
func (s *RecoverySession) IsActive() bool {
	return !s.IsExpired() && !s.MfaVerified && !s.IdentityVerified
}

// IsComplete checks if the recovery session is complete
func (s *RecoverySession) IsComplete() bool {
	return s.MfaVerified && s.IdentityVerified
}

// ValidateBasic performs basic validation of recovery session
func (s *RecoverySession) ValidateBasic() error {
	if s.Id == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "session ID cannot be empty")
	}
	if s.Did == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "DID cannot be empty")
	}
	if s.SocialProvider == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "social provider cannot be empty")
	}
	if s.SocialId == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "social ID cannot be empty")
	}
	if s.CreatedAt.IsZero() {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "creation time must be set")
	}
	if s.ExpiresAt.IsZero() {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "expiration time must be set")
	}
	if s.ExpiresAt.Before(s.CreatedAt) {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "expiration time must be after creation time")
	}
	if len(s.RecoveryData) == 0 {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "recovery data cannot be empty")
	}
	return nil
}

// ValidateBasic performs basic validation of recovery data
func (d *RecoveryData) ValidateBasic() error {
	if d.WalletId == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "wallet ID cannot be empty")
	}
	if d.Token == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "token cannot be empty")
	}
	if d.ExpiresAt == nil {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "expiration time must be set")
	}
	if d.ExpiresAt.Before(time.Now()) {
		return sdkerrors.Wrap(ErrRecoverySessionExpired, "recovery data has expired")
	}
	return nil
}

// ValidateBasic performs basic validation of MsgInitiateRecovery
func (msg *MsgInitiateRecovery) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	if msg.SocialProvider == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "social provider cannot be empty")
	}
	if msg.OauthToken == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "OAuth token cannot be empty")
	}
	return nil
}

// GetSigners returns the expected signers for MsgInitiateRecovery
func (msg *MsgInitiateRecovery) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ValidateBasic performs basic validation of MsgCompleteRecovery
func (msg *MsgCompleteRecovery) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	if msg.SessionId == "" {
		return sdkerrors.Wrap(ErrInvalidRecoverySession, "session ID cannot be empty")
	}
	if msg.MfaCode == "" {
		return sdkerrors.Wrap(ErrInvalidRecoveryCode, "MFA code cannot be empty")
	}
	return nil
}

// GetSigners returns the expected signers for MsgCompleteRecovery
func (msg *MsgCompleteRecovery) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
