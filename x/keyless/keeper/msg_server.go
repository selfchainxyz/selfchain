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

// CreateWallet creates a new wallet
func (k msgServer) CreateWallet(goCtx context.Context, msg *types.MsgCreateWallet) (*types.MsgCreateWalletResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the wallet already exists
	_, found := k.getWallet(ctx, msg.WalletAddress)
	if found {
		return nil, fmt.Errorf("wallet already exists: %s", msg.WalletAddress)
	}

	// Create wallet in store
	wallet := types.Wallet{
		Creator:       msg.Creator,
		PubKey:        msg.PubKey,
		WalletAddress: msg.WalletAddress,
		ChainId:       msg.ChainId,
	}

	k.setWallet(ctx, wallet)

	return &types.MsgCreateWalletResponse{
		WalletAddress: msg.WalletAddress,
	}, nil
}

// SignTransaction signs a transaction using the wallet's private key
func (k msgServer) SignTransaction(goCtx context.Context, msg *types.MsgSignTransaction) (*types.MsgSignTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the wallet exists and validate owner
	err := k.ValidateWalletOwner(ctx, msg.WalletAddress, msg.Creator)
	if err != nil {
		return nil, err
	}

	// TODO: Implement actual transaction signing logic here
	signedTx := msg.UnsignedTx // Placeholder for actual signing

	return &types.MsgSignTransactionResponse{
		SignedTx: signedTx,
	}, nil
}

// RecoverWallet recovers a wallet using recovery proof
func (k msgServer) RecoverWallet(goCtx context.Context, msg *types.MsgRecoverWallet) (*types.MsgRecoverWalletResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the wallet exists
	wallet, found := k.getWallet(ctx, msg.WalletAddress)
	if !found {
		return nil, fmt.Errorf("wallet not found: %s", msg.WalletAddress)
	}

	// TODO: Verify recovery proof
	// This should be implemented based on your recovery mechanism

	// Update the wallet with new public key
	wallet.PubKey = msg.NewPubKey
	k.setWallet(ctx, wallet)

	return &types.MsgRecoverWalletResponse{
		WalletAddress: msg.WalletAddress,
	}, nil
}
