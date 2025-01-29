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
	mockIdentityKeeper := mocks.NewIdentityKeeper(t)
	mockIdentityKeeper.On("GetIdentity", mock.Anything, mock.Anything).Return(&identitytypes.Identity{
		Id:     "test-id",
		Status: identitytypes.IdentityStatus_IDENTITY_STATUS_ACTIVE,
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
	mockTSSProtocol := mocks.NewTSSProtocol(t)
	mockTSSProtocol.On("ValidateRecoveryProof", mock.Anything, mock.Anything).Return(nil)
	mockTSSProtocol.On("GenerateKeygenSession", mock.Anything, mock.Anything).Return(&types.KeygenSession{
		Id:            "test-session",
		WalletAddress: "test-wallet",
		Status:        types.SessionStatus_SESSION_STATUS_ACTIVE,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
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
	store := k.stateStore.GetKVStore(k.storeKey)
	iter := store.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// MockSigningResponse creates a mock signing response
func MockSigningResponse(walletAddress string) *types.SigningResponse {
	now := time.Now().UTC()
	return &types.SigningResponse{
		WalletAddress: walletAddress,
		Signature:     []byte("test-signature"),
		Status:        types.SigningStatus_SIGNING_STATUS_COMPLETED,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
}
