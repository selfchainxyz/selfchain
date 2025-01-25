package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cometbft/cometbft-db"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
	"selfchain/x/keyless/testutil/mocks"
)

func TestWalletManagement(t *testing.T) {
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

	// Create keeper
	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		identityKeeper,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, nil)

	tests := []struct {
		name        string
		wallet      *types.Wallet
		expectError bool
	}{
		{
			name: "valid wallet creation",
			wallet: &types.Wallet{
				Creator:       "self1creator",
				WalletAddress: "self1wallet",
				ChainId:       "self-1",
				Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
				KeyVersion:    1,
			},
			expectError: false,
		},
		{
			name: "duplicate wallet address",
			wallet: &types.Wallet{
				Creator:       "self1creator",
				WalletAddress: "self1wallet",
				ChainId:       "self-1",
				Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
				KeyVersion:    1,
			},
			expectError: true,
		},
		{
			name: "invalid chain ID",
			wallet: &types.Wallet{
				Creator:       "self1creator",
				WalletAddress: "self1wallet3",
				ChainId:       "",
				Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
				KeyVersion:    1,
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := k.SaveWallet(ctx, tc.wallet)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				
				// Verify wallet was created correctly
				wallet, err := k.GetWallet(ctx, tc.wallet.WalletAddress)
				require.NoError(t, err)
				require.Equal(t, tc.wallet.Creator, wallet.Creator)
				require.Equal(t, tc.wallet.ChainId, wallet.ChainId)
				require.Equal(t, tc.wallet.Status, wallet.Status)
			}
		})
	}

	// Test getting all wallets
	wallets, err := k.GetAllWalletsFromStore(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, wallets)
}

func TestWalletAccess(t *testing.T) {
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

	// Create keeper
	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		identityKeeper,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, nil)

	// Create a test wallet
	wallet := &types.Wallet{
		Creator:       "self1creator",
		WalletAddress: "self1wallet",
		ChainId:       "self-1",
		Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
		KeyVersion:    1,
	}

	err := k.SaveWallet(ctx, wallet)
	require.NoError(t, err)

	// Test wallet access validation
	err = k.ValidateWalletAccess(ctx, wallet.WalletAddress, "sign")
	require.NoError(t, err)

	// Test wallet authorization
	authorized, err := k.IsWalletAuthorized(ctx, wallet.Creator, wallet.WalletAddress)
	require.NoError(t, err)
	require.True(t, authorized)

	// Test unauthorized access
	authorized, err = k.IsWalletAuthorized(ctx, "unauthorized_creator", wallet.WalletAddress)
	require.NoError(t, err)
	require.False(t, authorized)
}
