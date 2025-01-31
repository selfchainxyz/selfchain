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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cometbft/cometbft-db"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/testutil/mocks"
	"selfchain/x/keyless/types"
	identitytypes "selfchain/x/identity/types"
)

type TestKeeper struct {
	*keeper.Keeper
	identityKeeper *mocks.IdentityKeeper
	tssProtocol   *mocks.TSSProtocol
}

// NewKeeper creates a new keeper for testing
func NewKeeper(t testing.TB) (*TestKeeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
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

	// Mock the identity keeper
	identityKeeper := mocks.NewIdentityKeeper()
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
	}, nil)
	
	identityKeeper.On("VerifyDIDOwnership", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("VerifyRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("GetKeyShare", mock.Anything, mock.Anything).Return([]byte("test-key-share"), nil)
	identityKeeper.On("ReconstructWallet", mock.Anything, mock.Anything).Return([]byte("test-reconstructed-wallet"), nil)
	identityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("GenerateRecoveryToken", mock.Anything, mock.Anything).Return("test-recovery-token", nil)
	identityKeeper.On("ValidateRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Mock the TSS protocol
	tssProtocol := mocks.NewTSSProtocol()
	
	// Set up mock expectations for TSS protocol
	tssProtocol.On("GenerateKeyShares", mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.MatchedBy(func(req *types.KeyGenRequest) bool { return true })).Return(&types.KeyGenResponse{
		WalletAddress: "test-wallet",
		PublicKey:     []byte("test-public-key"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     now,
			LastRotated:   now,
			LastUsed:      now,
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
		},
	}, nil)

	tssProtocol.On("ProcessKeyGenRound", mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.Anything, mock.MatchedBy(func(partyData *types.PartyData) bool { return true })).Return(nil)

	tssProtocol.On("InitiateSigning", mock.MatchedBy(func(ctx context.Context) bool { return true }), mock.Anything, mock.Anything).Return(&types.SigningResponse{
		WalletAddress: "test-wallet",
		Signature:     []byte("test-signature"),
		Metadata: &types.SignatureMetadata{
			Timestamp: &now,
			ChainId:   "test-chain",
			SignType:  types.SignatureType_SIGNATURE_TYPE_ECDSA,
		},
	}, nil)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		identityKeeper,
		tssProtocol,
	)

	testKeeper := &TestKeeper{
		Keeper:         k,
		identityKeeper: identityKeeper,
		tssProtocol:   tssProtocol,
	}

	return testKeeper, ctx
}

// GetMockIdentityKeeper returns the mock identity keeper
func (k *TestKeeper) GetMockIdentityKeeper() *mocks.IdentityKeeper {
	return k.identityKeeper
}

// GetMockTSSProtocol returns the mock TSS protocol
func (k *TestKeeper) GetMockTSSProtocol() *mocks.TSSProtocol {
	return k.tssProtocol
}

// MockKeyGenResponse creates a mock key generation response
func MockKeyGenResponse(walletAddress string) *types.KeyGenResponse {
	now := time.Now().UTC()
	return &types.KeyGenResponse{
		WalletAddress: walletAddress,
		PublicKey:     []byte("test-public-key"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     now,
			LastRotated:   now,
			LastUsed:      now,
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
		},
	}
}
