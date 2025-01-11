package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgLinkSocialIdentity = "link_social_identity"
)

var _ sdk.Msg = &MsgLinkSocialIdentity{}

func NewMsgLinkSocialIdentity(did string, provider string, token string) *MsgLinkSocialIdentity {
	return &MsgLinkSocialIdentity{
		Did:      did,
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
	return []sdk.AccAddress{}
}

func (msg *MsgLinkSocialIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgLinkSocialIdentity) ValidateBasic() error {
	if msg.Did == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid DID address: %s", msg.Did)
	}
	if msg.Provider == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "provider cannot be empty")
	}
	if msg.Token == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token cannot be empty")
	}
	return nil
}
