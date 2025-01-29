package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgGrantPermission{}
var _ sdk.Msg = &MsgRevokePermission{}

// NewMsgGrantPermission creates a new MsgGrantPermission instance
func NewMsgGrantPermission(creator string, walletAddress string, grantee string, permissions []WalletPermission, expiresAt *time.Time) *MsgGrantPermission {
	return &MsgGrantPermission{
		Creator:      creator,
		WalletAddress: walletAddress,
		Grantee:      grantee,
		Permissions:  permissions,
		ExpiresAt:    expiresAt,
	}
}

// ValidateBasic runs stateless checks on the message
func (msg *MsgGrantPermission) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Grantee); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid grantee address (%s)", err)
	}

	if len(msg.Permissions) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "permissions cannot be empty")
	}

	for _, perm := range msg.Permissions {
		if perm == WalletPermission_WALLET_PERMISSION_UNSPECIFIED {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid permission")
		}
	}

	return nil
}

// GetSigners returns the expected signers for a MsgGrantPermission message
func (msg *MsgGrantPermission) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// NewMsgRevokePermission creates a new MsgRevokePermission instance
func NewMsgRevokePermission(creator string, walletAddress string, grantee string, permissions []WalletPermission) *MsgRevokePermission {
	return &MsgRevokePermission{
		Creator:      creator,
		WalletAddress: walletAddress,
		Grantee:      grantee,
		Permissions:  permissions,
	}
}

// ValidateBasic runs stateless checks on the message
func (msg *MsgRevokePermission) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Grantee); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid grantee address (%s)", err)
	}

	if len(msg.Permissions) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "permissions cannot be empty")
	}

	for _, perm := range msg.Permissions {
		if perm == WalletPermission_WALLET_PERMISSION_UNSPECIFIED {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid permission")
		}
	}

	return nil
}

// GetSigners returns the expected signers for a MsgRevokePermission message
func (msg *MsgRevokePermission) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
