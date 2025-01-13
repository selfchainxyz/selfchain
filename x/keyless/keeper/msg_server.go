package keeper

import (
    "context"
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "selfchain/x/keyless/types"
    "selfchain/x/keyless/tss"
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

// GenerateKey implements types.MsgServer
func (k msgServer) GenerateKey(goCtx context.Context, msg *types.MsgGenerateKey) (*types.MsgGenerateKeyResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Get the wallet
    wallet, err := k.Keeper.GetWalletState(ctx, msg.WalletAddress)
    if err != nil {
        return nil, fmt.Errorf("failed to get wallet: %w", err)
    }

    // Verify the creator is the wallet creator
    if wallet.Creator != msg.Creator {
        return nil, types.ErrUnauthorized
    }

    // Generate TSS key shares
    result, err := tss.GenerateKey(goCtx)
    if err != nil {
        return nil, fmt.Errorf("failed to generate key shares: %w", err)
    }

    // Update wallet with remote share and public key
    wallet.RemoteShare = string(result.RemoteShare)
    // TODO: Store public key in a proper format

    // Store updated wallet
    k.Keeper.SetWallet(ctx, wallet)

    return &types.MsgGenerateKeyResponse{
        PersonalShare: result.PersonalShare,
        PublicKey:    result.PublicKey,
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

    // Get the remote share from wallet
    remoteShare := []byte(wallet.RemoteShare)
    if len(remoteShare) == 0 {
        return nil, fmt.Errorf("wallet has no remote share")
    }

    // Sign the transaction using TSS
    result, err := tss.SignMessage(goCtx, msg.TransactionData, msg.PersonalShare, remoteShare)
    if err != nil {
        return nil, fmt.Errorf("failed to sign transaction: %w", err)
    }

    // Combine R and S into a single signature byte slice
    signature := append(result.R.Bytes(), result.S.Bytes()...)

    return &types.MsgSignTransactionResponse{
        Signature: signature,
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
