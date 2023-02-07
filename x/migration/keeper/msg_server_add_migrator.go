package keeper

import (
	"context"

	"frontier/x/migration/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddMigrator(goCtx context.Context, msg *types.MsgAddMigrator) (*types.MsgAddMigratorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgAddMigratorResponse{}, nil
}
