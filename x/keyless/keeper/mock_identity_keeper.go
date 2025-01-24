package keeper

import (
	identitytypes "selfchain/x/identity/types"
	"selfchain/x/keyless/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MockIdentityKeeper is a mock implementation of the identity keeper for testing
type MockIdentityKeeper struct {
	didDocs map[string]*identitytypes.DIDDocument
	shares  map[string][]byte
}

// NewMockIdentityKeeper creates a new mock identity keeper
func NewMockIdentityKeeper() *MockIdentityKeeper {
	return &MockIdentityKeeper{
		didDocs: make(map[string]*identitytypes.DIDDocument),
		shares:  make(map[string][]byte),
	}
}

// SetDIDDocument sets a DID document in the mock keeper
func (k *MockIdentityKeeper) SetDIDDocument(ctx sdk.Context, did string, doc *identitytypes.DIDDocument) {
	k.didDocs[did] = doc
}

// GetDIDDocument gets a DID document from the mock keeper
func (k *MockIdentityKeeper) GetDIDDocument(ctx sdk.Context, did string) (identitytypes.DIDDocument, bool) {
	doc, found := k.didDocs[did]
	if !found {
		return identitytypes.DIDDocument{}, false
	}
	return *doc, found
}

// VerifyDIDOwnership verifies DID ownership for testing
func (k *MockIdentityKeeper) VerifyDIDOwnership(ctx sdk.Context, did string, owner sdk.AccAddress) error {
	// For testing, always return nil to indicate valid ownership
	return nil
}

// VerifyOAuth2Token verifies OAuth2 token for testing
func (k *MockIdentityKeeper) VerifyOAuth2Token(ctx sdk.Context, did string, token string) error {
	// For testing, always return nil to indicate valid token
	return nil
}

// VerifyMFA verifies MFA for testing
func (k *MockIdentityKeeper) VerifyMFA(ctx sdk.Context, did string) error {
	// For testing, always return nil to indicate valid MFA
	return nil
}

// VerifyRecoveryToken verifies recovery token for testing
func (k *MockIdentityKeeper) VerifyRecoveryToken(ctx sdk.Context, did string, token string) error {
	// For testing, always return nil to indicate valid recovery token
	return nil
}

// GetKeyShare gets a key share for testing
func (k *MockIdentityKeeper) GetKeyShare(ctx sdk.Context, did string) ([]byte, bool) {
	share, found := k.shares[did]
	return share, found
}

// SetKeyShare sets a key share for testing
func (k *MockIdentityKeeper) SetKeyShare(ctx sdk.Context, did string, share []byte) {
	k.shares[did] = share
}

// ReconstructWallet reconstructs a wallet for testing
func (k *MockIdentityKeeper) ReconstructWallet(ctx sdk.Context, didDoc identitytypes.DIDDocument) (interface{}, error) {
	// For testing, return nil values
	return nil, nil
}

// CheckRateLimit implements rate limiting check for testing
func (k *MockIdentityKeeper) CheckRateLimit(ctx sdk.Context, did string, operation string) error {
	// For testing, always return nil to indicate no rate limiting
	return nil
}

// LogAuditEvent logs an audit event for testing
func (k *MockIdentityKeeper) LogAuditEvent(ctx sdk.Context, event *identitytypes.AuditEvent) error {
	// For testing, always return nil to indicate successful logging
	return nil
}

// GetIdentityKeeper returns the identity keeper
func (k *MockIdentityKeeper) GetIdentityKeeper() types.IdentityKeeper {
	return k
}
