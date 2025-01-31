package keeper

import (
	"context"

	"selfchain/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) VerifyOAuthToken(goCtx context.Context, msg *types.MsgVerifyOAuthToken) (*types.MsgVerifyOAuthTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get user info from token
	userInfo, err := k.Keeper.GetUserInfo(ctx, msg.Provider, msg.Token)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get user info from token")
	}

	// Get social identity by social ID
	socialIdentity, found := k.Keeper.GetSocialIdentityBySocialID(ctx, msg.Provider, userInfo.Id)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrSocialIdentityNotFound, "social identity not found")
	}

	return &types.MsgVerifyOAuthTokenResponse{
		Success: true,
		Id:      socialIdentity.Id,
	}, nil
}
