package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateWallet implements types.MsgServer
func (k msgServer) CreateWallet(goCtx context.Context, msg *types.MsgCreateWallet) (*types.MsgCreateWalletResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	// TODO: Implement wallet creation logic
	// This will be implemented in the next step

	return &types.MsgCreateWalletResponse{}, nil
}

// SignTransaction implements types.MsgServer
func (k msgServer) SignTransaction(goCtx context.Context, msg *types.MsgSignTransaction) (*types.MsgSignTransactionResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	// TODO: Implement transaction signing logic
	// This will be implemented in the next step

	return &types.MsgSignTransactionResponse{}, nil
}

// RecoverWallet implements types.MsgServer
func (k msgServer) RecoverWallet(goCtx context.Context, msg *types.MsgRecoverWallet) (*types.MsgRecoverWalletResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	// TODO: Implement wallet recovery logic
	// This will be implemented in the next step

	return &types.MsgRecoverWalletResponse{}, nil
}
