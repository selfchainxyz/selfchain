package testutil

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cometbft/cometbft-db"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/testutil/mocks"
	"selfchain/x/keyless/types"
	identitytypes "selfchain/x/identity/types"
)

// NewKeeper creates a new keeper for testing
func NewKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
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
		"KeylessParams",
	)

	// Mock the identity keeper
	identityKeeper := mocks.NewIdentityKeeper(t)
	identityKeeper.On("VerifyDIDOwnership", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	now := time.Now().UTC()
	identityKeeper.On("GetDIDDocument", mock.Anything, mock.Anything).Return(&identitytypes.DIDDocument{
		Id:         "test-did",
		Controller: []string{"test-controller"},
		Created:    &now,
		Updated:    &now,
		Status:     identitytypes.Status_STATUS_ACTIVE,
	}, nil)

	// Mock the TSS protocol
	tssProtocol := mocks.NewTSSProtocol(t)
	tssProtocol.On("GenerateKey", mock.Anything, mock.Anything).Return(&types.KeyGenResponse{
		WalletId:  "test-wallet",
		PublicKey: []byte("mock_public_key"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     time.Now(),
			LastRotated:   time.Now(),
			LastUsed:      time.Now(),
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
		},
	}, nil)
	tssProtocol.On("ProcessKeyGenRound", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	
	tssProtocol.On("InitiateSigning", mock.Anything, mock.Anything, mock.Anything).Return(&types.SigningResponse{
		WalletId:  "test-wallet",
		Signature: []byte("mock_signature"),
		Metadata: &types.SignatureMetadata{
			Timestamp: &now,
			ChainId:   "test-chain",
			SignType:  types.SignatureType_SIGNATURE_TYPE_ECDSA,
		},
	}, nil)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		identityKeeper,
		tssProtocol,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
