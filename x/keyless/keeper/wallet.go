package keeper

import (
	"fmt"

	"selfchain/x/keyless/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAllWalletsFromStore returns all wallets from the KVStore
func (k Keeper) GetAllWalletsFromStore(ctx sdk.Context) ([]types.Wallet, error) {
	store := k.GetWalletStore(ctx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var wallets []types.Wallet
	for ; iterator.Valid(); iterator.Next() {
		var wallet types.Wallet
		err := k.cdc.Unmarshal(iterator.Value(), &wallet)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal wallet: %v", err)
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

// ValidateWalletAccess validates if the permission exists in the wallet's permissions list
func (k Keeper) ValidateWalletAccess(ctx sdk.Context, walletAddress string, permission string) error {
	_, err := k.GetWallet(ctx, walletAddress)
	if err != nil {
		return fmt.Errorf("wallet not found: %s", walletAddress)
	}

	// TODO: Implement permission validation
	return nil
}

// IsWalletAuthorized checks if the creator is authorized to operate on the wallet
func (k Keeper) IsWalletAuthorized(ctx sdk.Context, creator string, walletAddress string) (bool, error) {
	wallet, err := k.GetWallet(ctx, walletAddress)
	if err != nil {
		return false, err
	}
	return wallet.Creator == creator, nil
}

// ValidateRecoveryProof validates a recovery proof for a wallet
func (k Keeper) ValidateRecoveryProof(ctx sdk.Context, walletId, recoveryProof string) error {
	// TODO: Implement recovery proof validation
	// This should integrate with the identity module for verification
	return nil
}

// SignWithTSS signs a transaction using TSS
func (k Keeper) SignWithTSS(ctx sdk.Context, wallet *types.Wallet, unsignedTx string) ([]byte, error) {
	if wallet == nil {
		return nil, fmt.Errorf("wallet cannot be nil")
	}

	// Get party data for the wallet
	_, err := k.GetPartyData(ctx, wallet.WalletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get party data: %w", err)
	}

	// Verify wallet is active
	if wallet.Status != types.WalletStatus_WALLET_STATUS_ACTIVE {
		return nil, fmt.Errorf("wallet is not active")
	}

	// Create signing session ID
	sessionID := fmt.Sprintf("%s-%d", wallet.WalletAddress, ctx.BlockHeight())

	// Create signing session
	session := &types.SigningSession{
		SessionId: sessionID,
		WalletId:  wallet.WalletAddress,
		Message:   []byte(unsignedTx),
		Status:    types.SigningStatus_SIGNING_STATUS_IN_PROGRESS,
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}

	// Save signing session
	if err := k.SaveSigningSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save signing session: %w", err)
	}

	// Create dummy signature for testing
	// TODO: Implement actual TSS signing protocol
	signature := []byte("dummy_signature")

	// Update session status
	session.Status = types.SigningStatus_SIGNING_STATUS_COMPLETED
	session.UpdatedAt = ctx.BlockTime()
	if err := k.SaveSigningSession(ctx, session); err != nil {
		ctx.Logger().Error("failed to update signing session status", "error", err)
	}

	// Update wallet metadata
	if err := k.updateWalletAfterSigning(ctx, wallet); err != nil {
		// Log error but don't fail the signing
		ctx.Logger().Error("failed to update wallet metadata after signing", "error", err)
	}

	return signature, nil
}

// GetSigningSession retrieves a signing session by ID
func (k Keeper) GetSigningSession(ctx sdk.Context, sessionID string) (*types.SigningSession, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixSigningSession)
	bz := store.Get([]byte(sessionID))
	if bz == nil {
		return nil, fmt.Errorf("signing session not found: %s", sessionID)
	}

	var session types.SigningSession
	if err := k.cdc.Unmarshal(bz, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal signing session: %w", err)
	}

	return &session, nil
}

// SaveSigningSession stores a signing session
func (k Keeper) SaveSigningSession(ctx sdk.Context, session *types.SigningSession) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixSigningSession)
	bz, err := k.cdc.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal signing session: %w", err)
	}
	store.Set([]byte(session.SessionId), bz)
	return nil
}

// updateWalletAfterSigning updates wallet metadata after successful signing
func (k Keeper) updateWalletAfterSigning(ctx sdk.Context, wallet *types.Wallet) error {
	// Update usage count and last used timestamp
	blockTime := ctx.BlockTime()
	wallet.LastUsed = &blockTime
	wallet.UsageCount++

	// Save updated wallet
	return k.SaveWallet(ctx, wallet)
}

// GetWalletStatus returns the current status of a wallet
func (k Keeper) GetWalletStatus(ctx sdk.Context, walletAddress string) (types.WalletStatus, error) {
	wallet, err := k.GetWallet(ctx, walletAddress)
	if err != nil {
		return types.WalletStatus_WALLET_STATUS_UNSPECIFIED, fmt.Errorf("failed to get wallet: %v", err)
	}

	return wallet.Status, nil
}

// SetWalletStatus updates the status of a wallet
func (k Keeper) SetWalletStatus(ctx sdk.Context, walletAddress string, status types.WalletStatus) error {
	wallet, err := k.GetWallet(ctx, walletAddress)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %v", err)
	}

	wallet.Status = status
	err = k.SaveWallet(ctx, wallet)
	if err != nil {
		return fmt.Errorf("failed to save wallet: %v", err)
	}

	return nil
}
