package testutil

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cometbft/cometbft-db"
)

// EncodingConfig specifies the concrete encoding types to use for the keyless module.
type EncodingConfig struct {
	Codec             codec.Codec
	Amino             *codec.LegacyAmino
	InterfaceRegistry codectypes.InterfaceRegistry
	Marshaler         codec.Codec // For backward compatibility
}

func MakeTestEncodingConfig() EncodingConfig {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	amino := codec.NewLegacyAmino()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	return EncodingConfig{
		Codec:             marshaler,
		Amino:             amino,
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler, // For backward compatibility
	}
}

func NewTestMultiStore(t *testing.T, storeKey storetypes.StoreKey) sdk.MultiStore {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.NoError(t, err)
	return ms
}

func NewTestContext(t *testing.T, storeKey storetypes.StoreKey) sdk.Context {
	return sdk.NewContext(
		NewTestMultiStore(t, storeKey),
		tmproto.Header{},
		false,
		log.NewNopLogger(),
	)
}
