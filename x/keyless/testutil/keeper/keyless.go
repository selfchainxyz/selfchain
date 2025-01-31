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
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/testutil/mocks"
	"selfchain/x/keyless/types"
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

	// Create mock keepers
	mockTSSProtocol := mocks.NewTSSProtocol()
	mockIdentityKeeper := mocks.NewIdentityKeeper()

	// Set up mock expectations
	now := time.Now().UTC()
	mockTSSProtocol.On("GenerateKeyShares", 
		mock.Anything, // sdk.Context
		mock.AnythingOfType("string"), // walletAddress
		mock.AnythingOfType("uint32"), // threshold
		mock.AnythingOfType("types.SecurityLevel"), // securityLevel
	).Return(&types.KeyGenResponse{
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

	mockTSSProtocol.On("ReconstructKey", 
		mock.Anything, // sdk.Context
		mock.AnythingOfType("[][]byte"), // shares
	).Return([]byte("test-pubkey"), nil)

	mockTSSProtocol.On("SignMessage", 
		mock.Anything, // sdk.Context
		mock.AnythingOfType("[]byte"), // message
		mock.AnythingOfType("[][]byte"), // shares
	).Return([]byte("test-signature"), nil)

	mockTSSProtocol.On("VerifyShare", 
		mock.Anything, // sdk.Context
		mock.AnythingOfType("[]byte"), // share
		mock.AnythingOfType("[]byte"), // publicKey
	).Return(nil)

	mockTSSProtocol.On("VerifySignature", 
		mock.Anything, // sdk.Context
		mock.AnythingOfType("[]byte"), // message
		mock.AnythingOfType("[]byte"), // signature
		mock.AnythingOfType("[]byte"), // publicKey
	).Return(nil)

	mockTSSProtocol.On("GetPartyData",
		mock.Anything, // sdk.Context
		mock.AnythingOfType("string"), // partyID
	).Return(&types.PartyData{
		PartyId:    "test-party",
		PartyShare: []byte("test-share"),
		Status:     "active",
	}, nil)

	mockTSSProtocol.On("SetPartyData",
		mock.Anything, // sdk.Context
		mock.AnythingOfType("*types.PartyData"), // data
	).Return(nil)

	mockIdentityKeeper.On("GetDIDDocument", mock.Anything, mock.Anything).Return(identitytypes.DIDDocument{}, true)
	mockIdentityKeeper.On("SaveDIDDocument", mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("ReconstructWalletFromDID", mock.Anything, mock.Anything).Return([]byte("test-wallet"), nil)
	mockIdentityKeeper.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil)

	paramsSubspace := paramtypes.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"KeylessParams",
	)

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
