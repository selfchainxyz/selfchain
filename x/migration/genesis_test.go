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

		TokenMigrationList: []types.TokenMigration{
			{
				MsgHash: "0",
			},
			{
				MsgHash: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MigrationKeeper(t)
	migration.InitGenesis(ctx, *k, genesisState)
	got := migration.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.TokenMigrationList, got.TokenMigrationList)
	// this line is used by starport scaffolding # genesis/test/assert
}
