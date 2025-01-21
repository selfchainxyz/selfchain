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

	// Generate a unique wallet ID
	walletId := fmt.Sprintf("%s-%s", msg.Creator, msg.WalletAddress)

	// Check if wallet already exists
	_, found := k.GetWalletFromStore(ctx, walletId)
	if found {
		return nil, fmt.Errorf("wallet already exists: %s", walletId)
	}

	// Create a new wallet
	wallet := types.NewWallet(
		walletId,
		msg.PubKey,
		1, // Initial security level
		2, // Initial threshold
		3, // Initial parties
	)

	// Add creator permission
	wallet.Permissions = append(wallet.Permissions, "owner")

	k.SetWallet(ctx, wallet)

	return &types.MsgCreateWalletResponse{
		WalletAddress: msg.WalletAddress,
	}, nil
}

// RecoverWallet recovers a wallet using recovery proof
func (k msgServer) RecoverWallet(goCtx context.Context, msg *types.MsgRecoverWallet) (*types.MsgRecoverWalletResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Generate wallet ID
	walletId := fmt.Sprintf("%s-%s", msg.Creator, msg.WalletAddress)

	// Validate recovery proof
	if err := k.ValidateRecoveryProof(ctx, walletId, msg.RecoveryProof); err != nil {
		return nil, err
	}

	// Get the existing wallet
	wallet, found := k.GetWalletFromStore(ctx, walletId)
	if !found {
		return nil, fmt.Errorf("wallet not found: %s", walletId)
	}

	// Update the wallet with new public key
	wallet.PublicKey = msg.NewPubKey
	wallet.KeyVersion++
	k.SetWallet(ctx, &wallet)

	return &types.MsgRecoverWalletResponse{
		WalletAddress: msg.WalletAddress,
	}, nil
}

// SignTransaction signs a transaction using the keyless wallet
func (k msgServer) SignTransaction(goCtx context.Context, msg *types.MsgSignTransaction) (*types.MsgSignTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Generate wallet ID
	walletId := fmt.Sprintf("%s-%s", msg.Creator, msg.WalletAddress)

	// Get the wallet
	wallet, found := k.GetWalletFromStore(ctx, walletId)
	if !found {
		return nil, fmt.Errorf("wallet not found: %s", walletId)
	}

	// Verify permissions
	if err := k.ValidateWalletAccess(ctx, walletId, "sign"); err != nil {
		return nil, err
	}

	// Sign the transaction using TSS
	signedTx, err := k.SignWithTSS(ctx, &wallet, msg.UnsignedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	return &types.MsgSignTransactionResponse{
		SignedTx: string(signedTx),
	}, nil
}

// RotateKey initiates key rotation for a wallet
func (k msgServer) RotateKey(goCtx context.Context, msg *types.MsgRotateKey) (*types.MsgRotateKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the wallet
	wallet, found := k.GetWalletFromStore(ctx, msg.WalletId)
	if !found {
		return nil, fmt.Errorf("wallet not found: %s", msg.WalletId)
	}

	// Verify permissions
	if err := k.ValidateWalletAccess(ctx, msg.WalletId, "rotate"); err != nil {
		return nil, err
	}

	// Increment key version
	newVersion := wallet.KeyVersion + 1

	return &types.MsgRotateKeyResponse{
		WalletId:   msg.WalletId,
		NewVersion: newVersion,
	}, nil
}

// BatchSign performs batch signing operation
func (k msgServer) BatchSign(goCtx context.Context, msg *types.MsgBatchSign) (*types.MsgBatchSignResponse, error) {
	// TODO: Implement batch signing
	return nil, fmt.Errorf("not implemented")
}

// InitiateKeyRotation initiates key rotation for a wallet
func (k msgServer) InitiateKeyRotation(goCtx context.Context, msg *types.MsgInitiateKeyRotation) (*types.MsgInitiateKeyRotationResponse, error) {
	// TODO: Implement key rotation initiation
	return nil, fmt.Errorf("not implemented")
}

// CompleteKeyRotation completes key rotation for a wallet
func (k msgServer) CompleteKeyRotation(goCtx context.Context, msg *types.MsgCompleteKeyRotation) (*types.MsgCompleteKeyRotationResponse, error) {
	// TODO: Implement key rotation completion
	return nil, fmt.Errorf("not implemented")
}

// CancelKeyRotation cancels key rotation for a wallet
func (k msgServer) CancelKeyRotation(goCtx context.Context, msg *types.MsgCancelKeyRotation) (*types.MsgCancelKeyRotationResponse, error) {
	// TODO: Implement key rotation cancellation
	return nil, fmt.Errorf("not implemented")
}
