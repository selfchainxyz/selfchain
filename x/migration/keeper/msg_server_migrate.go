package keeper

import (
	"context"

	"frontier/x/migration/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Migrate(goCtx context.Context, msg *types.MsgMigrate) (*types.MsgMigrateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgMigrateResponse{}, nil
}
