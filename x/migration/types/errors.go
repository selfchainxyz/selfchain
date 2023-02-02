package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/migration module sentinel errors
var (
	ErrInvalidMigrationAmount = sdkerrors.Register(ModuleName, 1100, "Migration amount must be be a positive integer")
)
