package identity_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "selfchain/testutil/keeper"
	"selfchain/x/identity"
	"selfchain/x/identity/types"
)

func TestGenesis(t *testing.T) {
	k, ctx := keepertest.IdentityKeeper(t)
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # genesis/test/state
	}

	identity.InitGenesis(ctx, *k, &genesisState)
	got := identity.ExportGenesis(ctx, *k)
	require.NotNil(t, got)
}
