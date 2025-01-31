package keeper

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
	"selfchain/x/keyless/testutil/mocks"
	identitytypes "selfchain/x/identity/types"
)

func KeylessKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"KeylessParams",
	)

	// Create mock keepers
	mockIdentityKeeper := mocks.NewIdentityKeeper()
	mockTSSProtocol := mocks.NewTSSProtocol()

	// Set up mock expectations
	now := time.Now().UTC()
	mockIdentityKeeper.On("GetDIDDocument", mock.Anything, mock.Anything).Return(&identitytypes.DIDDocument{
		Id:         "did:self:test",
		Controller: []string{"did:self:controller"},
		VerificationMethod: []*identitytypes.VerificationMethod{
			{
				Id:              "test_key",
				Type:            "Ed25519VerificationKey2020",
				PublicKeyBase58: "test_pubkey",
			},
		},
		Created:    &now,
		Updated:    &now,
		Status:     identitytypes.Status_STATUS_ACTIVE,
	}, true)

	mockIdentityKeeper.On("VerifyDIDOwnership", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("VerifyRecoveryProof", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("ValidateIdentityStatus", mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("GenerateRecoveryToken", mock.Anything, mock.Anything).Return("test_token", nil)
	mockIdentityKeeper.On("ValidateRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("GetKeyShare", mock.Anything, mock.Anything).Return([]byte("test_key_share"), true)
	mockIdentityKeeper.On("ReconstructWallet", mock.Anything, mock.Anything).Return(nil, nil)
	mockIdentityKeeper.On("VerifyRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mockTSSProtocol.On("GenerateKey", mock.Anything, mock.Anything, mock.Anything).Return([]byte("test_pubkey"), nil)
	mockTSSProtocol.On("Sign", mock.Anything, mock.Anything, mock.Anything).Return([]byte("test_signature"), nil)
	mockTSSProtocol.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockTSSProtocol.On("ValidateRecoveryProof", mock.Anything, mock.Anything).Return(nil)
	mockTSSProtocol.On("GenerateKeyShares", mock.Anything, mock.Anything).Return(&types.KeyGenResponse{
		WalletAddress: "test-wallet",
		PublicKey:     []byte("test-pubkey"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     now,
			LastRotated:   now,
			LastUsed:      now,
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
		},
	}, nil)
	mockTSSProtocol.On("InitiateSigning", mock.Anything, mock.Anything, mock.Anything).Return(&types.SigningResponse{
		WalletAddress: "test_wallet",
		Signature:     []byte("test_signature"),
		Metadata: &types.SignatureMetadata{
			Timestamp: &now,
			ChainId:   "test-chain",
			SignType:  types.SignatureType_SIGNATURE_TYPE_ECDSA,
		},
	}, nil)
	mockTSSProtocol.On("ProcessKeyGenRound", mock.Anything, mock.Anything, mock.AnythingOfType("*types.PartyData")).Return(nil)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		mockIdentityKeeper,
		mockTSSProtocol,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return k, ctx
}

// CreateTestSigningSession creates a signing session for testing
func CreateTestSigningSession(t testing.TB, k *keeper.Keeper, ctx sdk.Context, wallet *types.Wallet) *types.SigningSession {
	now := time.Now().UTC()
	session := &types.SigningSession{
		SessionId:  "test_session",
		WalletId:   wallet.WalletAddress,
		Message:    []byte("test_message"),
		Status:     types.SigningStatus_SIGNING_STATUS_IN_PROGRESS,
		Parties:    []string{"party1", "party2", "party3"},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	err := k.SaveSigningSession(ctx, session)
	require.NoError(t, err)

	return session
}
