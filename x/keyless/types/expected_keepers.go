package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	identitytypes "selfchain/x/identity/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// IdentityKeeper defines the expected interface for the identity module
type IdentityKeeper interface {
	// DID Operations
	GetDIDDocument(ctx sdk.Context, did string) (identitytypes.DIDDocument, bool)
	VerifyDIDOwnership(ctx sdk.Context, did string, owner sdk.AccAddress) error

	// OAuth2 and MFA
	VerifyOAuth2Token(ctx sdk.Context, did string, token string) error
	VerifyMFA(ctx sdk.Context, did string) error

	// Recovery Operations
	VerifyRecoveryToken(ctx sdk.Context, did string, token string) error
	GetKeyShare(ctx sdk.Context, did string) ([]byte, bool)
	ReconstructWallet(ctx sdk.Context, didDoc identitytypes.DIDDocument) (interface{}, error)

	// Security Features
	CheckRateLimit(ctx sdk.Context, did string, operation string) error
	LogAuditEvent(ctx sdk.Context, event *identitytypes.AuditEvent) error
}
