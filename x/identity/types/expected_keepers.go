package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper used for simulations
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account keeper should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank keeper should be defined here
}

// KeylessKeeper defines the expected interface for the keyless module
type KeylessKeeper interface {
	// ReconstructWallet reconstructs a wallet from a DID document
	ReconstructWallet(ctx sdk.Context, didDoc DIDDocument) ([]byte, error)

	// StoreKeyShare stores a key share for a DID
	StoreKeyShare(ctx sdk.Context, did string, keyShare []byte) error

	// GetKeyShare retrieves a key share for a DID
	GetKeyShare(ctx sdk.Context, did string) ([]byte, bool)

	// InitiateRecovery initiates the wallet recovery process
	InitiateRecovery(ctx sdk.Context, did string, recoveryToken string, recoveryAddress string) error
}
