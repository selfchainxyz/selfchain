package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRelease = "release"

var _ sdk.Msg = &MsgRelease{}

func NewMsgRelease(creator string, posIndex uint64) *MsgRelease {
	return &MsgRelease{
		Creator:  creator,
		PosIndex: posIndex,
	}
}

func (msg *MsgRelease) Route() string {
	return RouterKey
}

func (msg *MsgRelease) Type() string {
	return TypeMsgRelease
}

func (msg *MsgRelease) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRelease) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRelease) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
