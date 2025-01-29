package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"selfchain/x/keyless/types"
)

func (k msgServer) RevokePermission(goCtx context.Context, msg *types.MsgRevokePermission) (*types.MsgRevokePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate input parameters
	if msg.Grantee == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid grantee address: empty address string is not allowed")
	}
	if len(msg.Permissions) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one permission must be specified")
	}
	for _, p := range msg.Permissions {
		if p == types.WalletPermission_WALLET_PERMISSION_UNSPECIFIED {
			return nil, status.Error(codes.InvalidArgument, "invalid permission: WALLET_PERMISSION_UNSPECIFIED")
		}
	}

	// Validate creator is authorized
	authorized := k.IsWalletAuthorized(ctx, msg.Creator, msg.WalletAddress, types.WalletPermission_WALLET_PERMISSION_ADMIN)
	if !authorized {
		return nil, status.Error(codes.PermissionDenied, "not authorized")
	}

	// Check if creator has admin permission for revoking admin permissions
	for _, p := range msg.Permissions {
		if p == types.WalletPermission_WALLET_PERMISSION_ADMIN {
			hasAdmin := k.IsWalletAuthorized(ctx, msg.Creator, msg.WalletAddress, types.WalletPermission_WALLET_PERMISSION_ADMIN)
			if !hasAdmin {
				return nil, status.Error(codes.PermissionDenied, "admin permission requires additional verification")
			}
		}
	}

	// Revoke permissions
	err := k.Keeper.RevokePermission(ctx, msg.WalletAddress, msg.Grantee)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to revoke permission: %v", err))
	}

	// Convert permissions to strings for event
	permStrings := make([]string, len(msg.Permissions))
	for i, p := range msg.Permissions {
		permStrings[i] = p.String()
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePermissionRevoked,
			sdk.NewAttribute(types.AttributeKeyWalletAddress, msg.WalletAddress),
			sdk.NewAttribute(types.AttributeKeyGrantee, msg.Grantee),
			sdk.NewAttribute(types.AttributeKeyPermissions, fmt.Sprintf("%v", permStrings)),
			sdk.NewAttribute(types.AttributeKeyRevokedBy, msg.Creator),
		),
	)

	return &types.MsgRevokePermissionResponse{
		WalletAddress: msg.WalletAddress,
		Grantee:      msg.Grantee,
	}, nil
}
