package keeper

import (
	"context"

	"selfchain/x/keyless/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func (k Keeper) KeyRotationStatus(goCtx context.Context, req *types.QueryKeyRotationStatusRequest) (*types.QueryKeyRotationStatusResponse, error) {
	if req == nil {
		return nil, grpcstatus.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the key rotation status from state
	status, found := k.GetKeyRotationStatus(ctx, req.WalletId)
	if !found {
		return nil, grpcstatus.Error(codes.NotFound, "key rotation status not found")
	}

	// Get active rotation if status is not UNSPECIFIED
	var version uint64
	var newPubKey string
	var error string
	if status != types.KeyRotationStatus_KEY_ROTATION_STATUS_UNSPECIFIED {
		// TODO: Get active rotation details
		version = 0
		newPubKey = ""
	}

	if status == types.KeyRotationStatus_KEY_ROTATION_STATUS_FAILED {
		error = "Key rotation failed"
	}

	return &types.QueryKeyRotationStatusResponse{
		WalletId:   req.WalletId,
		Status:     status.String(),
		Version:    version,
		NewPubKey:  newPubKey,
		Error:      error,
	}, nil
}
