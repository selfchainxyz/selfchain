package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgLinkSocialIdentity = "link_social_identity"
)

var _ sdk.Msg = &MsgLinkSocialIdentity{}

// NewMsgLinkSocialIdentity creates a new MsgLinkSocialIdentity instance
func NewMsgLinkSocialIdentity(creator string, provider string, token string) *MsgLinkSocialIdentity {
	return &MsgLinkSocialIdentity{
		Creator:  creator,
		Provider: provider,
		Token:    token,
	}
}

func (msg *MsgLinkSocialIdentity) Route() string {
	return RouterKey
}

func (msg *MsgLinkSocialIdentity) Type() string {
	return TypeMsgLinkSocialIdentity
}

func (msg *MsgLinkSocialIdentity) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgLinkSocialIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgLinkSocialIdentity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Provider == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "provider cannot be empty")
	}
	if msg.Token == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token cannot be empty")
	}
	return nil
}

var _ sdk.Msg = &MsgUnlinkSocialIdentity{}

// NewMsgUnlinkSocialIdentity creates a new MsgUnlinkSocialIdentity instance
func NewMsgUnlinkSocialIdentity(creator string, provider string) *MsgUnlinkSocialIdentity {
	return &MsgUnlinkSocialIdentity{
		Creator:  creator,
		Provider: provider,
	}
}

func (msg *MsgUnlinkSocialIdentity) Route() string {
	return RouterKey
}

func (msg *MsgUnlinkSocialIdentity) Type() string {
	return "unlink_social_identity"
}

func (msg *MsgUnlinkSocialIdentity) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUnlinkSocialIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnlinkSocialIdentity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Provider == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "provider cannot be empty")
	}
	return nil
}
