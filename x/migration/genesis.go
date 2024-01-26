package migration

import (
	"selfchain/x/migration/keeper"
	"selfchain/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the tokenMigration
	for _, elem := range genState.TokenMigrationList {
		k.SetTokenMigration(ctx, elem)
	}
	// Set if defined
	if genState.Acl != nil {
		k.SetAcl(ctx, *genState.Acl)
	}
	// Set all the migrator
	for _, elem := range genState.MigratorList {
		k.SetMigrator(ctx, elem)
	}
	// Set if defined
	if genState.Config != nil {
		k.SetConfig(ctx, *genState.Config)
	}

	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.TokenMigrationList = k.GetAllTokenMigration(ctx)
	// Get all acl
	acl, found := k.GetAcl(ctx)
	if found {
		genesis.Acl = &acl
	}
	genesis.MigratorList = k.GetAllMigrator(ctx)
	// Get all config
	config, found := k.GetConfig(ctx)
	if found {
		genesis.Config = &config
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
