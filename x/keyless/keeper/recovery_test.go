package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
	identitytypes "selfchain/x/identity/types"
)

func TestRecoverWallet(t *testing.T) {
	k, ctx, identityKeeper := setupKeeper(t)

	// Create test wallet and DID
	walletID := "wallet123"
	did := "did:self:123"
	recoveryToken := "valid_token"

	// Setup DID document
	didDoc := &identitytypes.DIDDocument{
		Id:         did,
		Controller: "owner123",
	}
	err := identityKeeper.StoreDIDDocument(ctx, didDoc)
	require.NoError(t, err)

	// Setup recovery token
	err = identityKeeper.StoreRecoveryToken(ctx, did, recoveryToken)
	require.NoError(t, err)

	// Setup key share
	keyShare := []byte("test_key_share")
	err = identityKeeper.StoreKeyShare(ctx, did, keyShare)
	require.NoError(t, err)

	// Test successful recovery
	msg := &types.MsgRecoverWallet{
		WalletId:      walletID,
		Did:           did,
		RecoveryToken: recoveryToken,
	}
	err = k.RecoverWallet(ctx, msg)
	require.NoError(t, err)

	// Verify wallet was reconstructed
	wallet, err := k.GetWallet(ctx, walletID)
	require.NoError(t, err)
	require.NotNil(t, wallet)

	// Test invalid DID
	msg.Did = "invalid_did"
	err = k.RecoverWallet(ctx, msg)
	require.Error(t, err)

	// Test invalid recovery token
	msg.Did = did
	msg.RecoveryToken = "invalid_token"
	err = k.RecoverWallet(ctx, msg)
	require.Error(t, err)
}

func TestVerifyIdentity(t *testing.T) {
	k, ctx, identityKeeper := setupKeeper(t)

	did := "did:self:123"
	oauthToken := "valid_token"

	// Setup DID document
	didDoc := &identitytypes.DIDDocument{
		Id:         did,
		Controller: "owner123",
	}
	err := identityKeeper.StoreDIDDocument(ctx, didDoc)
	require.NoError(t, err)

	// Setup OAuth2 token
	err = identityKeeper.StoreOAuth2Token(ctx, did, oauthToken)
	require.NoError(t, err)

	// Test successful verification
	err = k.VerifyIdentity(ctx, did, oauthToken)
	require.NoError(t, err)

	// Test with MFA required
	k.SetSecurityLevel(ctx, types.SecurityLevelHigh)
	
	// Setup valid MFA session
	err = identityKeeper.StoreMFASession(ctx, did)
	require.NoError(t, err)

	err = k.VerifyIdentity(ctx, did, oauthToken)
	require.NoError(t, err)

	// Test rate limiting
	for i := 0; i < 100; i++ {
		_ = k.VerifyIdentity(ctx, did, oauthToken)
	}
	err = k.VerifyIdentity(ctx, did, oauthToken)
	require.Error(t, err)
	require.Contains(t, err.Error(), "rate limit exceeded")
}

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context, *mocks.IdentityKeeper) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := sdk.NewMemoryStoreKey(types.MemStoreKey)
	tKey := sdk.NewTransientStoreKey(types.TStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(tKey, sdk.StoreTypeTransient, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	identityKeeper := mocks.NewIdentityKeeper(t)

	paramsSubspace := typesparams.NewSubspace(cdc,
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
		identityKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return k, ctx, identityKeeper
}
