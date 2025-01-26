package keeper

import (
	"context"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/stretchr/testify/require"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/testutil/mocks"
	"selfchain/x/keyless/types"
)

// KeylessKeeper is a keeper that contains all the necessary information to test the keyless module
type KeylessKeeper struct {
	*keeper.Keeper
	Ctx sdk.Context
}

// NewKeylessKeeper creates a new keyless keeper for testing
func NewKeylessKeeper(t testing.TB) *KeylessKeeper {
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

	// Create mock TSS protocol
	mockTSS := &mocks.MockTSSProtocol{
		InitiateSigningFn: func(ctx context.Context, msg []byte, walletID string) (*types.SigningResponse, error) {
			return &types.SigningResponse{
				Signature: []byte("test_signature"),
			}, nil
		},
	}

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		identityKeeper,
		mockTSS,
	)

	header := tmproto.Header{
		ChainID: "test-chain",
		Height:  1,
		Time:    time.Now(),
	}
	ctx := sdk.NewContext(stateStore, header, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return &KeylessKeeper{
		Keeper: k,
		Ctx:    ctx,
	}
}
