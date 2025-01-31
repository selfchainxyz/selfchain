package testutil

import (
	"testing"
	"time"

	identitytypes "selfchain/x/identity/types"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/testutil/mocks"
	"selfchain/x/keyless/types"

	tmdb "github.com/cometbft/cometbft-db"
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

// NewTestKeeper creates a new keeper for testing
func NewTestKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey("mem_keyless")
	paramStoreKey := storetypes.NewKVStoreKey("params")

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(paramStoreKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := paramtypes.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"KeylessParams",
	)

	identityKeeper := mocks.NewIdentityKeeper()
	now := time.Now().UTC()
	identityKeeper.On("GetDid", mock.Anything, mock.Anything).Return(&identitytypes.DIDDocument{
		Id:         "test_did",
		Controller: []string{"test_creator"},
		Status:     identitytypes.Status_STATUS_ACTIVE,
	}, nil)

	identityKeeper.On("ValidateRecoveryProof", mock.Anything, mock.Anything).Return(&identitytypes.Proof{
		Type:               "test_type",
		VerificationMethod: "test_verifier",
		Purpose:            "recovery",
		Created:            &now,
		Value:              "test_signature",
	}, nil)

	identityKeeper.On("GetRecoveryProof", mock.Anything, mock.Anything).Return(&identitytypes.Proof{
		Type:               "test_type",
		VerificationMethod: "test_verifier",
		Purpose:            "recovery",
		Created:            &now,
		Value:              "test_signature",
	}, nil)

	identityKeeper.On("GetDidByVerifier", mock.Anything, mock.Anything).Return(&identitytypes.DIDDocument{
		Id:         "test_did",
		Controller: []string{"test_creator"},
		Status:     identitytypes.Status_STATUS_ACTIVE,
	}, nil)

	tssProtocol := mocks.NewTSSProtocol()

	tssProtocol.On("Sign", mock.Anything, mock.Anything).Return(&types.SigningResponse{
		WalletAddress: "test_wallet",
		Signature:     []byte("test_signature"),
		Metadata: &types.SignatureMetadata{
			Timestamp: &now,
			ChainId:   "test-chain",
			SignType:  types.SignatureType_SIGNATURE_TYPE_ECDSA,
		},
	}, nil)

	tssProtocol.On("GenerateKey", mock.Anything, mock.Anything).Return(&types.KeyGenResponse{
		WalletAddress: "test_wallet",
		PublicKey:     []byte("test_pubkey"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     now,
			LastRotated:   now,
			LastUsed:      now,
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
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

	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}

// CreateTestSigningSession creates a signing session for testing
func CreateTestSigningSession(t testing.TB, k *keeper.Keeper, ctx sdk.Context) *types.SigningSession {
	now := time.Now().UTC()
	session := &types.SigningSession{
		SessionId: "test_session",
		WalletId:  "test_wallet",
		Message:   []byte("test_message"),
		Status:    types.SigningStatus_SIGNING_STATUS_IN_PROGRESS,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := k.SaveSigningSession(ctx, session)
	require.NoError(t, err)

	return session
}

// CreateTestKeygenSession creates a keygen session for testing
func CreateTestKeygenSession(t testing.TB, k *keeper.Keeper, ctx sdk.Context, wallet *types.Wallet) *types.SigningSession {
	now := time.Now().UTC()
	session := &types.SigningSession{
		SessionId: "test_session",
		WalletId:  wallet.WalletAddress,
		Message:   []byte("test_message"),
		Status:    types.SigningStatus_SIGNING_STATUS_IN_PROGRESS,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := k.SaveSigningSession(ctx, session)
	require.NoError(t, err)

	return session
}

// CreateTestWallet2 creates a wallet for testing
func CreateTestWallet2(t testing.TB, k *keeper.Keeper, ctx sdk.Context) *types.Wallet {
	now := time.Now().UTC()
	wallet := &types.Wallet{
		WalletAddress: "test_wallet",
		Creator:       "test_creator",
		PublicKey:     "test_pubkey",
		ChainId:       "test-1",
		Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	err := k.SaveWallet(ctx, wallet)
	require.NoError(t, err)

	return wallet
}
