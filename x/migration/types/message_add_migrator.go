package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddMigrator = "add_migrator"

var _ sdk.Msg = &MsgAddMigrator{}

func NewMsgAddMigrator(creator string, migrator string) *MsgAddMigrator {
	return &MsgAddMigrator{
		Creator:  creator,
		Migrator: migrator,
	}
}

func (msg *MsgAddMigrator) Route() string {
	return RouterKey
}

func (msg *MsgAddMigrator) Type() string {
	return TypeMsgAddMigrator
}

func (msg *MsgAddMigrator) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddMigrator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddMigrator) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
