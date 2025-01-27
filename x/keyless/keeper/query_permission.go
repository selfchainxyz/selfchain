package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"selfchain/x/keyless/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implement gRPC query service
func (k Keeper) ListPermissions(goCtx context.Context, req *types.QueryListPermissionsRequest) (*types.QueryListPermissionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)
	permissionStore := prefix.NewStore(store, []byte(types.PermissionKey))

	var permissions []*types.Permission
	pageRes, err := query.Paginate(permissionStore, req.Pagination, func(key []byte, value []byte) error {
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

	return &types.QueryListPermissionsResponse{
		Permissions: permissions,
		Pagination: pageRes,
	}, nil
}
