package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cometbft/cometbft-db"
	
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

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
		Controller: []string{owner.String()},
		VerificationMethod: []*types.VerificationMethod{
			{
				Id:             did + "#key1",
				Type:          "Ed25519VerificationKey2020",
				Controller:    owner.String(),
				PublicKeyBase58: "test_key",
			},
		},
	}

	// Store DID document
	err := k.SetDIDDocument(ctx, did, doc)
	require.NoError(t, err)

	// Test valid ownership
	err = k.VerifyDIDOwnership(ctx, did, owner)
	require.NoError(t, err)

	// Test invalid ownership
	invalidOwner := sdk.AccAddress("invalid_owner")
	err = k.VerifyDIDOwnership(ctx, did, invalidOwner)
	require.Error(t, err)
}

func TestOAuthVerification(t *testing.T) {
	t.Skip("OAuth verification not implemented yet")
}

func TestMFAVerification(t *testing.T) {
	t.Skip("MFA verification not implemented yet")
}

func setupKeeper(t *testing.T) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
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
		"IdentityParams",
	)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		nil, // keyless keeper
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, nil)
	return k, ctx
}
