package types_test

import (
	"testing"

	"selfchain/x/migration/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

				TokenMigrationList: []types.TokenMigration{
					{
						MsgHash:   "0",
						Processed: false,
					},
					{
						MsgHash:   "1",
						Processed: false,
					},
				},
				Acl: &types.Acl{
					Admin: "25",
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
					VestingDuration:    87,
					VestingCliff:       91,
					MinMigrationAmount: 80,
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated tokenMigration",
			genState: &types.GenesisState{
				TokenMigrationList: []types.TokenMigration{
					{
						MsgHash:   "0",
						Processed: false,
					},
					{
						MsgHash:   "0",
						Processed: false,
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated migrator",
			genState: &types.GenesisState{
				MigratorList: []types.Migrator{
					{
						Migrator: "0",
					},
					{
						Migrator: "0",
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
