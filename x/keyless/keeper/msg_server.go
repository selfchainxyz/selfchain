package keeper

import (
	"context"
	"fmt"
	"crypto/sha256"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
	"selfchain/x/keyless/crypto/signing"
	"selfchain/x/keyless/networks"
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

	// Get wallet
	wallet, found := k.getWallet(ctx, msg.WalletAddress)
	if !found {
		return nil, fmt.Errorf("wallet not found: %s", msg.WalletAddress)
	}

	// Create signer factory
	signerFactory := signing.NewSignerFactory(networks.NewNetworkRegistry())

	// Convert unsigned transaction to bytes
	unsignedTx := []byte(msg.UnsignedTx)

	// Sign the transaction
	signature, err := signerFactory.Sign(ctx, wallet.ChainId, unsignedTx, map[string]interface{}{
		"wallet_address": msg.WalletAddress,
		"public_key":    wallet.PubKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return &types.MsgSignTransactionResponse{
		SignedTx: string(signature),
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

	// Verify recovery proof
	// 1. Hash the recovery data
	recoveryHash := sha256.Sum256([]byte(msg.RecoveryProof))
	
	// 2. Create signer factory for verification
	signerFactory := signing.NewSignerFactory(networks.NewNetworkRegistry())
	
	// 3. Verify the recovery proof signature
	pubKeyBytes := []byte(wallet.PubKey)
	if msg.Signature == "" {
		return nil, fmt.Errorf("recovery signature is required")
	}
	signatureBytes := []byte(msg.Signature)
	
	valid, err := signerFactory.Verify(wallet.ChainId, pubKeyBytes, recoveryHash[:], signatureBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to verify recovery proof: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid recovery proof")
	}

	// Update the wallet with new public key
	wallet.PubKey = msg.NewPubKey
	k.setWallet(ctx, wallet)

	return &types.MsgRecoverWalletResponse{
		WalletAddress: msg.WalletAddress,
	}, nil
}
