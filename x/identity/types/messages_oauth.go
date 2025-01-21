package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgLinkSocialIdentity = "link_social_identity"
	TypeMsgVerifyOAuthToken   = "verify_oauth_token"
)

var _ sdk.Msg = &MsgLinkSocialIdentity{}
var _ sdk.Msg = &MsgUnlinkSocialIdentity{}
var _ sdk.Msg = &MsgVerifyOAuthToken{}

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

// NewMsgVerifyOAuthToken creates a new MsgVerifyOAuthToken instance
func NewMsgVerifyOAuthToken(creator string, provider string, token string) *MsgVerifyOAuthToken {
	return &MsgVerifyOAuthToken{
		Creator:  creator,
		Provider: provider,
		Token:    token,
	}
}

func (msg *MsgVerifyOAuthToken) Route() string {
	return RouterKey
}

func (msg *MsgVerifyOAuthToken) Type() string {
	return TypeMsgVerifyOAuthToken
}

func (msg *MsgVerifyOAuthToken) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVerifyOAuthToken) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVerifyOAuthToken) ValidateBasic() error {
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
