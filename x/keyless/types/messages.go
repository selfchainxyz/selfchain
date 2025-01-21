package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgCreateWallet{}
	_ sdk.Msg = &MsgRecoverWallet{}
	_ sdk.Msg = &MsgSignTransaction{}
	_ sdk.Msg = &MsgBatchSign{}
	_ sdk.Msg = &MsgInitiateKeyRotation{}
	_ sdk.Msg = &MsgCompleteKeyRotation{}
	_ sdk.Msg = &MsgCancelKeyRotation{}
)

// Message type constants
const (
	TypeMsgCreateWallet          = "create_wallet"
	TypeMsgRecoverWallet         = "recover_wallet"
	TypeMsgSignTransaction       = "sign_transaction"
	TypeMsgBatchSign            = "batch_sign"
	TypeMsgInitiateKeyRotation  = "initiate_key_rotation"
	TypeMsgCompleteKeyRotation  = "complete_key_rotation"
	TypeMsgCancelKeyRotation    = "cancel_key_rotation"
)

// GetSigners returns the expected signers for MsgCreateWallet
func (msg *MsgCreateWallet) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ValidateBasic performs stateless validation on MsgCreateWallet
func (msg *MsgCreateWallet) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.PubKey == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "public key cannot be empty")
	}
	if msg.WalletAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wallet address cannot be empty")
	}
	if msg.ChainId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "chain ID cannot be empty")
	}
	return nil
}

// GetSigners returns the expected signers for MsgRecoverWallet
func (msg *MsgRecoverWallet) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ValidateBasic performs stateless validation on MsgRecoverWallet
func (msg *MsgRecoverWallet) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.WalletAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wallet address cannot be empty")
	}
	if msg.RecoveryProof == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "recovery proof cannot be empty")
	}
	if msg.NewPubKey == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "new public key cannot be empty")
	}
	if msg.Signature == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "signature cannot be empty")
	}
	return nil
}

// GetSigners returns the expected signers for MsgSignTransaction
func (msg *MsgSignTransaction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ValidateBasic performs stateless validation on MsgSignTransaction
func (msg *MsgSignTransaction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.WalletAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wallet address cannot be empty")
	}
	if msg.UnsignedTx == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unsigned transaction cannot be empty")
	}
	if msg.ChainId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "chain ID cannot be empty")
	}
	return nil
}

// GetSigners returns the expected signers for MsgBatchSign
func (msg *MsgBatchSign) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ValidateBasic performs stateless validation on MsgBatchSign
func (msg *MsgBatchSign) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.WalletAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wallet address cannot be empty")
	}
	if len(msg.UnsignedTxs) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unsigned transactions cannot be empty")
	}
	if msg.ChainId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "chain ID cannot be empty")
	}
	return nil
}

// GetSigners returns the expected signers for MsgInitiateKeyRotation
func (msg *MsgInitiateKeyRotation) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ValidateBasic performs stateless validation on MsgInitiateKeyRotation
func (msg *MsgInitiateKeyRotation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.WalletAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wallet address cannot be empty")
	}
	if msg.NewPubKey == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "new public key cannot be empty")
	}
	return nil
}

// GetSigners returns the expected signers for MsgCompleteKeyRotation
func (msg *MsgCompleteKeyRotation) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ValidateBasic performs stateless validation on MsgCompleteKeyRotation
func (msg *MsgCompleteKeyRotation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.WalletAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wallet address cannot be empty")
	}
	if msg.Signature == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "signature cannot be empty")
	}
	return nil
}

// GetSigners returns the expected signers for MsgCancelKeyRotation
func (msg *MsgCancelKeyRotation) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ValidateBasic performs stateless validation on MsgCancelKeyRotation
func (msg *MsgCancelKeyRotation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.WalletAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wallet address cannot be empty")
	}
	return nil
}
