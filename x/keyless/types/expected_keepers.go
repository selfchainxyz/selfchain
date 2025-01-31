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
	GetDIDDocument(ctx sdk.Context, did string) (identitytypes.DIDDocument, bool)
	SaveDIDDocument(ctx sdk.Context, doc identitytypes.DIDDocument) error
	ReconstructWalletFromDID(ctx sdk.Context, doc identitytypes.DIDDocument) ([]byte, error)
	
	// Authentication and Authorization
	VerifyOAuth2Token(ctx sdk.Context, token string, scope string) error
	VerifyMFA(ctx sdk.Context, code string) error
	CheckRateLimit(ctx sdk.Context, did string, action string) error
	
	// Audit and Logging
	LogAuditEvent(ctx sdk.Context, event *identitytypes.AuditEvent) error
}

// TSSProtocol defines the interface for TSS protocol operations
type TSSProtocol interface {
	// GenerateKeyShares generates key shares for a new wallet
	GenerateKeyShares(ctx sdk.Context, walletAddress string, threshold uint32, securityLevel SecurityLevel) (*KeyGenResponse, error)

	// ReconstructKey reconstructs a key from shares
	ReconstructKey(ctx sdk.Context, shares [][]byte) ([]byte, error)

	// VerifyShare verifies a share's validity
	VerifyShare(ctx sdk.Context, share []byte, publicKey []byte) error

	// SignMessage signs a message using TSS
	SignMessage(ctx sdk.Context, message []byte, shares [][]byte) ([]byte, error)

	// VerifySignature verifies a TSS signature
	VerifySignature(ctx sdk.Context, message []byte, signature []byte, publicKey []byte) error

	// GetPartyData gets TSS party data
	GetPartyData(ctx sdk.Context, partyID string) (*PartyData, error)

	// SetPartyData sets TSS party data
	SetPartyData(ctx sdk.Context, data *PartyData) error
}

// KeylessKeeper defines the expected interface for other modules
type KeylessKeeper interface {
	// ReconstructWallet reconstructs a wallet from a DID document
	ReconstructWallet(ctx sdk.Context, didDoc identitytypes.DIDDocument) ([]byte, error)

	// StoreKeyShare stores a key share for a DID
	StoreKeyShare(ctx sdk.Context, did string, keyShare []byte) error

	// GetKeyShare retrieves a key share for a DID
	GetKeyShare(ctx sdk.Context, did string) ([]byte, bool)

	// InitiateRecovery initiates the wallet recovery process
	InitiateRecovery(ctx sdk.Context, did string, recoveryToken string, recoveryAddress string) error
}
