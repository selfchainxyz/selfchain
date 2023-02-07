package keeper

import (
	"context"

	"frontier/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RemoveMigrator(goCtx context.Context, msg *types.MsgRemoveMigrator) (*types.MsgRemoveMigratorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	acl, aclExists := k.GetAcl(ctx)
	if !aclExists {
		panic("ACL does not exist")
	}

	if acl.Admin != msg.Creator {
		return nil, types.ErrOnlyAdmin
	}

	k.Keeper.RemoveMigrator(ctx,  msg.Migrator);

	return &types.MsgRemoveMigratorResponse{}, nil
}
