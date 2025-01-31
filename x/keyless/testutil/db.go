package testutil

import (
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// SetupTestDB creates a new in-memory database for testing
func SetupTestDB(t *testing.T) storetypes.KVStore {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	key := storetypes.NewKVStoreKey("test")
	cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	err := cms.LoadLatestVersion()
	require.NoError(t, err)

	header := tmproto.Header{
		ChainID: "test-chain",
		Height:  1,
	}
	ctx := sdk.NewContext(cms, header, false, log.NewNopLogger())
	return prefix.NewStore(ctx.KVStore(key), []byte("test"))
}
