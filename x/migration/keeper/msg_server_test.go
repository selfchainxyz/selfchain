package keeper_test

import (
	"context"
	"testing"

	commontest "frontier/testutil"
	keepertest "frontier/testutil/keeper"
	"frontier/testutil/sample"
	"frontier/x/migration/keeper"
	"frontier/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setup(t testing.TB) (types.MsgServer, context.Context) {
	commontest.InitSDKConfig()

	k, ctx := keepertest.MigrationKeeper(t)
	
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func TestShouldFailIfInvalidMigrator(t *testing.T) {
	// create a couple of migrators
	setup(t)
	migrator1 := sample.AccAddress()
	migrator1 := sample.AccAddress()
}
