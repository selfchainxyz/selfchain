package keeper

import (
    "context"
    "fmt"

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

    ctx := sdk.UnwrapSDKContext(c)
    
    wallet, err := k.GetWalletState(ctx, req.Address)
    if err != nil {
        return nil, fmt.Errorf("failed to get wallet: %w", err)
    }

    return &types.QueryGetWalletResponse{
        Address: wallet.Address,
        Did:     wallet.Did,
        Status:  wallet.Status,
    }, nil
}

// GetWalletByDID implements the Query/GetWalletByDID gRPC method
func (k Keeper) GetWalletByDID(c context.Context, req *types.QueryGetWalletByDIDRequest) (*types.QueryGetWalletByDIDResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    ctx := sdk.UnwrapSDKContext(c)
    
    wallet, err := k.GetWalletStateByDID(ctx, req.Did)
    if err != nil {
        return nil, fmt.Errorf("failed to get wallet by DID: %w", err)
    }

    return &types.QueryGetWalletByDIDResponse{
        Address: wallet.Address,
        Did:     wallet.Did,
        Status:  wallet.Status,
    }, nil
}
