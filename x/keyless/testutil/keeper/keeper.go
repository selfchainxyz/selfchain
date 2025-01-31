package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cometbft/cometbft-db"
)

type TestKeeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	memKey     storetypes.StoreKey
	paramstore typesparams.Subspace
}

func NewTestKeeper(t testing.TB) (*TestKeeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey("keyless")
	memStoreKey := storetypes.NewMemoryStoreKey("mem_keyless")

	testDB := db.NewMemDB()
	stateStore := store.NewCommitMultiStore(testDB)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, testDB)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		codec.NewLegacyAmino(),
		storeKey,
		memStoreKey,
		"KeylessParams",
	)

	k := &TestKeeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memStoreKey,
		paramstore: paramsSubspace,
	}

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	return k, ctx
}

func (k *TestKeeper) GetCodec() codec.BinaryCodec {
	return k.cdc
}

func (k *TestKeeper) GetStoreKey() storetypes.StoreKey {
	return k.storeKey
}

func (k *TestKeeper) GetMemKey() storetypes.StoreKey {
	return k.memKey
}

func (k *TestKeeper) GetParamStore() typesparams.Subspace {
	return k.paramstore
}
