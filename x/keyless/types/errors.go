package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/keyless module sentinel errors
var (
	ErrInvalidWalletAddress = sdkerrors.Register(ModuleName, 1100, "invalid wallet address")
	ErrInvalidDID          = sdkerrors.Register(ModuleName, 1101, "invalid DID")
	ErrWalletNotFound      = sdkerrors.Register(ModuleName, 1102, "wallet not found")
	ErrInvalidCreator      = sdkerrors.Register(ModuleName, 1103, "invalid creator")
	ErrInvalidSignature    = sdkerrors.Register(ModuleName, 1104, "invalid signature")
	ErrInvalidKeyShare     = sdkerrors.Register(ModuleName, 1105, "invalid key share")
	ErrWalletExists        = sdkerrors.Register(ModuleName, 1106, "wallet already exists")
	ErrInvalidStatus       = sdkerrors.Register(ModuleName, 1107, "invalid wallet status")
	ErrUnauthorized        = sdkerrors.Register(ModuleName, 1108, "unauthorized")
)
