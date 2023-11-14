package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/selfvesting module sentinel errors
var (
	ErrNoVestingPositions       = sdkerrors.Register(ModuleName, 1100, "Current account has not vesting positions")
	ErrPositionIndexOutOfBounds = sdkerrors.Register(ModuleName, 1101, "Position index out of bounds")
	ErrPositionFullyClaimed     = sdkerrors.Register(ModuleName, 1102, "Tokens fully claimed")
	ErrCliffViolation           = sdkerrors.Register(ModuleName, 1103, "Cliff period violation")
)
