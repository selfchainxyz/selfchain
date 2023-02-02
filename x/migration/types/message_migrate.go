package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	uint256 "github.com/holiman/uint256"
)

const TypeMsgMigrate = "migrate"

var _ sdk.Msg = &MsgMigrate{}

func NewMsgMigrate(creator string, ethAddress string, destAddress string, amount string, token uint64) *MsgMigrate {
	return &MsgMigrate{
		Creator:     creator,
		EthAddress:  ethAddress,
		DestAddress: destAddress,
		Amount:      amount,
		Token:       token,
	}
}

func (msg *MsgMigrate) Route() string {
	return RouterKey
}

func (msg *MsgMigrate) Type() string {
	return TypeMsgMigrate
}

func (msg *MsgMigrate) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMigrate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMigrate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator); if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err2 := sdk.AccAddressFromBech32(msg.DestAddress); if err2 != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid destination address (%s)", err2)
	}

	// we don't want to get spammed people who migrate small amounts
	amount, err := uint256.FromHex(msg.Amount); if amount.Lt(getMinMigrationAmount()) || err != nil {
		return ErrInvalidMigrationAmount
	}

	// check that token is supported
	if msg.Token != uint64(Front) && msg.Token != uint64(Hotcross) {
		return ErrTokenNotSupported
	}

	if msg.TxHash == "" || msg.EthAddress == "" {
		return ErrEmptyStringValue
	}

	return nil
}
