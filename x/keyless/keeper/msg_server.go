package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	selfchainTss "selfchain/x/keyless/tss"
	"selfchain/x/keyless/types"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	// Get network configuration
	networkRegistry := networks.DefaultRegistry()
	networkConfig, err := networkRegistry.GetNetwork(msg.ChainID)
	if err != nil {
		return nil, fmt.Errorf("unsupported chain ID: %w", err)
	}

	// Generate encryption key for personal share
	encryptionKey, err := crypto.NewEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}

	// Generate pre-parameters for TSS
	preParams, err := keygen.GeneratePreParams(time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate pre-parameters: %w", err)
	}

	// Generate TSS key shares
	result, err := selfchainTss.GenerateKey(goCtx, preParams, msg.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key shares: %w", err)
	}

	// Encrypt party2's data for storage
	remoteShare, err := selfchainTss.EncryptShare(encryptionKey, result.Party2Data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt remote share: %w", err)
	}

	// Encrypt party1's data for response
	personalShare, err := selfchainTss.EncryptShare(encryptionKey, result.Party1Data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt personal share: %w", err)
	}

	// Get the public key from either party (they should be the same)
	publicKeyBytes, err := json.Marshal(result.Party1Data.ECDSAPub)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize public key: %w", err)
	}

	// Update wallet with remote share
	wallet.RemoteShare = remoteShare.EncryptedData
	wallet.ChainID = msg.ChainID
	
	// Store updated wallet
	k.Keeper.SetWallet(ctx, wallet)

	return &types.MsgGenerateKeyResponse{
		PersonalShare: []byte(personalShare.EncryptedData),
		PublicKey:    publicKeyBytes,
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

	// Get network configuration
	networkRegistry := networks.DefaultRegistry()
	networkConfig, err := networkRegistry.GetNetwork(wallet.ChainID)
	if err != nil {
		return nil, fmt.Errorf("unsupported chain ID: %w", err)
	}

	// Generate encryption key (in production this should be securely stored/retrieved)
	encryptionKey, err := crypto.NewEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}

	// Get the remote share from wallet and decrypt it
	remoteShare := &selfchainTss.EncryptedShare{
		EncryptedData: wallet.RemoteShare,
		ChainID:      wallet.ChainID,
	}
	remoteShareData, err := selfchainTss.DecryptShare(encryptionKey, remoteShare)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt remote share: %w", err)
	}

	// Decrypt personal share
	personalShare := &selfchainTss.EncryptedShare{
		EncryptedData: string(msg.PersonalShare),
		ChainID:      wallet.ChainID,
	}
	personalShareData, err := selfchainTss.DecryptShare(encryptionKey, personalShare)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt personal share: %w", err)
	}

	// Sign the transaction using TSS with appropriate algorithm
	result, err := selfchainTss.SignMessage(goCtx, msg.TransactionData, personalShareData, remoteShareData)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Format signature based on network requirements
	var signature []byte
	switch networkConfig.Algorithm {
	case networks.ECDSA:
		signature = append(result.R.Bytes(), result.S.Bytes()...)
	case networks.EdDSA:
		// Add EdDSA signature formatting
		return nil, fmt.Errorf("EdDSA signing not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported signing algorithm: %s", networkConfig.Algorithm)
	}

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

	// Verify the creator is the original wallet creator
	if wallet.Creator != msg.Creator {
		return nil, types.ErrUnauthorized
	}

	return &types.MsgRecoverWalletResponse{
		Address: wallet.Address,
	}, nil
}
