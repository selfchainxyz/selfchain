package keyless_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"selfchain/x/keyless"
	"selfchain/x/keyless/testutil"
	"selfchain/x/keyless/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k, ctx := testutil.NewTestKeeper(t)
	keyless.InitGenesis(ctx, *k, genesisState)
	got := keyless.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	// Verify params
	require.Equal(t, genesisState.Params, got.Params)
}
