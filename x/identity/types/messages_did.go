package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateDID = "create_did"
	TypeMsgUpdateDID = "update_did"
	TypeMsgDeleteDID = "delete_did"
)

var (
	_ sdk.Msg = &MsgCreateDID{}
	_ sdk.Msg = &MsgUpdateDID{}
	_ sdk.Msg = &MsgDeleteDID{}
)

// Route implements sdk.Msg
func (msg *MsgCreateDID) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg *MsgCreateDID) Type() string { return TypeMsgCreateDID }

// ValidateBasic implements sdk.Msg
func (msg *MsgCreateDID) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	if msg.Id == "" {
		return sdkerrors.Wrap(ErrInvalidDIDID, "DID identifier cannot be empty")
	}
	if len(msg.Controller) == 0 {
		return sdkerrors.Wrap(ErrInvalidDIDController, "at least one controller is required")
	}

	// Validate verification methods
	for _, vm := range msg.VerificationMethod {
		if err := vm.ValidateBasic(); err != nil {
			return err
		}
	}

	// Validate service endpoints
	for _, svc := range msg.Service {
		if err := svc.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgCreateDID) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// Route implements sdk.Msg
func (msg *MsgUpdateDID) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg *MsgUpdateDID) Type() string { return TypeMsgUpdateDID }

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateDID) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	if msg.Id == "" {
		return sdkerrors.Wrap(ErrInvalidDIDID, "DID identifier cannot be empty")
	}
	if len(msg.Controller) == 0 {
		return sdkerrors.Wrap(ErrInvalidDIDController, "at least one controller is required")
	}

	// Validate verification methods
	for _, vm := range msg.VerificationMethod {
		if err := vm.ValidateBasic(); err != nil {
			return err
		}
	}

	// Validate service endpoints
	for _, svc := range msg.Service {
		if err := svc.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateDID) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// Route implements sdk.Msg
func (msg *MsgDeleteDID) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg *MsgDeleteDID) Type() string { return TypeMsgDeleteDID }

// ValidateBasic implements sdk.Msg
func (msg *MsgDeleteDID) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	if msg.Id == "" {
		return sdkerrors.Wrap(ErrInvalidDIDID, "DID identifier cannot be empty")
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg *MsgDeleteDID) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
