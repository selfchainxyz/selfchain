package keeper

import (
    "context"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "selfchain/x/keyless/types"
)

var _ types.QueryServer = Keeper{}

// Params implements the Query/Params gRPC method
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    ctx := sdk.UnwrapSDKContext(c)

    return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// GetWallet implements the Query/GetWallet gRPC method
func (k Keeper) GetWallet(c context.Context, req *types.QueryGetWalletRequest) (*types.QueryGetWalletResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    _ = sdk.UnwrapSDKContext(c)
    // TODO: Implement get wallet logic
    // This will be implemented in the next step

    return &types.QueryGetWalletResponse{}, nil
}

// GetWalletByDID implements the Query/GetWalletByDID gRPC method
func (k Keeper) GetWalletByDID(c context.Context, req *types.QueryGetWalletByDIDRequest) (*types.QueryGetWalletByDIDResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    _ = sdk.UnwrapSDKContext(c)
    // TODO: Implement get wallet by DID logic
    // This will be implemented in the next step

    return &types.QueryGetWalletByDIDResponse{}, nil
}
