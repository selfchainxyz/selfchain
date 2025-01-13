package keeper

import (
    "context"
    "fmt"

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
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Create the wallet using keeper method
    wallet, err := k.Keeper.CreateWallet(ctx, msg.Creator, msg.Did)
    if err != nil {
        return nil, fmt.Errorf("failed to create wallet: %w", err)
    }

    return &types.MsgCreateWalletResponse{
        Address: wallet.Address,
    }, nil
}

// SignTransaction implements types.MsgServer
func (k msgServer) SignTransaction(goCtx context.Context, msg *types.MsgSignTransaction) (*types.MsgSignTransactionResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Get the wallet
    wallet, err := k.Keeper.GetWalletState(ctx, msg.WalletAddress)
    if err != nil {
        return nil, fmt.Errorf("failed to get wallet: %w", err)
    }

    // Verify the signer is the wallet creator
    if wallet.Creator != msg.Creator {
        return nil, types.ErrUnauthorized
    }

    // TODO: Implement MPC-TSS signing logic
    // For now, return empty signature
    return &types.MsgSignTransactionResponse{
        Signature: []byte{},
    }, nil
}

// RecoverWallet implements types.MsgServer
func (k msgServer) RecoverWallet(goCtx context.Context, msg *types.MsgRecoverWallet) (*types.MsgRecoverWalletResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Get the wallet by DID
    wallet, err := k.Keeper.GetWalletStateByDID(ctx, msg.Did)
    if err != nil {
        return nil, fmt.Errorf("failed to get wallet: %w", err)
    }

    // TODO: Implement wallet recovery logic using DID verification
    // For now, just verify the creator
    if wallet.Creator != msg.Creator {
        return nil, types.ErrUnauthorized
    }

    return &types.MsgRecoverWalletResponse{
        Address: wallet.Address,
    }, nil
}
