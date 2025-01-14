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

// CreateWallet creates a new keyless wallet
func (k msgServer) CreateWallet(goCtx context.Context, msg *types.MsgCreateWallet) (*types.MsgCreateWalletResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if wallet already exists
	_, found := k.GetWalletFromStore(ctx, msg.WalletAddress)
	if found {
		return nil, fmt.Errorf("wallet already exists: %s", msg.WalletAddress)
	}

	wallet := types.NewWallet(
		msg.Creator,
		msg.WalletAddress,
		msg.PubKey,
		msg.ChainId,
	)

	k.SetWallet(ctx, wallet)

	return &types.MsgCreateWalletResponse{
		WalletAddress: msg.WalletAddress,
	}, nil
}

// RecoverWallet recovers a wallet using recovery proof
func (k msgServer) RecoverWallet(goCtx context.Context, msg *types.MsgRecoverWallet) (*types.MsgRecoverWalletResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate recovery proof and signature
	if err := k.ValidateRecoveryProof(ctx, msg.WalletAddress, msg.RecoveryProof, msg.Signature); err != nil {
		return nil, err
	}

	// Get the existing wallet
	wallet, found := k.GetWalletFromStore(ctx, msg.WalletAddress)
	if !found {
		return nil, fmt.Errorf("wallet not found: %s", msg.WalletAddress)
	}

	// Update the wallet with new public key
	wallet.PubKey = msg.NewPubKey
	k.SetWallet(ctx, &wallet)

	return &types.MsgRecoverWalletResponse{
		WalletAddress: msg.WalletAddress,
	}, nil
}

// SignTransaction signs a transaction using the keyless wallet
func (k msgServer) SignTransaction(goCtx context.Context, msg *types.MsgSignTransaction) (*types.MsgSignTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the wallet
	wallet, found := k.GetWalletFromStore(ctx, msg.WalletAddress)
	if !found {
		return nil, fmt.Errorf("wallet not found: %s", msg.WalletAddress)
	}

	// Verify chain ID
	if wallet.ChainId != msg.ChainId {
		return nil, fmt.Errorf("chain ID mismatch: expected %s, got %s", wallet.ChainId, msg.ChainId)
	}

	// Sign the transaction using TSS
	signedTx, err := k.SignWithTSS(ctx, &wallet, msg.UnsignedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	return &types.MsgSignTransactionResponse{
		SignedTx: signedTx,
	}, nil
}
