package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

func (k msgServer) RevokePermission(goCtx context.Context, msg *types.MsgRevokePermission) (*types.MsgRevokePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Revoke permission
	err := k.Keeper.RevokePermission(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke permission: %v", err)
	}

	return &types.MsgRevokePermissionResponse{
		WalletId: msg.WalletId,
		Grantee:  msg.Grantee,
	}, nil
}
