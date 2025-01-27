package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

func (k msgServer) GrantPermission(goCtx context.Context, msg *types.MsgGrantPermission) (*types.MsgGrantPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Grant permission
	permission, err := k.Keeper.GrantPermission(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to grant permission: %v", err)
	}

	return &types.MsgGrantPermissionResponse{
		WalletId: permission.WalletId,
		Grantee:  permission.Grantee,
	}, nil
}
