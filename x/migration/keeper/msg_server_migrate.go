package keeper

import (
	"context"
	"crypto/sha256"
	"fmt"
	"frontier/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Migrate(goCtx context.Context, msg *types.MsgMigrate) (*types.MsgMigrateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. Make sure signer is in the list of migrators
	_, migratorExist := k.GetMigrator(ctx, msg.Creator); if !migratorExist {
		return nil, types.ErrUnknownMigrator
	}

	// 2. Create a hash of the message
	encodedMsg := fmt.Sprintf(
		"%s|%s|%s|%d|%s",
		msg.EthAddress,
		msg.DestAddress,
		msg.Amount,
		msg.Token,
		msg.TxHash,
	)
	msgHash := fmt.Sprintf("%x", sha256.Sum256([]byte(encodedMsg)))

	// 3. Check if message i.e. migration request has been processed already
	_, migrationExists := k.GetTokenMigration(ctx, msgHash); if migrationExists {
		return nil, types.ErrMigrationProcessed
	}

	// 4. Mint new tokens to the destAddress

	// TODO: Handling the message
	_ = ctx

	return &types.MsgMigrateResponse{}, nil
}
