package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateConfig = "update_config"

var _ sdk.Msg = &MsgUpdateConfig{}

func NewMsgUpdateConfig(creator string, vestingDuration uint64, vestingCliff uint64, minMigrationAmount uint64) *MsgUpdateConfig {
	return &MsgUpdateConfig{
		Creator:            creator,
		VestingDuration:    vestingDuration,
		VestingCliff:       vestingCliff,
		MinMigrationAmount: minMigrationAmount,
	}
}

func (msg *MsgUpdateConfig) Route() string {
	return RouterKey
}

func (msg *MsgUpdateConfig) Type() string {
	return TypeMsgUpdateConfig
}

func (msg *MsgUpdateConfig) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateConfig) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateConfig) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
