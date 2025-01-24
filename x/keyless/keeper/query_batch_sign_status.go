package keeper

import (
	"context"

	"selfchain/x/keyless/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func (k Keeper) BatchSignStatus(goCtx context.Context, req *types.QueryBatchSignStatusRequest) (*types.QueryBatchSignStatusResponse, error) {
	if req == nil {
		return nil, grpcstatus.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the batch sign status from state
	status, found := k.GetBatchSignStatus(ctx, req.WalletId, req.BatchId)
	if !found {
		return nil, grpcstatus.Error(codes.NotFound, "batch sign status not found")
	}

	// Get signatures if available
	var signatures []string
	if status == types.BatchSignStatus_BATCH_SIGN_STATUS_COMPLETED {
		// TODO: Implement signature retrieval
		signatures = []string{}
	}

	var error string
	if status == types.BatchSignStatus_BATCH_SIGN_STATUS_FAILED {
		error = "Batch sign failed"
	}

	return &types.QueryBatchSignStatusResponse{
		WalletId:   req.WalletId,
		BatchId:    req.BatchId,
		Status:     status.String(),
		Signatures: signatures,
		Error:      error,
	}, nil
}
