package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cometbft/cometbft-db"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
	identitytypes "selfchain/x/identity/types"
	"selfchain/x/keyless/testutil/mocks"
)

func TestRecoverWallet(t *testing.T) {
	// Setup test environment
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

	// Create test data
	walletAddr := "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
	creator := "cosmos1creator"
	recoveryProof := "valid_proof"
	newPubKey := "new_pub_key"
	signature := "valid_signature"

	// Create DID document
	didDoc := identitytypes.DIDDocument{
		Id:         creator,
		Controller: []string{"owner123"},
	}

	// Setup mock expectations with any context
	identityKeeper.On("GetDIDDocument", mock.Anything, creator).Return(didDoc, true)
	identityKeeper.On("VerifyRecoveryToken", mock.Anything, creator, recoveryProof).Return(nil)
	identityKeeper.On("VerifyDIDOwnership", mock.Anything, creator, sdk.AccAddress{}).Return(nil)
	identityKeeper.On("GetKeyShare", mock.Anything, creator).Return([]byte("key_share"), true)
	identityKeeper.On("ReconstructWallet", mock.Anything, didDoc).Return(&types.Wallet{
		PublicKey:   newPubKey,
		KeyVersion:  1,
		ChainId:     "test-1",
	}, nil)

	// Create keeper
	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		identityKeeper,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, nil)

	// Test successful recovery
	msg := &types.MsgRecoverWallet{
		Creator:       creator,
		WalletAddress: walletAddr,
		RecoveryProof: recoveryProof,
		NewPubKey:     newPubKey,
		Signature:     signature,
	}
	err := k.RecoverWallet(ctx, msg)
	require.NoError(t, err)

	// Test invalid recovery proof
	msg.RecoveryProof = "invalid_proof"
	identityKeeper.On("VerifyRecoveryToken", mock.Anything, creator, "invalid_proof").Return(types.ErrInvalidRecoveryProof)
	err = k.RecoverWallet(ctx, msg)
	require.Error(t, err)
}

func TestVerifyIdentity(t *testing.T) {
	// Setup test environment
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

	// Create test data
	did := "did:self:123"
	oauthToken := "valid_token"

	// Create DID document
	didDoc := identitytypes.DIDDocument{
		Id:         did,
		Controller: []string{"owner123"},
	}

	// Setup mock expectations with any context
	identityKeeper.On("GetDIDDocument", mock.Anything, did).Return(didDoc, true)
	identityKeeper.On("VerifyOAuth2Token", mock.Anything, did, oauthToken).Return(nil)
	identityKeeper.On("VerifyMFA", mock.Anything, did).Return(nil)

	// First call to CheckRateLimit should succeed
	identityKeeper.On("CheckRateLimit", mock.Anything, did, "identity_verification").Return(nil).Once()
	identityKeeper.On("LogAuditEvent", mock.Anything, mock.MatchedBy(func(event *identitytypes.AuditEvent) bool {
		return event.Did == did && event.EventType == "identity_verification" && event.Success
	})).Return(nil)

	// Create keeper
	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		identityKeeper,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, nil)

	// Test successful verification
	err := k.VerifyIdentity(ctx, did, oauthToken)
	require.NoError(t, err)

	// Test with high security level
	params := types.DefaultParams()
	params.MaxSecurityLevel = types.DefaultMaxSecurityLevel
	k.SetParams(ctx, params)

	// Second call to CheckRateLimit should succeed
	identityKeeper.On("CheckRateLimit", mock.Anything, did, "identity_verification").Return(nil).Once()
	err = k.VerifyIdentity(ctx, did, oauthToken)
	require.NoError(t, err)

	// Test rate limiting - third call should fail
	identityKeeper.On("CheckRateLimit", mock.Anything, did, "identity_verification").Return(types.ErrRateLimitExceeded).Once()
	err = k.VerifyIdentity(ctx, did, oauthToken)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrRateLimitExceeded)
}
