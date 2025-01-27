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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	identitytypes "selfchain/x/identity/types"
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
	
	// Set up mock expectations
	identityKeeper.On("GetDIDDocument", mock.Anything, mock.Anything).Return(identitytypes.DIDDocument{}, true)
	identityKeeper.On("VerifyDIDOwnership", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("VerifyRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("GetKeyShare", mock.Anything, mock.Anything).Return([]byte("test_key_share"), true)
	identityKeeper.On("ReconstructWallet", mock.Anything, mock.Anything).Return(nil, nil)
	identityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil)
	identityKeeper.On("GenerateRecoveryToken", mock.Anything, mock.Anything).Return("test_recovery_token", nil)
	identityKeeper.On("ValidateRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create mock TSS protocol
	mockTSS := &mocks.MockTSSProtocol{
		InitiateSigningFn: func(ctx context.Context, msg []byte, walletID string) (*types.SigningResponse, error) {
			now := time.Now().UTC()
			return &types.SigningResponse{
				WalletId:  walletID,
				Signature: []byte("test_signature"),
				Metadata: &types.SignatureMetadata{
					Timestamp: &now,
					ChainId:   "test-chain",
					SignType:  types.SignatureType_SIGNATURE_TYPE_ECDSA,
				},
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
		Time:    time.Now().UTC(),
	}
	ctx := sdk.NewContext(stateStore, header, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return &KeylessKeeper{
		Keeper: k,
		Ctx:    ctx,
	}
}
