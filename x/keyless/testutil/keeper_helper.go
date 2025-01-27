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
	"github.com/stretchr/testify/mock"

	identitytypes "selfchain/x/identity/types"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/testutil/mocks"
	"selfchain/x/keyless/types"
)

// MockTSSProtocol is a mock implementation of the TSS protocol for testing
type MockTSSProtocol struct{}

// NewMockTSSProtocol creates a new instance of MockTSSProtocol
func NewMockTSSProtocol() *MockTSSProtocol {
	return &MockTSSProtocol{}
}

// GenerateKeyShares returns mock key shares
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

// ProcessKeyGenRound returns nil for mock key generation round processing
func (m *MockTSSProtocol) ProcessKeyGenRound(ctx context.Context, sessionID string, partyData *types.PartyData) error {
	return nil
}

// InitiateSigning returns a mock signing response
func (m *MockTSSProtocol) InitiateSigning(ctx context.Context, msg []byte, walletID string) (*types.SigningResponse, error) {
	now := time.Now().UTC()
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

	// Create mock identity keeper
	identityKeeper := mocks.NewIdentityKeeper(t)
	
	// Set up mock expectations
	identityKeeper.On("GetDIDDocument", mock.Anything, mock.Anything).Return(identitytypes.DIDDocument{}, true)
	identityKeeper.On("VerifyDIDOwnership", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("VerifyRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("GetKeyShare", mock.Anything, mock.Anything).Return([]byte("mock_key_share"), true)
	identityKeeper.On("ReconstructWallet", mock.Anything, mock.Anything).Return(nil, nil)
	identityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("GenerateRecoveryToken", mock.Anything, mock.Anything).Return("mock_recovery_token", nil)
	identityKeeper.On("ValidateRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockTSS := NewMockTSSProtocol()

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		identityKeeper,
		mockTSS,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{
		ChainID: "test-chain",
		Height:  1,
		Time:    time.Now().UTC(),
	}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
