package storage

import (
	"testing"
	"time"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	db "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestStorage(t *testing.T) (Storage, sdk.Context) {
	key := storetypes.NewKVStoreKey("test")
	db := db.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	err := cms.LoadLatestVersion()
	require.NoError(t, err)

	ctx := sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger())
	prefixStore := prefix.NewStore(ctx.KVStore(key), []byte("test"))
	storage := NewStorage(prefixStore)

	return storage, ctx
}

func TestStorage(t *testing.T) {
	storage, ctx := setupTestStorage(t)

	// Generate test party data
	preParams, err := keygen.GeneratePreParams(time.Minute)
	require.NoError(t, err)
	require.NotNil(t, preParams)

	party1Data := &keygen.LocalPartySaveData{
		LocalPreParams: *preParams,
	}
	party2Data := &keygen.LocalPartySaveData{
		LocalPreParams: *preParams,
	}

	walletAddress := "test_wallet"

	t.Run("Save and retrieve party data", func(t *testing.T) {
		// Save party data
		err := storage.SavePartyData(ctx, walletAddress, party1Data, party2Data)
		require.NoError(t, err)

		// Retrieve party data
		p1, p2, err := storage.GetPartyData(ctx, walletAddress)
		require.NoError(t, err)
		require.NotNil(t, p1)
		require.NotNil(t, p2)

		// Verify data matches
		assert.Equal(t, party1Data.LocalPreParams.PaillierSK, p1.LocalPreParams.PaillierSK)
		assert.Equal(t, party2Data.LocalPreParams.PaillierSK, p2.LocalPreParams.PaillierSK)
	})

	t.Run("Delete party data", func(t *testing.T) {
		// Delete party data
		err := storage.DeletePartyData(ctx, walletAddress)
		require.NoError(t, err)

		// Try to retrieve deleted data
		p1, p2, err := storage.GetPartyData(ctx, walletAddress)
		assert.Error(t, err)
		assert.Nil(t, p1)
		assert.Nil(t, p2)
	})

	t.Run("Get non-existent party data", func(t *testing.T) {
		p1, p2, err := storage.GetPartyData(ctx, "non_existent_wallet")
		assert.Error(t, err)
		assert.Nil(t, p1)
		assert.Nil(t, p2)
	})
}
