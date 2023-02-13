package keeper

import (
	v002 "frontier/x/migration/migrations/v002"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc"
)

type Migrator struct {
	keeper      Keeper
	queryServer grpc.Server
}

func NewMigrator(keeper Keeper, queryServer grpc.Server) Migrator {
	return Migrator{keeper: keeper, queryServer: queryServer}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v002.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}
