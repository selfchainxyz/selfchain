package types

import (
	selfvestingTypes "selfchain/x/selfvesting/types"

	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	MintCoins(ctx context.Context, moduleName string, amounts sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// SelfvestingKeeper defines the expected interface needed to interact with the selfvesting module
type SelfvestingKeeper interface {
	AddBeneficiary(ctx sdk.Context, req selfvestingTypes.AddBeneficiaryRequest) (*selfvestingTypes.VestingInfo, error)
}
