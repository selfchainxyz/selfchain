package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/keeper"
	"selfchain/x/identity/types"
)

func TestVerifyDIDOwnership(t *testing.T) {
	k, ctx := setupKeeper(t)
	
	// Create test DID document
	did := "did:self:123"
	owner := sdk.AccAddress("test_owner")
	doc := types.DIDDocument{
		Id:         did,
		Controller: owner.String(),
	}

	// Store DID document
	store := k.DidStore(ctx)
	store.Set([]byte(did), k.Cdc().MustMarshal(&doc))

	// Test valid ownership
	err := k.VerifyDIDOwnership(ctx, did, owner)
	require.NoError(t, err)

	// Test invalid ownership
	invalidOwner := sdk.AccAddress("invalid_owner")
	err = k.VerifyDIDOwnership(ctx, did, invalidOwner)
	require.Error(t, err)
}

func TestVerifyOAuth2Token(t *testing.T) {
	k, ctx := setupKeeper(t)

	// Create test DID document
	did := "did:self:123"
	doc := types.DIDDocument{
		Id: did,
	}
	store := k.DidStore(ctx)
	store.Set([]byte(did), k.Cdc().MustMarshal(&doc))

	// Create valid token
	validToken := createTestToken(t, did)
	err := k.VerifyOAuth2Token(ctx, did, validToken)
	require.NoError(t, err)

	// Test invalid token
	invalidToken := "invalid.token"
	err = k.VerifyOAuth2Token(ctx, did, invalidToken)
	require.Error(t, err)
}

func TestVerifyMFA(t *testing.T) {
	k, ctx := setupKeeper(t)

	did := "did:self:123"
	
	// Create valid MFA session
	session := types.MFASession{
		Did:       did,
		CreatedAt: sdk.NewInt(time.Now().Unix()),
	}
	store := k.MfaStore(ctx)
	store.Set([]byte(did), k.Cdc().MustMarshal(&session))

	// Test valid MFA
	err := k.VerifyMFA(ctx, did)
	require.NoError(t, err)

	// Test expired MFA
	expiredSession := types.MFASession{
		Did:       did,
		CreatedAt: sdk.NewInt(time.Now().Add(-31 * time.Minute).Unix()),
	}
	store.Set([]byte(did), k.Cdc().MustMarshal(&expiredSession))
	err = k.VerifyMFA(ctx, did)
	require.Error(t, err)
}

func TestCheckRateLimit(t *testing.T) {
	k, ctx := setupKeeper(t)

	did := "did:self:123"
	operation := "test_operation"

	// Test first operation
	err := k.CheckRateLimit(ctx, did, operation)
	require.NoError(t, err)

	// Test rate limit exceeded
	for i := 0; i < 100; i++ {
		_ = k.CheckRateLimit(ctx, did, operation)
	}
	err = k.CheckRateLimit(ctx, did, operation)
	require.Error(t, err)
}

func TestLogAuditEvent(t *testing.T) {
	k, ctx := setupKeeper(t)

	event := &types.AuditEvent{
		Did:       "did:self:123",
		EventType: "test_event",
		Success:   true,
		Metadata: map[string]string{
			"test": "value",
		},
	}

	err := k.LogAuditEvent(ctx, event)
	require.NoError(t, err)

	// Verify event was stored
	store := k.AuditStore(ctx)
	key := []byte(fmt.Sprintf("%d:%s", event.Timestamp, event.Did))
	bz := store.Get(key)
	require.NotNil(t, bz)

	var storedEvent types.AuditEvent
	k.Cdc().MustUnmarshal(bz, &storedEvent)
	require.Equal(t, event.Did, storedEvent.Did)
	require.Equal(t, event.EventType, storedEvent.EventType)
	require.Equal(t, event.Success, storedEvent.Success)
}

func setupKeeper(t *testing.T) (keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := sdk.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return k, ctx
}
