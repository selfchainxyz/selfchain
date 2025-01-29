package keeper

import (
	"context"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"selfchain/x/keyless/types"
)

func (k msgServer) GrantPermission(goCtx context.Context, msg *types.MsgGrantPermission) (*types.MsgGrantPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Check if the wallet exists
	wallet, err := k.GetWallet(ctx, msg.WalletAddress)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "wallet not found: %s", msg.WalletAddress)
	}

	// Check if the sender is the wallet owner
	if msg.Creator != wallet.Creator {
		return nil, status.Error(codes.PermissionDenied, "only wallet owner can grant permissions")
	}

	// Check if the permission already exists
	existingPerm, err := k.GetPermission(ctx, msg.WalletAddress, msg.Grantee)
	if err == nil && !existingPerm.IsExpired() && !existingPerm.IsRevoked() {
		return nil, status.Error(codes.AlreadyExists, "permission already exists and is active")
	}

	// Check if expiry time is valid
	if msg.ExpiresAt == nil {
		return nil, status.Error(codes.InvalidArgument, "expiry time must be specified")
	}
	if msg.ExpiresAt.Before(time.Now()) {
		return nil, status.Error(codes.InvalidArgument, "expiry time must be in the future")
	}

	// Convert permissions to strings
	permStrings := make([]string, len(msg.Permissions))
	for i, p := range msg.Permissions {
		permStrings[i] = p.String()
	}

	// Create current time as pointer
	now := time.Now()

	// Create and store the permission
	permission := &types.Permission{
		WalletAddress: msg.WalletAddress,
		Grantee:      msg.Grantee,
		Permissions:  permStrings,
		GrantedAt:    &now,
		ExpiresAt:    msg.ExpiresAt,
		Revoked:      false,
	}

	err = k.Keeper.GrantPermission(ctx, permission)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store permission: %v", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeGrantPermission,
			sdk.NewAttribute(types.AttributeKeyWalletAddress, msg.WalletAddress),
			sdk.NewAttribute(types.AttributeKeyGrantee, msg.Grantee),
			sdk.NewAttribute(types.AttributeKeyPermissions, strings.Join(permStrings, ",")),
		),
	)

	return &types.MsgGrantPermissionResponse{}, nil
}
