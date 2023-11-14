package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgMigrate = "migrate"

var _ sdk.Msg = &MsgMigrate{}

func NewMsgMigrate(
	creator string,
	txHash string,
	ethAddress string,
	destAddress string,
	amount string,
	token uint64,
	logIndex uint64,
) *MsgMigrate {
	return &MsgMigrate{
		Creator:     creator,
		TxHash:      txHash,
		EthAddress:  ethAddress,
		DestAddress: destAddress,
		Amount:      amount,
		Token:       token,
		LogIndex:    logIndex,
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
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.DestAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid destination address (%s)", err)
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
