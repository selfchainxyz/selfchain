package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/keyless module sentinel errors
var (
	ErrInvalidWalletAddress = sdkerrors.Register(ModuleName, 1100, "invalid wallet address")
	ErrInvalidWalletDID     = sdkerrors.Register(ModuleName, 1101, "invalid wallet DID")
	ErrInvalidWalletCreator = sdkerrors.Register(ModuleName, 1102, "invalid wallet creator")
	ErrWalletExists         = sdkerrors.Register(ModuleName, 1103, "wallet already exists")
	ErrWalletNotFound       = sdkerrors.Register(ModuleName, 1104, "wallet not found")
	ErrUnauthorized         = sdkerrors.Register(ModuleName, 1105, "unauthorized")
)
