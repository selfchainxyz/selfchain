package keeper

import (
	"context"

	"selfchain/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddMigrator(goCtx context.Context, msg *types.MsgAddMigrator) (*types.MsgAddMigratorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	acl, aclExists := k.GetAcl(ctx)
	if !aclExists {
		panic("ACL does not exist")
	}

	if acl.Admin != msg.Creator {
		return nil, types.ErrOnlyAdmin
	}

	// Store migrator. If it exists it will simply overwrite it
	k.SetMigrator(ctx, types.Migrator{
		Migrator: msg.Migrator,
		Exists:   true,
	})

	return &types.MsgAddMigratorResponse{}, nil
}
