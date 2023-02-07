package keeper

import (
	"context"

	"frontier/x/migration/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RemoveMigrator(goCtx context.Context, msg *types.MsgRemoveMigrator) (*types.MsgRemoveMigratorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgRemoveMigratorResponse{}, nil
}
