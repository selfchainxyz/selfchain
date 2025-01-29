package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"selfchain/x/keyless/types"
)

// Permissions returns all permissions for a wallet
func (k Keeper) Permissions(goCtx context.Context, req *types.QueryPermissionsRequest) (*types.QueryPermissionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validation
	if req.WalletId == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet ID cannot be empty")
	}

	// Check if wallet exists
	_, err := k.GetWallet(ctx, req.WalletId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "wallet not found: %s", req.WalletId)
	}

	// Get all permissions for the wallet
	store := ctx.KVStore(k.storeKey)
	permissionsStore := prefix.NewStore(store, types.GetPermissionPrefix(req.WalletId))

	var permissions []*types.Permission
	pageRes, err := query.Paginate(permissionsStore, req.Pagination, func(key []byte, value []byte) error {
		var permission types.Permission
		if err := k.cdc.Unmarshal(value, &permission); err != nil {
			return err
		}

		permissions = append(permissions, &permission)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPermissionsResponse{
		Permissions: permissions,
		Pagination: pageRes,
	}, nil
}

// Permission returns a specific permission for a wallet and grantee
func (k Keeper) Permission(goCtx context.Context, req *types.QueryPermissionRequest) (*types.QueryPermissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validation
	if req.WalletId == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet ID cannot be empty")
	}

	if req.Grantee == "" {
		return nil, status.Error(codes.InvalidArgument, "grantee cannot be empty")
	}

	// Check if wallet exists
	_, err := k.GetWallet(ctx, req.WalletId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "wallet not found: %s", req.WalletId)
	}

	// Get permission
	permission, err := k.GetPermission(ctx, req.WalletId, req.Grantee)
	if err != nil {
		return nil, status.Error(codes.NotFound, "permission not found")
	}

	return &types.QueryPermissionResponse{
		Permission: permission,
	}, nil
}
