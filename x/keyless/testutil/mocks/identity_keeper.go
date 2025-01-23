package mocks

import (
	"selfchain/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
)

// IdentityKeeper is a mock implementation of the identity keeper
type IdentityKeeper struct {
	mock.Mock
}

func NewIdentityKeeper(t mock.TestingT) *IdentityKeeper {
	mock := &IdentityKeeper{}
	mock.Test(t)
	return mock
}

func (m *IdentityKeeper) GetDIDDocument(ctx sdk.Context, did string) (*types.DIDDocument, error) {
	args := m.Called(ctx, did)
	if doc := args.Get(0); doc != nil {
		return doc.(*types.DIDDocument), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *IdentityKeeper) VerifyDIDOwnership(ctx sdk.Context, did string, owner sdk.AccAddress) error {
	args := m.Called(ctx, did, owner)
	return args.Error(0)
}

func (m *IdentityKeeper) VerifyOAuth2Token(ctx sdk.Context, did string, token string) error {
	args := m.Called(ctx, did, token)
	return args.Error(0)
}

func (m *IdentityKeeper) VerifyMFA(ctx sdk.Context, did string) error {
	args := m.Called(ctx, did)
	return args.Error(0)
}

func (m *IdentityKeeper) VerifyRecoveryToken(ctx sdk.Context, did string, token string) error {
	args := m.Called(ctx, did, token)
	return args.Error(0)
}

func (m *IdentityKeeper) GetKeyShare(ctx sdk.Context, did string) ([]byte, bool) {
	args := m.Called(ctx, did)
	return args.Get(0).([]byte), args.Bool(1)
}

func (m *IdentityKeeper) ReconstructWallet(ctx sdk.Context, didDoc types.DIDDocument) (interface{}, error) {
	args := m.Called(ctx, didDoc)
	return args.Get(0), args.Error(1)
}

func (m *IdentityKeeper) CheckRateLimit(ctx sdk.Context, did string, operation string) error {
	args := m.Called(ctx, did, operation)
	return args.Error(0)
}

func (m *IdentityKeeper) LogAuditEvent(ctx sdk.Context, event *types.AuditEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// Helper methods for tests

func (m *IdentityKeeper) StoreDIDDocument(ctx sdk.Context, doc *types.DIDDocument) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *IdentityKeeper) StoreOAuth2Token(ctx sdk.Context, did string, token string) error {
	args := m.Called(ctx, did, token)
	return args.Error(0)
}

func (m *IdentityKeeper) StoreMFASession(ctx sdk.Context, did string) error {
	args := m.Called(ctx, did)
	return args.Error(0)
}

func (m *IdentityKeeper) StoreRecoveryToken(ctx sdk.Context, did string, token string) error {
	args := m.Called(ctx, did, token)
	return args.Error(0)
}

func (m *IdentityKeeper) StoreKeyShare(ctx sdk.Context, did string, share []byte) error {
	args := m.Called(ctx, did, share)
	return args.Error(0)
}
