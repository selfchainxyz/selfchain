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

func (k Keeper) TokenMigrationAll(goCtx context.Context, req *types.QueryAllTokenMigrationRequest) (*types.QueryAllTokenMigrationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var tokenMigrations []types.TokenMigration
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	tokenMigrationStore := prefix.NewStore(store, types.KeyPrefix(types.TokenMigrationKeyPrefix))

	pageRes, err := query.Paginate(tokenMigrationStore, req.Pagination, func(key []byte, value []byte) error {
		var tokenMigration types.TokenMigration
		if err := k.cdc.Unmarshal(value, &tokenMigration); err != nil {
			return err
		}

		tokenMigrations = append(tokenMigrations, tokenMigration)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTokenMigrationResponse{TokenMigration: tokenMigrations, Pagination: pageRes}, nil
}

func (k Keeper) TokenMigration(goCtx context.Context, req *types.QueryGetTokenMigrationRequest) (*types.QueryGetTokenMigrationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetTokenMigration(
		ctx,
		req.MsgHash,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetTokenMigrationResponse{TokenMigration: val}, nil
}
