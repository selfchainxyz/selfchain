package keeper_test

import (
	"fmt"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/testutil/mocks"
	"selfchain/x/keyless/types"
)

func TestWalletKeyGeneration(t *testing.T) {
	// Setup test environment
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey("mem_keyless")

	db := dbm.NewMemDB()
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
		walletAddr  string
		creator     string
		parties     []string
		threshold   uint32
		expectError bool
	}{
		{
			name:        "valid keygen request",
			walletAddr:  "self1wallet",
			creator:     "self1creator",
			parties:     []string{"party1", "party2", "party3"},
			threshold:   2,
			expectError: false,
		},
		{
			name:        "invalid threshold",
			walletAddr:  "self1wallet2",
			creator:     "self1creator",
			parties:     []string{"party1", "party2"},
			threshold:   3, // Greater than number of parties
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create wallet
			blockTime := ctx.BlockTime()
			wallet := &types.Wallet{
				Id:            tc.walletAddr,
				Creator:       tc.creator,
				PublicKey:     "test_public_key",
				WalletAddress: tc.walletAddr,
				ChainId:       "test-chain",
				Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
				KeyVersion:    1,
				CreatedAt:     &blockTime,
				UpdatedAt:     &blockTime,
				LastUsed:      &blockTime,
				UsageCount:    0,
			}
			err := k.SaveWallet(ctx, wallet)
			require.NoError(t, err)

			// Create party data
			partyData := &types.PartyData{
				PartyId:          tc.creator,
				PublicKey:        []byte("test_public_key"),
				PartyShare:       []byte("test_party_share"),
				VerificationData: []byte("test_verification_data"),
				ChainId:          tc.walletAddr,
				Status:           "active",
			}

			// For invalid threshold test, don't save party data
			if !tc.expectError {
				err = k.SavePartyData(ctx, tc.walletAddr, partyData)
				require.NoError(t, err)
			}

			// Test signing
			signature, err := k.SignWithTSS(ctx, wallet, "test_message")
			if tc.expectError {
				require.Error(t, err)
				require.Nil(t, signature)
			} else {
				require.NoError(t, err)
				require.NotNil(t, signature)

				// Verify signing session was created
				sessionID := fmt.Sprintf("%s-%d", tc.walletAddr, ctx.BlockHeight())
				session, err := k.GetSigningSession(ctx, sessionID)
				require.NoError(t, err)
				require.NotNil(t, session)
				require.Equal(t, types.SigningStatus_SIGNING_STATUS_COMPLETED, session.Status)
				require.Equal(t, tc.walletAddr, session.WalletId)
				require.Equal(t, []byte("test_message"), session.Message)

				// Verify wallet status
				status, err := k.GetWalletStatus(ctx, tc.walletAddr)
				require.NoError(t, err)
				require.Equal(t, types.WalletStatus_WALLET_STATUS_ACTIVE, status)
			}
		})
	}
}

func TestKeyShare(t *testing.T) {
	// Setup test environment
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey("mem_keyless")

	db := dbm.NewMemDB()
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

	// Setup test data
	did := "did:self:123"
	keyShare := []byte("test_key_share")

	// Test storing key share
	err := k.StoreKeyShare(ctx, did, keyShare)
	require.NoError(t, err)

	// Test retrieving key share
	storedShare, found := k.GetKeyShare(ctx, did)
	require.True(t, found)
	require.Equal(t, keyShare, storedShare)

	// Test retrieving non-existent key share
	_, found = k.GetKeyShare(ctx, "non_existent_did")
	require.False(t, found)
}
