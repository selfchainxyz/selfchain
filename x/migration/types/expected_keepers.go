package types

import (
	selfvestingTypes "selfchain/x/selfvesting/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amounts sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// SelfvestingKeeper defines the expected interface needed to interact with the selfvesting module
type SelfvestingKeeper interface {
	AddBeneficiary(ctx sdk.Context, req selfvestingTypes.AddBeneficiaryRequest) (*selfvestingTypes.VestingInfo, error)
}
