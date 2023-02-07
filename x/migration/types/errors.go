package types

// DONTCOVER

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/migration module sentinel errors
var (
	ErrInvalidMigrationAmount = sdkerrors.Register(ModuleName, 1100, fmt.Sprintf("Migration amount should be at least %s", GetMinMigrationAmount().String()))
	ErrTokenNotSupported      = sdkerrors.Register(ModuleName, 1101, "Token should either Front or Hotcross")
	ErrEmptyStringValue       = sdkerrors.Register(ModuleName, 1102, "TxHash or EthAddress have empty string value")
	ErrUnknownMigrator        = sdkerrors.Register(ModuleName, 1103, "Unknown migrator")
	ErrMigrationProcessed     = sdkerrors.Register(ModuleName, 1104, "The given migration message has been previously processed")
	ErrOnlyAdmin              = sdkerrors.Register(ModuleName, 1105, "Only admin can update migrators")
)
