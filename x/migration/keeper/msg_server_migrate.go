package keeper

import (
	"context"

	"frontier/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Migrate(goCtx context.Context, msg *types.MsgMigrate) (*types.MsgMigrateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. Make sure signer is in the list of migrators
	_, found := k.GetMigrator(ctx, msg.Creator); if !found {
		return nil, types.ErrUnknownMigrator
	}
	// 2. Create a hash of the message

	// 3. Check if message i.e. migration request has been processed already

	// 4. Mint new tokens to the destAddress

	// TODO: Handling the message
	_ = ctx

	return &types.MsgMigrateResponse{}, nil
}
