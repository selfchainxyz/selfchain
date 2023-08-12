package selfvesting_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "selfchain/testutil/keeper"
	"selfchain/testutil/nullify"
	"selfchain/x/selfvesting"
	"selfchain/x/selfvesting/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		VestingPositionsList: []types.VestingPositions{
			{
				Beneficiary: "0",
			},
			{
				Beneficiary: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SelfvestingKeeper(t)
	selfvesting.InitGenesis(ctx, *k, genesisState)
	got := selfvesting.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.VestingPositionsList, got.VestingPositionsList)
	// this line is used by starport scaffolding # genesis/test/assert
}
