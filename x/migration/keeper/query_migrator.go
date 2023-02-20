package keeper

import (
	"context"

	"selfchain/x/migration/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) MigratorAll(goCtx context.Context, req *types.QueryAllMigratorRequest) (*types.QueryAllMigratorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var migrators []types.Migrator
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	migratorStore := prefix.NewStore(store, types.KeyPrefix(types.MigratorKeyPrefix))

	pageRes, err := query.Paginate(migratorStore, req.Pagination, func(key []byte, value []byte) error {
		var migrator types.Migrator
		if err := k.cdc.Unmarshal(value, &migrator); err != nil {
			return err
		}

		migrators = append(migrators, migrator)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllMigratorResponse{Migrator: migrators, Pagination: pageRes}, nil
}

func (k Keeper) Migrator(goCtx context.Context, req *types.QueryGetMigratorRequest) (*types.QueryGetMigratorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetMigrator(
		ctx,
		req.Migrator,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetMigratorResponse{Migrator: val}, nil
}
