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
	MockIdentityKeeper *mocks.IdentityKeeper
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

	// Create mock identity keeper
	mockIdentityKeeper := mocks.NewIdentityKeeper()
	now := time.Now().UTC()
	mockIdentityKeeper.On("GetDIDDocument", mock.Anything, mock.Anything).Return(&identitytypes.DIDDocument{
		Id:         "test-did",
		Controller: []string{"test-controller"},
		Created:    &now,
		Updated:    &now,
		Status:     identitytypes.Status_STATUS_ACTIVE,
	}, nil)
	mockIdentityKeeper.On("VerifyDIDOwnership", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("VerifyRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("GenerateRecoveryToken", mock.Anything, mock.Anything).Return("test-token", nil)
	mockIdentityKeeper.On("ValidateRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("GetKeyShare", mock.Anything, mock.Anything).Return([]byte("test-key-share"), true)
	mockIdentityKeeper.On("ReconstructWallet", mock.Anything, mock.Anything).Return(nil, nil)
	mockIdentityKeeper.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create mock TSS protocol
	mockTSSProtocol := mocks.NewTSSProtocol()
	mockTSSProtocol.On("ValidateRecoveryProof", mock.Anything, mock.Anything).Return(nil)
	mockTSSProtocol.On("GenerateKeyShares", mock.Anything, mock.Anything).Return(&types.KeyGenResponse{
		WalletAddress: "test-wallet",
		PublicKey:     []byte("test-pubkey"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     now.UTC(),
			LastRotated:   now.UTC(),
			LastUsed:      now.UTC(),
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
		},
	}, nil)

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
		Keeper:            k,
		Ctx:              ctx,
		MockIdentityKeeper: mockIdentityKeeper,
		MockTSSProtocol:    mockTSSProtocol,
		storeKey:         storeKey,
		memKey:           memStoreKey,
		db:               db,
		stateStore:       stateStore,
	}
}

// GetIdentityKeeper returns the mock identity keeper
func (k *KeylessKeeper) GetIdentityKeeper() *mocks.IdentityKeeper {
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
	k.MockIdentityKeeper = &mocks.IdentityKeeper{}
	k.MockIdentityKeeper.On("GetDIDDocument", mock.Anything, mock.Anything).Return(&identitytypes.DIDDocument{}, nil)
	k.MockIdentityKeeper.On("VerifyDIDOwnership", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	k.MockIdentityKeeper.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	k.MockIdentityKeeper.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil)
	k.MockIdentityKeeper.On("VerifyRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	k.MockIdentityKeeper.On("GetKeyShare", mock.Anything, mock.Anything).Return([]byte("mock_key_share"), true)
	k.MockIdentityKeeper.On("ReconstructWallet", mock.Anything, mock.Anything).Return(nil, nil)
	k.MockIdentityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	k.MockIdentityKeeper.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil)
	k.MockIdentityKeeper.On("GenerateRecoveryToken", mock.Anything, mock.Anything).Return("mock_token", nil)
	k.MockIdentityKeeper.On("ValidateRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Reset mock TSS protocol
	k.MockTSSProtocol = &mocks.TSSProtocol{}
	k.MockTSSProtocol.On("GenerateKeyShare", mock.Anything).Return([]byte("mock_key_share"), nil)
	k.MockTSSProtocol.On("ValidateRecoveryProof", mock.Anything, mock.Anything).Return(true, nil)
	k.MockTSSProtocol.On("Sign", mock.Anything, mock.Anything).Return([]byte("mock_signature"), nil)
	k.MockTSSProtocol.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	k.MockTSSProtocol.On("BatchSign", mock.Anything, mock.Anything).Return([]byte("mock_batch_signature"), nil)
	k.MockTSSProtocol.On("VerifyBatchSignature", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)

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
