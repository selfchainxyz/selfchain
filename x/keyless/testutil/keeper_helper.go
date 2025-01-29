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
		WalletAddress: req.WalletAddress,
		PublicKey:    []byte("mock_public_key"),
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
		WalletAddress: walletID,
		Signature:    []byte("mock_signature"),
		Metadata: &types.SignatureMetadata{
			Timestamp: &now,
			ChainId:   "test-chain",
			SignType:  types.SignatureType_SIGNATURE_TYPE_ECDSA,
		},
	}, nil
}

// MockKeyGenRequest creates a mock key generation request
func MockKeyGenRequest(walletAddress string) *types.KeyGenRequest {
	return &types.KeyGenRequest{
		WalletAddress:  walletAddress,
		ChainId:       "test-chain",
		SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_STANDARD,
	}
}

// MockSigningResponse creates a mock signing response
func MockSigningResponse(walletAddress string) *types.SigningResponse {
	now := time.Now().UTC()
	timestamp := now
	return &types.SigningResponse{
		WalletAddress: walletAddress,
		Signature:    []byte("mock-signature"),
		Metadata: &types.SignatureMetadata{
			Timestamp: &timestamp,
			ChainId:   "test-chain",
			SignType:  types.SignatureType_SIGNATURE_TYPE_ECDSA,
		},
	}
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
	now := time.Now().UTC()
	
	// Set up mock expectations for identity keeper
	identityKeeper.On("GetDIDDocument", mock.Anything, mock.Anything).Return(&identitytypes.DIDDocument{
		Id:         "test-did",
		Controller: []string{"test-controller"},
		VerificationMethod: []*identitytypes.VerificationMethod{
			{
				Id:              "test-verification-method",
				Type:            "Ed25519VerificationKey2020",
				Controller:      "test-controller",
				PublicKeyBase58: "test-public-key",
			},
		},
		Created: &now,
		Updated: &now,
		Status:  identitytypes.Status_STATUS_ACTIVE,
	}, nil).Maybe()
	
	// Basic verification mocks
	identityKeeper.On("VerifyDIDOwnership", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	identityKeeper.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	identityKeeper.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil).Maybe()
	identityKeeper.On("VerifyRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	
	// Key management mocks
	identityKeeper.On("GetKeyShare", mock.Anything, mock.Anything).Return([]byte("test-key-share"), nil).Maybe()
	identityKeeper.On("ReconstructWallet", mock.Anything, mock.Anything).Return([]byte("test-reconstructed-wallet"), nil).Maybe()
	identityKeeper.On("StoreKeyShare", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	identityKeeper.On("DeleteKeyShare", mock.Anything, mock.Anything).Return(nil).Maybe()
	
	// Security and audit mocks
	identityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	identityKeeper.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil).Maybe()
	identityKeeper.On("GenerateRecoveryToken", mock.Anything, mock.Anything).Return("test-recovery-token", nil).Maybe()
	identityKeeper.On("ValidateRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	// Create mock TSS protocol
	tssProtocol := mocks.NewTSSProtocol(t)
	
	// Key generation mocks
	tssProtocol.On("GenerateKey", mock.Anything, mock.Anything).Return(MockKeyGenRequest("test-wallet"), nil).Maybe()
	
	// Key generation round processing mocks
	tssProtocol.On("ProcessKeyGenRound", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	
	// Signing mocks
	tssProtocol.On("InitiateSigning", mock.Anything, mock.Anything, mock.Anything).Return(MockSigningResponse("test-wallet"), nil).Maybe()
	
	// Key rotation mocks
	tssProtocol.On("InitiateKeyRotation", mock.Anything, mock.Anything).Return(&types.KeyRotationResponse{
		WalletAddress: "test-wallet",
		NewPublicKey:  []byte("test-new-public-key"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     now,
			LastRotated:   now,
			LastUsed:      now,
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
		},
	}, nil).Maybe()
	
	tssProtocol.On("ProcessKeyRotationRound", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		identityKeeper,
		tssProtocol,
	)

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
