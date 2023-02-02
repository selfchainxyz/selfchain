package migration_test

import (
	"testing"

	keepertest "frontier/testutil/keeper"
	"frontier/testutil/nullify"
	"frontier/x/migration"
	"frontier/x/migration/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MigrationKeeper(t)
	migration.InitGenesis(ctx, *k, genesisState)
	got := migration.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
