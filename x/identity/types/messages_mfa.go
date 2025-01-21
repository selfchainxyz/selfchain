package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgAddMFA    = "add_mfa"
	TypeMsgRemoveMFA = "remove_mfa"
	TypeMsgVerifyMFA = "verify_mfa"
)

var _ sdk.Msg = &MsgAddMFA{}
var _ sdk.Msg = &MsgRemoveMFA{}
var _ sdk.Msg = &MsgVerifyMFA{}

// NewMsgAddMFA creates a new MsgAddMFA instance
func NewMsgAddMFA(creator string, did string, method string, secret string) *MsgAddMFA {
	return &MsgAddMFA{
		Creator: creator,
		Did:     did,
		Method:  method,
		Secret:  secret,
	}
}

func (msg *MsgAddMFA) Route() string {
	return RouterKey
}

func (msg *MsgAddMFA) Type() string {
	return TypeMsgAddMFA
}

func (msg *MsgAddMFA) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddMFA) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddMFA) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Did == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "DID cannot be empty")
	}
	if msg.Method == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "method cannot be empty")
	}
	if msg.Secret == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "secret cannot be empty")
	}
	return nil
}

// NewMsgRemoveMFA creates a new MsgRemoveMFA instance
func NewMsgRemoveMFA(creator string, did string, method string) *MsgRemoveMFA {
	return &MsgRemoveMFA{
		Creator: creator,
		Did:     did,
		Method:  method,
	}
}

func (msg *MsgRemoveMFA) Route() string {
	return RouterKey
}

func (msg *MsgRemoveMFA) Type() string {
	return TypeMsgRemoveMFA
}

func (msg *MsgRemoveMFA) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveMFA) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveMFA) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Did == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "DID cannot be empty")
	}
	if msg.Method == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "method cannot be empty")
	}
	return nil
}

// NewMsgVerifyMFA creates a new MsgVerifyMFA instance
func NewMsgVerifyMFA(creator string, did string, method string, code string) *MsgVerifyMFA {
	return &MsgVerifyMFA{
		Creator: creator,
		Did:     did,
		Method:  method,
		Code:    code,
	}
}

func (msg *MsgVerifyMFA) Route() string {
	return RouterKey
}

func (msg *MsgVerifyMFA) Type() string {
	return TypeMsgVerifyMFA
}

func (msg *MsgVerifyMFA) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVerifyMFA) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVerifyMFA) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Did == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "DID cannot be empty")
	}
	if msg.Method == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "method cannot be empty")
	}
	if msg.Code == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "code cannot be empty")
	}
	return nil
}
