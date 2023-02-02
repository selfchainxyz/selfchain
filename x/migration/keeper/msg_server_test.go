package keeper_test

import (
	"context"
	"testing"

	keepertest "frontier/testutil/keeper"
	"frontier/x/migration/keeper"
	"frontier/x/migration/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.MigrationKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
