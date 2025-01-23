package mock

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	identitytypes "selfchain/x/identity/types"
)

// MockIdentityKeeper is a mock implementation of the identity keeper for testing
type MockIdentityKeeper struct {
	didDocuments map[string]identitytypes.DIDDocument
	keyShares    map[string][]byte
}

// NewMockIdentityKeeper creates a new mock identity keeper
func NewMockIdentityKeeper() *MockIdentityKeeper {
	return &MockIdentityKeeper{
		didDocuments: make(map[string]identitytypes.DIDDocument),
		keyShares:    make(map[string][]byte),
	}
}

// GetDIDDocument implements IdentityKeeper
func (m *MockIdentityKeeper) GetDIDDocument(ctx sdk.Context, did string) (identitytypes.DIDDocument, bool) {
	doc, found := m.didDocuments[did]
	return doc, found
}

// VerifyDIDOwnership implements IdentityKeeper
func (m *MockIdentityKeeper) VerifyDIDOwnership(ctx sdk.Context, did string, owner sdk.AccAddress) error {
	doc, found := m.GetDIDDocument(ctx, did)
	if !found {
		return fmt.Errorf("DID document not found: %s", did)
	}
	if !doc.HasController(owner.String()) {
		return fmt.Errorf("address %s is not the controller of DID %s", owner.String(), did)
	}
	return nil
}

// VerifyOAuth2Token implements IdentityKeeper
func (m *MockIdentityKeeper) VerifyOAuth2Token(ctx sdk.Context, did string, token string) error {
	// For testing, accept any token
	return nil
}

// VerifyMFA implements IdentityKeeper
func (m *MockIdentityKeeper) VerifyMFA(ctx sdk.Context, did string) error {
	// For testing, always pass MFA
	return nil
}

// VerifyRecoveryToken implements IdentityKeeper
func (m *MockIdentityKeeper) VerifyRecoveryToken(ctx sdk.Context, did string, token string) error {
	// For testing, accept any token
	return nil
}

// GetKeyShare implements IdentityKeeper
func (m *MockIdentityKeeper) GetKeyShare(ctx sdk.Context, did string) ([]byte, bool) {
	share, found := m.keyShares[did]
	return share, found
}

// ReconstructWallet implements IdentityKeeper
func (m *MockIdentityKeeper) ReconstructWallet(ctx sdk.Context, didDoc identitytypes.DIDDocument) (interface{}, error) {
	// For testing, return a dummy wallet
	return []byte("reconstructed_wallet"), nil
}

// CheckRateLimit implements IdentityKeeper
func (m *MockIdentityKeeper) CheckRateLimit(ctx sdk.Context, did string, operation string) error {
	// For testing, no rate limits
	return nil
}

// LogAuditEvent implements IdentityKeeper
func (m *MockIdentityKeeper) LogAuditEvent(ctx sdk.Context, event *identitytypes.AuditEvent) error {
	// For testing, just return nil
	return nil
}

// Helper methods for testing

// SetDIDDocument adds a DID document for testing
func (m *MockIdentityKeeper) SetDIDDocument(did string, doc identitytypes.DIDDocument) {
	m.didDocuments[did] = doc
}

// SetKeyShare adds a key share for testing
func (m *MockIdentityKeeper) SetKeyShare(did string, share []byte) {
	m.keyShares[did] = share
}
