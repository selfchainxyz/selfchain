package keeper

import (
	"testing"
	"time"

	identitytypes "selfchain/x/identity/types"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/testutil/mocks"
	"selfchain/x/keyless/types"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// KeylessKeeper is a keeper that contains all the necessary information to test the keyless module
type KeylessKeeper struct {
	*keeper.Keeper
	Ctx                sdk.Context
	MockIdentityKeeper *MockIdentityKeeper
	MockTSSProtocol    *mocks.TSSProtocol
	storeKey           storetypes.StoreKey
	memKey             storetypes.StoreKey
	db                 *dbm.MemDB
	stateStore         store.CommitMultiStore
}

// NewKeylessKeeper creates a new keyless keeper for testing
func NewKeylessKeeper(t testing.TB) *KeylessKeeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey("mem_keyless")

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	paramsSubspace := paramtypes.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"KeylessParams",
	)

	// Register the interfaces
	types.RegisterInterfaces(registry)
	identitytypes.RegisterInterfaces(registry)

	mockIdentityKeeper := NewMockIdentityKeeper()
	mockTSSProtocol := mocks.NewTSSProtocol()

	// Set up mock expectations for identity keeper
	mockIdentityKeeper.On("GetDIDDocument", mock.Anything, mock.Anything).Return(identitytypes.DIDDocument{
		Id: "test_did",
	}, true)
	mockIdentityKeeper.On("SaveDIDDocument", mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("ReconstructWalletFromDID", mock.Anything, mock.Anything).Return([]byte("test_wallet"), nil)
	mockIdentityKeeper.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil)

	// Set up mock expectations for TSS protocol
	mockTSSProtocol.On("GenerateKeyShares", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&types.KeyGenResponse{
		WalletAddress: "test_wallet",
		PublicKey:     []byte("test_pubkey"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     time.Now().UTC(),
			LastRotated:   time.Now().UTC(),
			LastUsed:      time.Now().UTC(),
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
		},
	}, nil)
	mockTSSProtocol.On("SignMessage", mock.Anything, mock.Anything, mock.Anything).Return([]byte("test_signature"), nil)
	mockTSSProtocol.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockTSSProtocol.On("VerifyShare", mock.Anything, mock.Anything, mock.Anything).Return(nil)

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

	return &KeylessKeeper{
		Keeper:             k,
		Ctx:               ctx,
		MockIdentityKeeper: mockIdentityKeeper,
		MockTSSProtocol:    mockTSSProtocol,
		storeKey:           storeKey,
		memKey:             memStoreKey,
		db:                 db,
		stateStore:         stateStore,
	}
}

// GetIdentityKeeper returns the mock identity keeper
func (k *KeylessKeeper) GetIdentityKeeper() *MockIdentityKeeper {
	return k.MockIdentityKeeper
}

// GetTSSProtocol returns the mock TSS protocol
func (k *KeylessKeeper) GetTSSProtocol() *mocks.TSSProtocol {
	return k.MockTSSProtocol
}

// ClearStore clears all data from the store
func (k *KeylessKeeper) ClearStore() {
	// Clear main store
	store := k.stateStore.GetKVStore(k.storeKey)
	iter := store.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}

	// Clear memory store
	memStore := k.stateStore.GetKVStore(k.memKey)
	memIter := memStore.Iterator(nil, nil)
	defer memIter.Close()

	for ; memIter.Valid(); memIter.Next() {
		memStore.Delete(memIter.Key())
	}

	// Reset mock identity keeper
	k.MockIdentityKeeper = NewMockIdentityKeeper()
	k.MockIdentityKeeper.On("GetDIDDocument", mock.Anything, mock.Anything).Return(identitytypes.DIDDocument{
		Id: "test_did",
	}, true)
	k.MockIdentityKeeper.On("SaveDIDDocument", mock.Anything, mock.Anything).Return(nil)
	k.MockIdentityKeeper.On("ReconstructWalletFromDID", mock.Anything, mock.Anything).Return([]byte("test_wallet"), nil)
	k.MockIdentityKeeper.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	k.MockIdentityKeeper.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil)
	k.MockIdentityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	k.MockIdentityKeeper.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil)

	// Reset mock TSS protocol
	k.MockTSSProtocol = mocks.NewTSSProtocol()
	k.MockTSSProtocol.On("GenerateKeyShares", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&types.KeyGenResponse{
		WalletAddress: "test_wallet",
		PublicKey:     []byte("test_pubkey"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     time.Now().UTC(),
			LastRotated:   time.Now().UTC(),
			LastUsed:      time.Now().UTC(),
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
		},
	}, nil)
	k.MockTSSProtocol.On("SignMessage", mock.Anything, mock.Anything, mock.Anything).Return([]byte("test_signature"), nil)
	k.MockTSSProtocol.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	k.MockTSSProtocol.On("VerifyShare", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Reset context with fresh store
	k.stateStore.Commit()
	k.Ctx = sdk.NewContext(k.stateStore, tmproto.Header{}, false, log.NewNopLogger())
}

// MockSigningResponse creates a mock signing response
func MockSigningResponse(walletAddress string) *types.SigningOutput {
	now := time.Now().UTC()
	return &types.SigningOutput{
		Signature:  []byte("test-signature"),
		PublicKey:  []byte("test-public-key"),
		SignedAt:   now.UTC(),
		KeyVersion: 1,
	}
}

// Mock implementation of IdentityKeeper
type MockIdentityKeeper struct {
	mock.Mock
}

func (m *MockIdentityKeeper) GetDIDDocument(ctx sdk.Context, did string) (identitytypes.DIDDocument, bool) {
	args := m.Called(ctx, did)
	return args.Get(0).(identitytypes.DIDDocument), args.Bool(1)
}

func (m *MockIdentityKeeper) SaveDIDDocument(ctx sdk.Context, doc identitytypes.DIDDocument) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockIdentityKeeper) ReconstructWalletFromDID(ctx sdk.Context, doc identitytypes.DIDDocument) ([]byte, error) {
	args := m.Called(ctx, doc)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockIdentityKeeper) VerifyOAuth2Token(ctx sdk.Context, token string, scope string) error {
	args := m.Called(ctx, token, scope)
	return args.Error(0)
}

func (m *MockIdentityKeeper) VerifyMFA(ctx sdk.Context, code string) error {
	args := m.Called(ctx, code)
	return args.Error(0)
}

func (m *MockIdentityKeeper) CheckRateLimit(ctx sdk.Context, did string, action string) error {
	args := m.Called(ctx, did, action)
	return args.Error(0)
}

func (m *MockIdentityKeeper) LogAuditEvent(ctx sdk.Context, event *identitytypes.AuditEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func NewMockIdentityKeeper() *MockIdentityKeeper {
	return &MockIdentityKeeper{}
}
