package keeper

import (
	"context"
	"fmt"
	"strings"

	"selfchain/x/keyless/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrantPermission handles granting permissions to a wallet
func (k msgServer) GrantPermission(goCtx context.Context, msg *types.MsgGrantPermission) (*types.MsgGrantPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic
	if err := msg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid message: %v", err)
	}

	// Check if wallet exists
	wallet, found := k.GetWallet(ctx, msg.WalletAddress)
	if !found {
		return nil, fmt.Errorf("wallet not found: %s", msg.WalletAddress)
	}

	// Check if sender is wallet owner
	if wallet.Creator != msg.Creator {
		return nil, fmt.Errorf("only wallet owner can grant permissions")
	}

	// Check if grantee exists
	if _, err := sdk.AccAddressFromBech32(msg.Grantee); err != nil {
		return nil, fmt.Errorf("invalid grantee address: %v", err)
	}

	// Check if permission already exists
	hasPermission, err := k.HasPermission(ctx, msg.WalletAddress, msg.Grantee, string(msg.Permissions[0]))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if hasPermission {
		return nil, status.Error(codes.AlreadyExists, "permission already exists")
	}

	// Create new permission
	now := ctx.BlockTime()
	permission := &types.Permission{
		WalletAddress: msg.WalletAddress,
		Grantee:      msg.Grantee,
		Permissions:  make([]string, len(msg.Permissions)),
		GrantedAt:    &now,
		ExpiresAt:    msg.ExpiresAt,
		Revoked:      false,
		RevokedAt:    nil,
	}

	// Convert WalletPermission enum values to strings
	for i, p := range msg.Permissions {
		permission.Permissions[i] = p.String()
	}

	// Validate and store permission
	if err := k.ValidateAndGrantPermission(ctx, permission); err != nil {
		return nil, fmt.Errorf("failed to grant permission: %v", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeGrantPermission,
			sdk.NewAttribute(types.AttributeKeyWalletAddress, msg.WalletAddress),
			sdk.NewAttribute(types.AttributeKeyGrantee, msg.Grantee),
			sdk.NewAttribute(types.AttributeKeyPermissions, strings.Join(permission.Permissions, ",")),
			sdk.NewAttribute(types.AttributeKeyTimestamp, now.String()),
		),
	)

	return &types.MsgGrantPermissionResponse{}, nil
}
