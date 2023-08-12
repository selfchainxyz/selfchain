package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"selfchain/x/selfvesting/types"
)

func (k Keeper) VestingPositionsAll(goCtx context.Context, req *types.QueryAllVestingPositionsRequest) (*types.QueryAllVestingPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var vestingPositionss []types.VestingPositions
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	vestingPositionsStore := prefix.NewStore(store, types.KeyPrefix(types.VestingPositionsKeyPrefix))

	pageRes, err := query.Paginate(vestingPositionsStore, req.Pagination, func(key []byte, value []byte) error {
		var vestingPositions types.VestingPositions
		if err := k.cdc.Unmarshal(value, &vestingPositions); err != nil {
			return err
		}

		vestingPositionss = append(vestingPositionss, vestingPositions)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllVestingPositionsResponse{VestingPositions: vestingPositionss, Pagination: pageRes}, nil
}

func (k Keeper) VestingPositions(goCtx context.Context, req *types.QueryGetVestingPositionsRequest) (*types.QueryGetVestingPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetVestingPositions(
		ctx,
		req.Beneficiary,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetVestingPositionsResponse{VestingPositions: val}, nil
}
