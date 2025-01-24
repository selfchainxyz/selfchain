package testutil

import (
	"testing"

	identitytypes "selfchain/x/identity/types"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
)

// MockIdentityKeeper is a mock implementation of the identity keeper for testing
type MockIdentityKeeper struct {
	didDocs map[string]*identitytypes.DIDDocument
}

// NewMockIdentityKeeper creates a new instance of MockIdentityKeeper
func NewMockIdentityKeeper() *MockIdentityKeeper {
	return &MockIdentityKeeper{
		didDocs: make(map[string]*identitytypes.DIDDocument),
	}
}

// SetDIDDocument sets a DID document in the mock keeper
func (m *MockIdentityKeeper) SetDIDDocument(doc *identitytypes.DIDDocument) {
	m.didDocs[doc.Id] = doc
}

// GetDIDDocument returns a mock DID document
func (m *MockIdentityKeeper) GetDIDDocument(ctx sdk.Context, did string) (identitytypes.DIDDocument, bool) {
	if doc, ok := m.didDocs[did]; ok {
		return *doc, true
	}
	return identitytypes.DIDDocument{}, false
}

// VerifyDIDOwnership returns nil for mock DID ownership verification
func (m *MockIdentityKeeper) VerifyDIDOwnership(ctx sdk.Context, did string, owner sdk.AccAddress) error {
	return nil
}

// VerifyOAuth2Token returns nil for mock OAuth2 token verification
func (m *MockIdentityKeeper) VerifyOAuth2Token(ctx sdk.Context, did string, token string) error {
	return nil
}

// VerifyMFA returns nil for mock MFA verification
func (m *MockIdentityKeeper) VerifyMFA(ctx sdk.Context, did string) error {
	return nil
}

// VerifyRecoveryToken returns nil for mock recovery token verification
func (m *MockIdentityKeeper) VerifyRecoveryToken(ctx sdk.Context, did string, token string) error {
	return nil
}

// GetKeyShare returns a mock key share if DID document exists
func (m *MockIdentityKeeper) GetKeyShare(ctx sdk.Context, did string) ([]byte, bool) {
	if doc, ok := m.didDocs[did]; ok {
		return []byte("mock_key_share_for_" + doc.Id), true
	}
	return nil, false
}

// ReconstructWallet returns a mock reconstructed wallet
func (m *MockIdentityKeeper) ReconstructWallet(ctx sdk.Context, didDoc identitytypes.DIDDocument) (interface{}, error) {
	return []byte("mock_reconstructed_wallet_for_" + didDoc.Id), nil
}

// CheckRateLimit returns nil for mock rate limit check
func (m *MockIdentityKeeper) CheckRateLimit(ctx sdk.Context, did string, operation string) error {
	return nil
}

// LogAuditEvent returns nil for mock audit event logging
func (m *MockIdentityKeeper) LogAuditEvent(ctx sdk.Context, event *identitytypes.AuditEvent) error {
	return nil
}

// NewTestKeeper creates a new keeper for testing
func NewTestKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	paramsStoreKey := storetypes.NewKVStoreKey("params")

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(paramsStoreKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"KeylessParams",
	)

	identityKeeper := NewMockIdentityKeeper()

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		identityKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
