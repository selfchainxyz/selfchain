package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveMigrator = "remove_migrator"

var _ sdk.Msg = &MsgRemoveMigrator{}

func NewMsgRemoveMigrator(creator string, migrator string) *MsgRemoveMigrator {
	return &MsgRemoveMigrator{
		Creator:  creator,
		Migrator: migrator,
	}
}

func (msg *MsgRemoveMigrator) Route() string {
	return RouterKey
}

func (msg *MsgRemoveMigrator) Type() string {
	return TypeMsgRemoveMigrator
}

func (msg *MsgRemoveMigrator) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveMigrator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveMigrator) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
