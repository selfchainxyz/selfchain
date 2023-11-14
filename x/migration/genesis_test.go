package migration_test

import (
	"testing"

	keepertest "selfchain/testutil/keeper"
	"selfchain/testutil/nullify"
	"selfchain/x/migration"
	"selfchain/x/migration/types"

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
		Acl: &types.Acl{
			Admin: "56",
		},
		MigratorList: []types.Migrator{
			{
				Migrator: "0",
			},
			{
				Migrator: "1",
			},
		},
		Config: &types.Config{
			VestingDuration:    52,
			VestingCliff:       79,
			MinMigrationAmount: 70,
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
	require.Equal(t, genesisState.Acl, got.Acl)
	require.ElementsMatch(t, genesisState.MigratorList, got.MigratorList)
	require.Equal(t, genesisState.Config, got.Config)
	// this line is used by starport scaffolding # genesis/test/assert
}
