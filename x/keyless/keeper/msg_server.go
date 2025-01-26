package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// Create wallet
	now := ctx.BlockTime()
	wallet := types.NewWallet(
		msg.Creator,
		msg.PubKey,
		msg.WalletAddress,
		msg.ChainId,
		types.WalletStatus_WALLET_STATUS_ACTIVE,
		0, // Initial key version
	)
	wallet.CreatedAt = &now
	wallet.UpdatedAt = &now
	wallet.LastUsed = &now

	// Save wallet to store
	err := k.SaveWallet(ctx, wallet)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.MsgCreateWalletResponse{
		WalletAddress: msg.WalletAddress,
	}, nil
}

// SignTransaction signs a transaction using TSS
func (k msgServer) SignTransaction(goCtx context.Context, msg *types.MsgSignTransaction) (*types.MsgSignTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate creator is authorized
	authorized, err := k.IsWalletAuthorized(ctx, msg.Creator, msg.WalletAddress)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if !authorized {
		return nil, status.Error(codes.PermissionDenied, "not authorized")
	}

	// TODO: Implement TSS signing logic
	signedTx := "dummy_signed_tx"

	return &types.MsgSignTransactionResponse{
		SignedTx: signedTx,
	}, nil
}

// BatchSign initiates batch signing process
func (k msgServer) BatchSign(goCtx context.Context, msg *types.MsgBatchSignRequest) (*types.MsgBatchSignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate creator is authorized
	authorized, err := k.IsWalletAuthorized(ctx, msg.Creator, msg.WalletAddress)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if !authorized {
		return nil, status.Error(codes.PermissionDenied, "not authorized")
	}

	// Get wallet
	wallet, err := k.GetWallet(ctx, msg.WalletAddress)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Update wallet status
	wallet.Status = types.WalletStatus_WALLET_STATUS_ROTATING
	err = k.SaveWallet(ctx, wallet)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Create batch sign status
	batchStatus := &types.BatchSignStatusInfo{
		WalletAddress: msg.WalletAddress,
		Messages:      msg.Messages,
		Status:        types.BatchSignStatus_BATCH_SIGN_STATUS_IN_PROGRESS,
		Signatures:    make([]string, 0),
	}

	// Save batch sign status
	err = k.SetBatchSignStatus(ctx, batchStatus)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.MsgBatchSignResponse{}, nil
}

// CompleteKeyRotation completes key rotation process
func (k msgServer) CompleteKeyRotation(goCtx context.Context, msg *types.MsgCompleteKeyRotation) (*types.MsgCompleteKeyRotationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate creator is authorized
	authorized, err := k.IsWalletAuthorized(ctx, msg.Creator, msg.WalletAddress)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if !authorized {
		return nil, status.Error(codes.PermissionDenied, "not authorized")
	}

	// Get wallet
	wallet, err := k.GetWallet(ctx, msg.WalletAddress)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Update wallet with new public key
	wallet.PublicKey = msg.NewPubKey
	wallet.Status = types.WalletStatus_WALLET_STATUS_ACTIVE
	err = k.SaveWallet(ctx, wallet)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.MsgCompleteKeyRotationResponse{
		WalletAddress: msg.WalletAddress,
		Version:       msg.Version,
	}, nil
}

// InitiateKeyRotation initiates key rotation process
func (k msgServer) InitiateKeyRotation(goCtx context.Context, msg *types.MsgInitiateKeyRotation) (*types.MsgInitiateKeyRotationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate creator is authorized
	authorized, err := k.IsWalletAuthorized(ctx, msg.Creator, msg.WalletAddress)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if !authorized {
		return nil, status.Error(codes.PermissionDenied, "not authorized")
	}

	// Get wallet
	wallet, err := k.GetWallet(ctx, msg.WalletAddress)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Update wallet status
	wallet.Status = types.WalletStatus_WALLET_STATUS_ROTATING
	wallet.KeyVersion++
	err = k.SaveWallet(ctx, wallet)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.MsgInitiateKeyRotationResponse{
		WalletAddress: msg.WalletAddress,
		NewVersion:    wallet.KeyVersion,
	}, nil
}

// RecoverWallet recovers a wallet
func (k msgServer) RecoverWallet(goCtx context.Context, msg *types.MsgRecoverWallet) (*types.MsgRecoverWalletResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic message fields
	if err := msg.ValidateBasic(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Recover the wallet
	err := k.Keeper.RecoverWallet(ctx, msg.WalletAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.MsgRecoverWalletResponse{
		WalletAddress: msg.WalletAddress,
	}, nil
}
