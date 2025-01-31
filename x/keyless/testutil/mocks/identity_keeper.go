package mocks

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	identitytypes "selfchain/x/identity/types"
)

type IdentityKeeper struct {
	mock.Mock
}

// NewIdentityKeeper creates a new mock IdentityKeeper
func NewIdentityKeeper() *IdentityKeeper {
	return &IdentityKeeper{}
}

func (m *IdentityKeeper) GetDIDDocument(ctx sdk.Context, did string) (identitytypes.DIDDocument, bool) {
	args := m.Called(ctx, did)
	return args.Get(0).(identitytypes.DIDDocument), args.Bool(1)
}

func (m *IdentityKeeper) VerifyDIDOwnership(ctx sdk.Context, did string, owner sdk.AccAddress) error {
	args := m.Called(ctx, did, owner)
	return args.Error(0)
}

func (m *IdentityKeeper) VerifyRecoveryProof(ctx sdk.Context, did string, proof []byte) error {
	args := m.Called(ctx, did, proof)
	return args.Error(0)
}

func (m *IdentityKeeper) ValidateIdentityStatus(ctx sdk.Context, did string) error {
	args := m.Called(ctx, did)
	return args.Error(0)
}

func (m *IdentityKeeper) CheckRateLimit(ctx sdk.Context, did string, operation string) error {
	args := m.Called(ctx, did, operation)
	return args.Error(0)
}

func (m *IdentityKeeper) LogAuditEvent(ctx sdk.Context, event *identitytypes.AuditEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *IdentityKeeper) GenerateRecoveryToken(ctx sdk.Context, walletID string) (string, error) {
	args := m.Called(ctx, walletID)
	return args.String(0), args.Error(1)
}

func (m *IdentityKeeper) ValidateRecoveryToken(ctx sdk.Context, walletID string, token string) error {
	args := m.Called(ctx, walletID, token)
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

func (m *IdentityKeeper) VerifyMFA(ctx sdk.Context, did string) error {
	args := m.Called(ctx, did)
	return args.Error(0)
}

func (m *IdentityKeeper) VerifyOAuth2Token(ctx sdk.Context, did string, token string) error {
	args := m.Called(ctx, did, token)
	return args.Error(0)
}

func (m *IdentityKeeper) ReconstructWallet(ctx sdk.Context, didDoc identitytypes.DIDDocument) (interface{}, error) {
	args := m.Called(ctx, didDoc)
	return args.Get(0), args.Error(1)
}
