package keeper

import (
	"context"

	"selfchain/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateConfig(goCtx context.Context, msg *types.MsgUpdateConfig) (*types.MsgUpdateConfigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	acl, aclExists := k.GetAcl(ctx)
	if !aclExists {
		panic("ACL does not exist")
	}

	if acl.Admin != msg.Creator {
		return nil, types.ErrOnlyAdmin
	}

	// Store new config. If it exists it will simply overwrite it
	k.SetConfig(ctx, types.Config{
		VestingDuration:    msg.VestingDuration,
		VestingCliff:       msg.VestingCliff,
		MinMigrationAmount: msg.MinMigrationAmount,
	})

	return &types.MsgUpdateConfigResponse{}, nil
}
