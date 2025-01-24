package testutil

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/stretchr/testify/require"
	tmdb "github.com/cometbft/cometbft-db"
)

// NewTestMultiStore creates a new MultiStore for testing
func NewTestMultiStore(t *testing.T, key storetypes.StoreKey) storetypes.MultiStore {
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	err := stateStore.LoadLatestVersion()
	require.NoError(t, err)
	return stateStore
}

// MakeTestEncodingConfig creates a test encoding config
func MakeTestEncodingConfig() codec.ProtoCodecMarshaler {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	return codec.NewProtoCodec(interfaceRegistry)
}
