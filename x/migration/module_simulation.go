package migration

import (
	"math/rand"

	"selfchain/testutil/sample"
	migrationsimulation "selfchain/x/migration/simulation"
	"selfchain/x/migration/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = migrationsimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgMigrate = "op_weight_msg_migrate"
	// TODO: Determine the simulation weight value
	defaultWeightMsgMigrate int = 100

	opWeightMsgAddMigrator = "op_weight_msg_add_migrator"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddMigrator int = 100

	opWeightMsgRemoveMigrator = "op_weight_msg_remove_migrator"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRemoveMigrator int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	migrationGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&migrationGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgMigrate int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgMigrate, &weightMsgMigrate, nil,
		func(_ *rand.Rand) {
			weightMsgMigrate = defaultWeightMsgMigrate
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMigrate,
		migrationsimulation.SimulateMsgMigrate(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddMigrator int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAddMigrator, &weightMsgAddMigrator, nil,
		func(_ *rand.Rand) {
			weightMsgAddMigrator = defaultWeightMsgAddMigrator
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddMigrator,
		migrationsimulation.SimulateMsgAddMigrator(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgRemoveMigrator int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgRemoveMigrator, &weightMsgRemoveMigrator, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveMigrator = defaultWeightMsgRemoveMigrator
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRemoveMigrator,
		migrationsimulation.SimulateMsgRemoveMigrator(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
