package keeper

import (
	"fmt"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	identitytypes "selfchain/x/identity/types"
)

// MockIdentityKeeper is a mock implementation of the identity keeper interface
type MockIdentityKeeper struct {
	didDocuments sync.Map
	keyShares    sync.Map
}

// NewMockIdentityKeeper creates a new mock identity keeper
func NewMockIdentityKeeper() *MockIdentityKeeper {
	return &MockIdentityKeeper{}
}

// GetDIDDocument returns a DID document
func (k *MockIdentityKeeper) GetDIDDocument(ctx sdk.Context, did string) (identitytypes.DIDDocument, bool) {
	if doc, ok := k.didDocuments.Load(did); ok {
		return *doc.(*identitytypes.DIDDocument), true
	}
	return identitytypes.DIDDocument{}, false
}

// SetDIDDocument sets a DID document
func (k *MockIdentityKeeper) SetDIDDocument(ctx sdk.Context, did string, doc *identitytypes.DIDDocument) {
	k.didDocuments.Store(did, doc)
}

// GetKeyShare returns a key share
func (k *MockIdentityKeeper) GetKeyShare(ctx sdk.Context, did string) ([]byte, bool) {
	if share, ok := k.keyShares.Load(did); ok {
		return share.([]byte), true
	}
	return nil, false
}

// SetKeyShare sets a key share
func (k *MockIdentityKeeper) SetKeyShare(ctx sdk.Context, did string, share []byte) {
	k.keyShares.Store(did, share)
}

// VerifyRecoveryToken verifies a recovery token
func (k *MockIdentityKeeper) VerifyRecoveryToken(ctx sdk.Context, did, token string) error {
	// Mock implementation always succeeds
	return nil
}

// ReconstructWallet reconstructs a wallet
func (k *MockIdentityKeeper) ReconstructWallet(ctx sdk.Context, doc identitytypes.DIDDocument) (interface{}, error) {
	// Mock implementation returns a simple wallet with the first verification method's public key
	if len(doc.VerificationMethod) == 0 {
		return nil, fmt.Errorf("no verification methods found")
	}
	return &struct {
		PubKey string
	}{
		PubKey: doc.VerificationMethod[0].PublicKeyBase58,
	}, nil
}

// VerifyDIDOwnership verifies DID ownership
func (k *MockIdentityKeeper) VerifyDIDOwnership(ctx sdk.Context, did string, owner sdk.AccAddress) error {
	// Mock implementation always succeeds
	return nil
}

// VerifyOAuth2Token verifies an OAuth2 token
func (k *MockIdentityKeeper) VerifyOAuth2Token(ctx sdk.Context, did string, token string) error {
	// Mock implementation always succeeds
	return nil
}

// VerifyMFA verifies MFA
func (k *MockIdentityKeeper) VerifyMFA(ctx sdk.Context, did string) error {
	// Mock implementation always succeeds
	return nil
}

// CheckRateLimit checks rate limit for an operation
func (k *MockIdentityKeeper) CheckRateLimit(ctx sdk.Context, did string, operation string) error {
	// Mock implementation always succeeds
	return nil
}

// LogAuditEvent logs an audit event
func (k *MockIdentityKeeper) LogAuditEvent(ctx sdk.Context, event *identitytypes.AuditEvent) error {
	// Mock implementation always succeeds
	return nil
}
