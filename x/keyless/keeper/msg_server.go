package keeper

import (
	"context"
	"encoding/json"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"selfchain/x/keyless/types"
)

type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
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
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to save wallet: %v", err)
	}

	return &types.MsgCreateWalletResponse{
		WalletAddress: msg.WalletAddress,
	}, nil
}

// SignTransaction handles MsgSignTransaction
func (k msgServer) SignTransaction(goCtx context.Context, msg *types.MsgSignTransaction) (*types.MsgSignTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get wallet
	wallet, found := k.Keeper.GetWallet(ctx, msg.WalletAddress)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrWalletNotFound, "wallet %s not found", msg.WalletAddress)
	}

	// Verify creator has permission
	if wallet.Creator != msg.Creator {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "creator %s not authorized", msg.Creator)
	}

	// Verify wallet is active
	if wallet.Status != types.WalletStatus_WALLET_STATUS_ACTIVE {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "wallet %s is not active", msg.WalletAddress)
	}

	// Get key shares
	personalShare, found := k.Keeper.GetKeyShare(ctx, msg.Creator)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "personal share not found for creator %s", msg.Creator)
	}

	remoteShare, found := k.Keeper.GetKeyShare(ctx, msg.WalletAddress)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "remote share not found for wallet %s", msg.WalletAddress)
	}

	// Unmarshal key shares
	var personalPartyData keygen.LocalPartySaveData
	if err := json.Unmarshal(personalShare, &personalPartyData); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to unmarshal personal share: %v", err)
	}

	var remotePartyData keygen.LocalPartySaveData
	if err := json.Unmarshal(remoteShare, &remotePartyData); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to unmarshal remote share: %v", err)
	}

	// Sign transaction using TSS protocol
	protocol := k.GetTSSProtocol()
	if protocol == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "TSS protocol not configured")
	}

	signResult, err := protocol.SignMessage(ctx, []byte(msg.UnsignedTx), personalShare, remoteShare)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidSignature, "failed to sign transaction: %v", err)
	}

	// Return the signature
	return &types.MsgSignTransactionResponse{
		SignedTx: string(signResult),
	}, nil
}

// BatchSign initiates batch signing process
func (k msgServer) BatchSign(goCtx context.Context, msg *types.MsgBatchSignRequest) (*types.MsgBatchSignResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate creator is authorized
	authorized := k.IsWalletAuthorized(ctx, msg.WalletAddress, msg.Creator, types.WalletPermission_WALLET_PERMISSION_SIGN)
	if !authorized {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "creator %s not authorized", msg.Creator)
	}

	// Get wallet
	wallet, found := k.GetWallet(ctx, msg.WalletAddress)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrWalletNotFound, "wallet %s not found", msg.WalletAddress)
	}

	// Update wallet status
	wallet.Status = types.WalletStatus_WALLET_STATUS_ROTATING
	k.SetWallet(ctx, wallet)

	// Create batch sign status
	status := &types.BatchSignStatusInfo{
		WalletAddress: msg.WalletAddress,
		Messages:     msg.Messages,
		Status:       types.BatchSignStatus_BATCH_SIGN_STATUS_IN_PROGRESS,
	}

	return &types.MsgBatchSignResponse{
		Status: status,
	}, nil
}

// CompleteKeyRotation completes key rotation process
func (k msgServer) CompleteKeyRotation(goCtx context.Context, msg *types.MsgCompleteKeyRotation) (*types.MsgCompleteKeyRotationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate creator is authorized
	authorized := k.IsWalletAuthorized(ctx, msg.WalletAddress, msg.Creator, types.WalletPermission_WALLET_PERMISSION_ROTATE)
	if !authorized {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "creator %s not authorized", msg.Creator)
	}

	// Get wallet
	wallet, found := k.GetWallet(ctx, msg.WalletAddress)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrWalletNotFound, "wallet %s not found", msg.WalletAddress)
	}

	// Update wallet with new public key
	wallet.PublicKey = msg.NewPubKey
	wallet.Status = types.WalletStatus_WALLET_STATUS_ACTIVE
	k.SetWallet(ctx, wallet)

	return &types.MsgCompleteKeyRotationResponse{
		WalletAddress: msg.WalletAddress,
		Version:       msg.Version,
	}, nil
}

// InitiateKeyRotation initiates key rotation process
func (k msgServer) InitiateKeyRotation(goCtx context.Context, msg *types.MsgInitiateKeyRotation) (*types.MsgInitiateKeyRotationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate creator is authorized
	authorized := k.IsWalletAuthorized(ctx, msg.WalletAddress, msg.Creator, types.WalletPermission_WALLET_PERMISSION_ROTATE)
	if !authorized {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "creator %s not authorized", msg.Creator)
	}

	// Get wallet
	wallet, found := k.GetWallet(ctx, msg.WalletAddress)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrWalletNotFound, "wallet %s not found", msg.WalletAddress)
	}

	// Update wallet status
	wallet.Status = types.WalletStatus_WALLET_STATUS_ROTATING
	wallet.KeyVersion++
	k.SetWallet(ctx, wallet)

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
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid request: %v", err)
	}

	// Recover the wallet
	err := k.Keeper.RecoverWallet(ctx, msg)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to recover wallet: %v", err)
	}

	return &types.MsgRecoverWalletResponse{
		WalletAddress: msg.WalletAddress,
	}, nil
}

func (k msgServer) GetTSSProtocol() *mockTSSProtocol {
	return &mockTSSProtocol{}
}

type mockTSSProtocol struct{}

func (m *mockTSSProtocol) SignMessage(ctx context.Context, message []byte, personalShare, remoteShare []byte) ([]byte, error) {
	// Mock signing logic
	return []byte("mocked_signature"), nil
}
