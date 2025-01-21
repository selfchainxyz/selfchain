package keeper

import (
	"context"

	"selfchain/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) VerifyOAuthToken(goCtx context.Context, msg *types.MsgVerifyOAuthToken) (*types.MsgVerifyOAuthTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify the token
	socialId, err := k.Keeper.VerifyOAuthToken(ctx, msg.Provider, msg.Token)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to verify OAuth token")
	}

	return &types.MsgVerifyOAuthTokenResponse{
		Success: true,
		Id:      socialId,
	}, nil
}
