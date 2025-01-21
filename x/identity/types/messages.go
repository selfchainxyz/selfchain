package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateCredential     = "create_credential"
	TypeMsgUpdateCredential     = "update_credential"
	TypeMsgRevokeCredential     = "revoke_credential"
	TypeMsgUnlinkSocialIdentity = "unlink_social_identity"
)

var (
	_ sdk.Msg = &MsgCreateCredential{}
	_ sdk.Msg = &MsgUpdateCredential{}
	_ sdk.Msg = &MsgRevokeCredential{}
	_ sdk.Msg = &MsgUnlinkSocialIdentity{}
)

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
