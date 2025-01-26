package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cometbft/cometbft-db"

	identitytypes "selfchain/x/identity/types"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
)

// MockIdentityKeeper is a mock implementation of the identity keeper for testing
type MockIdentityKeeper struct {
	didDocs map[string]*identitytypes.DIDDocument
}

// NewMockIdentityKeeper creates a new instance of MockIdentityKeeper
func NewMockIdentityKeeper() *MockIdentityKeeper {
	return &MockIdentityKeeper{
		didDocs: make(map[string]*identitytypes.DIDDocument),
	}
}

// GetDIDDocument returns a mock DID document
func (m *MockIdentityKeeper) GetDIDDocument(ctx sdk.Context, did string) (identitytypes.DIDDocument, bool) {
	if doc, ok := m.didDocs[did]; ok {
		return *doc, true
	}
	return identitytypes.DIDDocument{}, false
}

// VerifyDIDOwnership returns nil for mock DID ownership verification
func (m *MockIdentityKeeper) VerifyDIDOwnership(ctx sdk.Context, did string, owner sdk.AccAddress) error {
	return nil
}

// VerifyMFA returns nil for mock MFA verification
func (m *MockIdentityKeeper) VerifyMFA(ctx sdk.Context, did string) error {
	return nil
}

// VerifyRecoveryToken returns nil for mock recovery token verification
func (m *MockIdentityKeeper) VerifyRecoveryToken(ctx sdk.Context, did string, token string) error {
	return nil
}

// VerifyOAuth2Token returns nil for mock OAuth2 token verification
func (m *MockIdentityKeeper) VerifyOAuth2Token(ctx sdk.Context, did string, token string) error {
	return nil
}

// GetKeyShare returns a mock key share if DID document exists
func (m *MockIdentityKeeper) GetKeyShare(ctx sdk.Context, did string) ([]byte, bool) {
	if _, ok := m.didDocs[did]; ok {
		return []byte("mock_key_share"), true
	}
	return nil, false
}

// ReconstructWallet returns a mock reconstructed wallet
func (m *MockIdentityKeeper) ReconstructWallet(ctx sdk.Context, didDoc identitytypes.DIDDocument) (interface{}, error) {
	return nil, nil
}

// CheckRateLimit returns nil for mock rate limit check
func (m *MockIdentityKeeper) CheckRateLimit(ctx sdk.Context, did string, operation string) error {
	return nil
}

// LogAuditEvent returns nil for mock audit event logging
func (m *MockIdentityKeeper) LogAuditEvent(ctx sdk.Context, event *identitytypes.AuditEvent) error {
	return nil
}

// MockTSSProtocol is a mock implementation of the TSS protocol for testing
type MockTSSProtocol struct{}

func NewMockTSSProtocol() *MockTSSProtocol {
	return &MockTSSProtocol{}
}

func (m *MockTSSProtocol) GenerateKeyShares(ctx context.Context, req *types.KeyGenRequest) (*types.KeyGenResponse, error) {
	now := time.Now()
	return &types.KeyGenResponse{
		WalletId:  req.WalletId,
		PublicKey: []byte("mock_public_key"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     now,
			LastRotated:   now,
			LastUsed:      now,
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: req.SecurityLevel,
		},
	}, nil
}

func (m *MockTSSProtocol) ProcessKeyGenRound(ctx context.Context, sessionID string, partyData *types.PartyData) error {
	return nil
}

func (m *MockTSSProtocol) InitiateSigning(ctx context.Context, msg []byte, walletID string) (*types.SigningResponse, error) {
	now := time.Now()
	return &types.SigningResponse{
		WalletId:  walletID,
		Signature: []byte("mock_signature"),
		Metadata: &types.SignatureMetadata{
			Timestamp: &now,
			ChainId:   "test-chain",
			SignType:  types.SignatureType_SIGNATURE_TYPE_ECDSA,
		},
	}, nil
}

// NewTestKeeper creates a new keeper for testing
func NewTestKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Create a new memory database for each test
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := paramtypes.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"KeylessParams",
	)

	mockIdentityKeeper := NewMockIdentityKeeper()
	mockTSSProtocol := NewMockTSSProtocol()

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		mockIdentityKeeper,
		mockTSSProtocol,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
