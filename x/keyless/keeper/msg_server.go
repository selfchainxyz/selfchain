package keeper

import (
	"context"
	"fmt"
	"crypto/sha256"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"selfchain/x/keyless/types"
	"selfchain/x/keyless/crypto/signing"
	"selfchain/x/keyless/networks"
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

// getTSSPartyData retrieves TSS party data for a wallet
func (k msgServer) getTSSPartyData(ctx sdk.Context, wallet types.Wallet) (*keygen.LocalPartySaveData, *keygen.LocalPartySaveData, error) {
	// TODO: Implement TSS party data retrieval from wallet
	// This should retrieve the TSS key shares from secure storage
	// For now, we'll return an error
	return nil, nil, fmt.Errorf("TSS party data retrieval not implemented")
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

	// Get TSS party data from wallet
	party1Data, party2Data, err := k.getTSSPartyData(ctx, wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to get TSS party data: %w", err)
	}

	// Sign using TSS
	signResult, err := tss.SignMessage(ctx, unsignedTx, party1Data, party2Data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign with TSS: %w", err)
	}

	// Format signature according to network
	signature, err := signerFactory.Sign(ctx, wallet.ChainId, unsignedTx, map[string]interface{}{
		"wallet_address": msg.WalletAddress,
		"public_key":    wallet.PubKey,
	}, signResult)
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
