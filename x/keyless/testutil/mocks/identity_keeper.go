package mocks

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	identitytypes "selfchain/x/identity/types"
)

type IdentityKeeper struct {
	mock.Mock
}

// NewIdentityKeeper creates a new mock IdentityKeeper
func NewIdentityKeeper() *IdentityKeeper {
	m := &IdentityKeeper{}
	now := time.Now().UTC()

	// Set up default mock expectations
	m.On("GetDIDDocument", mock.Anything, mock.Anything).Return(identitytypes.DIDDocument{
		Id:         "test-did",
		Controller: []string{"test-controller"},
		Created:    &now,
		Updated:    &now,
		Status:     identitytypes.Status_STATUS_ACTIVE,
	}, true)

	m.On("SaveDIDDocument", mock.Anything, mock.Anything).Return(nil)
	m.On("ReconstructWalletFromDID", mock.Anything, mock.Anything).Return([]byte("test-wallet"), nil)
	m.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil)
	m.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil)

	return m
}

func (m *IdentityKeeper) GetDIDDocument(ctx sdk.Context, did string) (identitytypes.DIDDocument, bool) {
	args := m.Called(ctx, did)
	return args.Get(0).(identitytypes.DIDDocument), args.Bool(1)
}

func (m *IdentityKeeper) SaveDIDDocument(ctx sdk.Context, doc identitytypes.DIDDocument) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *IdentityKeeper) ReconstructWalletFromDID(ctx sdk.Context, doc identitytypes.DIDDocument) ([]byte, error) {
	args := m.Called(ctx, doc)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *IdentityKeeper) VerifyOAuth2Token(ctx sdk.Context, token string, scope string) error {
	args := m.Called(ctx, token, scope)
	return args.Error(0)
}

func (m *IdentityKeeper) VerifyMFA(ctx sdk.Context, code string) error {
	args := m.Called(ctx, code)
	return args.Error(0)
}

func (m *IdentityKeeper) CheckRateLimit(ctx sdk.Context, did string, action string) error {
	args := m.Called(ctx, did, action)
	return args.Error(0)
}

func (m *IdentityKeeper) LogAuditEvent(ctx sdk.Context, event *identitytypes.AuditEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
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

func (m *IdentityKeeper) ReconstructWallet(ctx sdk.Context, didDoc identitytypes.DIDDocument) (interface{}, error) {
	args := m.Called(ctx, didDoc)
	return args.Get(0), args.Error(1)
}
