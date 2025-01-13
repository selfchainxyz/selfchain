package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateCredential     = "create_credential"
	TypeMsgUpdateCredential     = "update_credential"
	TypeMsgRevokeCredential     = "revoke_credential"
	TypeMsgUnlinkSocialIdentity = "unlink_social_identity"
	TypeMsgConfigureMFA         = "configure_mfa"
	TypeMsgVerifyMFA            = "verify_mfa"
)

var (
	_ sdk.Msg = &MsgCreateCredential{}
	_ sdk.Msg = &MsgUpdateCredential{}
	_ sdk.Msg = &MsgRevokeCredential{}
	_ sdk.Msg = &MsgUnlinkSocialIdentity{}
	_ sdk.Msg = &MsgConfigureMFA{}
	_ sdk.Msg = &MsgVerifyMFA{}
)

// MsgConfigureMFA defines the ConfigureMFA message
type MsgConfigureMFA struct {
	Creator string   `json:"creator"`
	Did     string   `json:"did"`
	Methods []string `json:"methods"`
}

// MsgConfigureMFAResponse defines the response for ConfigureMFA
type MsgConfigureMFAResponse struct {
	Success bool `json:"success"`
}

// Route implements sdk.Msg
func (msg *MsgCreateCredential) Route() string { return RouterKey }

// MsgType implements sdk.Msg
func (msg *MsgCreateCredential) MsgType() string { return TypeMsgCreateCredential }

// ValidateBasic implements sdk.Msg
func (msg *MsgCreateCredential) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	if msg.Id == "" {
		return sdkerrors.Wrap(ErrInvalidCredentialID, "credential ID cannot be empty")
	}
	if msg.GetType() == "" {
		return sdkerrors.Wrap(ErrInvalidCredentialType, "credential type cannot be empty")
	}
	if msg.Issuer == "" {
		return sdkerrors.Wrap(ErrInvalidCredentialIssuer, "issuer cannot be empty")
	}
	if msg.Subject == "" {
		return sdkerrors.Wrap(ErrInvalidCredentialSubject, "subject cannot be empty")
	}
	if len(msg.Claims) == 0 {
		return sdkerrors.Wrap(ErrInvalidCredentialClaims, "claims cannot be empty")
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgCreateCredential) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// Route implements sdk.Msg
func (msg *MsgUpdateCredential) Route() string { return RouterKey }

// MsgType implements sdk.Msg
func (msg *MsgUpdateCredential) MsgType() string { return TypeMsgUpdateCredential }

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateCredential) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	if msg.Id == "" {
		return sdkerrors.Wrap(ErrInvalidCredentialID, "credential ID cannot be empty")
	}
	if len(msg.Claims) == 0 {
		return sdkerrors.Wrap(ErrInvalidCredentialClaims, "claims cannot be empty")
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateCredential) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// Route implements sdk.Msg
func (msg *MsgConfigureMFA) Route() string { return RouterKey }

// MsgType implements sdk.Msg
func (msg *MsgConfigureMFA) MsgType() string { return TypeMsgConfigureMFA }

// ValidateBasic implements sdk.Msg
func (msg *MsgConfigureMFA) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	if len(msg.Methods) == 0 {
		return sdkerrors.Wrap(ErrInvalidMFAMethod, "at least one MFA method is required")
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgConfigureMFA) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// String implements fmt.Stringer
func (msg *MsgConfigureMFA) String() string {
	return fmt.Sprintf("MsgConfigureMFA{Creator: %s, DID: %s, Methods: %v}", msg.Creator, msg.Did, msg.Methods)
}

// ProtoMessage implements proto.Message
func (msg *MsgConfigureMFA) ProtoMessage() {}

// Reset implements proto.Message
func (msg *MsgConfigureMFA) Reset() {
	*msg = MsgConfigureMFA{}
}

// Route implements sdk.Msg
func (msg *MsgVerifyMFA) Route() string { return RouterKey }

// MsgType implements sdk.Msg
func (msg *MsgVerifyMFA) MsgType() string { return TypeMsgVerifyMFA }

// ValidateBasic implements sdk.Msg
func (msg *MsgVerifyMFA) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	// Validate MFA fields
	if msg.Method == "" {
		return sdkerrors.Wrap(ErrInvalidMFAMethod, "method cannot be empty")
	}
	if msg.Code == "" {
		return sdkerrors.Wrap(ErrInvalidMFAMethod, "verification code cannot be empty")
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgVerifyMFA) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// Route implements sdk.Msg
func (msg *MsgRevokeCredential) Route() string { return RouterKey }

// MsgType implements sdk.Msg
func (msg *MsgRevokeCredential) MsgType() string { return TypeMsgRevokeCredential }

// ValidateBasic implements sdk.Msg
func (msg *MsgRevokeCredential) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	if msg.Id == "" {
		return sdkerrors.Wrap(ErrInvalidCredentialID, "credential ID cannot be empty")
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgRevokeCredential) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
