package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

// CreateRecoverySession creates a new recovery session for a wallet
func (k Keeper) CreateRecoverySession(ctx sdk.Context, creator, walletAddress string) error {
	// Validate the wallet exists
	wallet, err := k.GetWallet(ctx, walletAddress)
	if err != nil {
		return fmt.Errorf("wallet not found: %s", walletAddress)
	}

	// Check if recovery is already in progress
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("recovery/%s", walletAddress))
	if store.Has(key) {
		return types.ErrRecoveryInProgress
	}

	// Verify creator has permission to initiate recovery
	if err := k.ValidateWalletAccess(ctx, walletAddress, creator, "recovery"); err != nil {
		return err
	}

	// Create recovery info
	now := ctx.BlockTime()
	recoveryInfo := &types.RecoveryInfo{
		Did:             wallet.Creator,
		RecoveryAddress: creator,
		Status:          types.RecoveryStatus_RECOVERY_STATUS_PENDING,
		CreatedAt:       now,
	}

	// Generate recovery token using identity module
	token, err := k.identityKeeper.GenerateRecoveryToken(ctx, wallet.Creator)
	if err != nil {
		return fmt.Errorf("failed to generate recovery token: %v", err)
	}
	recoveryInfo.RecoveryToken = token

	// Store recovery info
	bz := k.cdc.MustMarshal(recoveryInfo)
	store.Set(key, bz)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRecoveryStarted,
			sdk.NewAttribute(types.AttributeKeyWalletAddress, walletAddress),
			sdk.NewAttribute(types.AttributeKeyRecoveryAddress, creator),
			sdk.NewAttribute(types.AttributeKeyTimestamp, now.String()),
		),
	)

	return nil
}

// ValidateRecoverySession validates a recovery session
func (k Keeper) ValidateRecoverySession(ctx sdk.Context, creator, walletAddress string) error {
	// Get recovery info
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("recovery/%s", walletAddress))
	bz := store.Get(key)
	if bz == nil {
		return types.ErrRecoveryNotAllowed
	}

	var recoveryInfo types.RecoveryInfo
	k.cdc.MustUnmarshal(bz, &recoveryInfo)

	// Check recovery status
	if recoveryInfo.Status != types.RecoveryStatus_RECOVERY_STATUS_PENDING {
		return types.ErrRecoveryNotAllowed
	}

	// Verify creator matches recovery address
	if recoveryInfo.RecoveryAddress != creator {
		return types.ErrUnauthorized
	}

	// Check timelock (24 hours)
	timeLock := time.Duration(24) * time.Hour
	if ctx.BlockTime().Sub(recoveryInfo.CreatedAt) < timeLock {
		return types.ErrRecoveryNotAllowed
	}

	return nil
}

// verifyRecoveryProof verifies the recovery proof for a wallet
func (k Keeper) verifyRecoveryProof(ctx sdk.Context, wallet *types.Wallet, proof string) bool {
	// Get recovery info
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("recovery/%s", wallet.WalletAddress))
	bz := store.Get(key)
	if bz == nil {
		return false
	}

	var recoveryInfo types.RecoveryInfo
	k.cdc.MustUnmarshal(bz, &recoveryInfo)

	// Verify recovery token
	if err := k.identityKeeper.VerifyRecoveryToken(ctx, recoveryInfo.Did, recoveryInfo.RecoveryToken); err != nil {
		return false
	}

	// Verify identity proof
	if err := k.identityKeeper.VerifyDIDOwnership(ctx, recoveryInfo.Did, sdk.AccAddress(recoveryInfo.RecoveryAddress)); err != nil {
		return false
	}

	// Verify MFA if enabled
	if err := k.identityKeeper.VerifyMFA(ctx, recoveryInfo.Did); err != nil {
		return false
	}

	// Verify OAuth2 token if provided
	if err := k.identityKeeper.VerifyOAuth2Token(ctx, recoveryInfo.Did, proof); err != nil {
		return false
	}

	return true
}

// RecoverWallet recovers a wallet by its address
func (k Keeper) RecoverWallet(ctx sdk.Context, msg *types.MsgRecoverWallet) error {
	// Get the wallet
	wallet, err := k.GetWallet(ctx, msg.WalletAddress)
	if err != nil {
		return fmt.Errorf("wallet not found: %v", err)
	}

	// Validate recovery session
	if err := k.ValidateRecoverySession(ctx, msg.Creator, msg.WalletAddress); err != nil {
		return err
	}

	// Verify recovery proof
	if !k.verifyRecoveryProof(ctx, wallet, msg.RecoveryProof) {
		return types.ErrInvalidRecoveryProof
	}

	// Delete the old wallet
	k.DeleteWallet(ctx, msg.WalletAddress)

	// Create new wallet with updated owner and public key
	now := ctx.BlockTime()
	newWallet := types.NewWallet(
		msg.Creator,
		msg.NewPubKey,
		msg.WalletAddress,
		wallet.ChainId,
		types.WalletStatus_WALLET_STATUS_ACTIVE,
		wallet.KeyVersion+1,
	)
	newWallet.CreatedAt = &now
	newWallet.UpdatedAt = &now
	newWallet.LastUsed = &now

	// Save new wallet
	if err := k.SaveWallet(ctx, newWallet); err != nil {
		return fmt.Errorf("failed to save wallet: %v", err)
	}

	// Delete recovery info
	if err := k.DeleteRecoveryInfo(ctx, msg.WalletAddress); err != nil {
		return fmt.Errorf("failed to delete recovery info: %v", err)
	}

	// Emit recovery completed event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRecoverWallet,
			sdk.NewAttribute(types.AttributeKeyWalletID, newWallet.Id),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyWalletAddress, newWallet.WalletAddress),
		),
	)

	return nil
}

// GetRecoveryInfo gets recovery info for a wallet
func (k Keeper) GetRecoveryInfo(ctx sdk.Context, walletId string) (*types.RecoveryInfo, error) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("recovery/%s", walletId))
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("recovery info not found")
	}

	var info types.RecoveryInfo
	k.cdc.MustUnmarshal(bz, &info)
	return &info, nil
}

// SaveRecoveryInfo saves recovery info for a wallet
func (k Keeper) SaveRecoveryInfo(ctx sdk.Context, info *types.RecoveryInfo) error {
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("recovery/%s", info.Did))
	bz := k.cdc.MustMarshal(info)
	store.Set(key, bz)
	return nil
}

// DeleteRecoveryInfo deletes recovery info for a wallet
func (k Keeper) DeleteRecoveryInfo(ctx sdk.Context, walletId string) error {
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("recovery/%s", walletId))
	store.Delete(key)
	return nil
}
