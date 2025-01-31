package keeper

import (
	"fmt"
	"time"

	"selfchain/x/keyless/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateRecoverySession creates a new recovery session for a wallet
func (k Keeper) CreateRecoverySession(ctx sdk.Context, walletAddress string) (*types.RecoveryInfo, error) {
	if walletAddress == "" {
		return nil, fmt.Errorf("wallet address cannot be empty")
	}

	// Get wallet
	wallet, found := k.GetWallet(ctx, walletAddress)
	if !found {
		return nil, fmt.Errorf("wallet not found: %s", walletAddress)
	}

	// Create recovery info
	now := time.Now().UTC()
	recoveryInfo := &types.RecoveryInfo{
		Did:             wallet.Creator,
		RecoveryToken:   wallet.Creator, // TODO: implement proper token generation
		RecoveryAddress: wallet.Creator,
		Status:          types.RecoveryStatus_RECOVERY_STATUS_PENDING,
		CreatedAt:       now,
	}

	// Store recovery info
	k.StoreRecoveryInfo(ctx, recoveryInfo)

	return recoveryInfo, nil
}

// VerifyRecoverySession verifies a recovery session
func (k Keeper) VerifyRecoverySession(ctx sdk.Context, walletAddress string, token string) error {
	if walletAddress == "" {
		return fmt.Errorf("wallet address cannot be empty")
	}
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Get recovery info
	info, err := k.GetRecoveryInfo(ctx, walletAddress)
	if err != nil {
		return fmt.Errorf("failed to get recovery info: %v", err)
	}

	// Verify token
	if info.RecoveryToken != token {
		return fmt.Errorf("invalid recovery token")
	}

	// Check status
	if info.Status != types.RecoveryStatus_RECOVERY_STATUS_PENDING {
		return fmt.Errorf("recovery not in pending state")
	}

	return nil
}

// CompleteRecoverySession completes a recovery session
func (k Keeper) CompleteRecoverySession(ctx sdk.Context, walletAddress string, token string) error {
	if walletAddress == "" {
		return fmt.Errorf("wallet address cannot be empty")
	}
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Get recovery info
	info, err := k.GetRecoveryInfo(ctx, walletAddress)
	if err != nil {
		return fmt.Errorf("failed to get recovery info: %v", err)
	}

	// Verify token
	if info.RecoveryToken != token {
		return fmt.Errorf("invalid recovery token")
	}

	// Check status
	if info.Status != types.RecoveryStatus_RECOVERY_STATUS_PENDING {
		return fmt.Errorf("recovery not in pending state")
	}

	// Update status
	info.Status = types.RecoveryStatus_RECOVERY_STATUS_COMPLETED
	k.StoreRecoveryInfo(ctx, info)

	return nil
}

// GetRecoveryInfo gets recovery info for a wallet
func (k Keeper) GetRecoveryInfo(ctx sdk.Context, walletAddress string) (*types.RecoveryInfo, error) {
	if walletAddress == "" {
		return nil, fmt.Errorf("wallet address cannot be empty")
	}

	store := ctx.KVStore(k.storeKey)
	key := k.GetRecoveryKey(walletAddress)
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("recovery info not found for wallet address: %s", walletAddress)
	}

	var info types.RecoveryInfo
	if err := k.cdc.Unmarshal(bz, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recovery info: %v", err)
	}

	return &info, nil
}

// StoreRecoveryInfo stores recovery info for a wallet
func (k Keeper) StoreRecoveryInfo(ctx sdk.Context, info *types.RecoveryInfo) {
	store := ctx.KVStore(k.storeKey)
	key := k.GetRecoveryKey(info.Did)
	bz := k.cdc.MustMarshal(info)
	store.Set(key, bz)
}

// GetRecoveryKey returns the key for storing recovery info
func (k Keeper) GetRecoveryKey(did string) []byte {
	return []byte(fmt.Sprintf("recovery/%s", did))
}

// DeleteRecoveryInfo deletes recovery info for a wallet
func (k Keeper) DeleteRecoveryInfo(ctx sdk.Context, did string) error {
	store := ctx.KVStore(k.storeKey)
	key := k.GetRecoveryKey(did)
	store.Delete(key)
	return nil
}

// RecoverWallet recovers a wallet by its address
func (k Keeper) RecoverWallet(ctx sdk.Context, msg *types.MsgRecoverWallet) error {
	// Get wallet
	wallet, found := k.GetWallet(ctx, msg.WalletAddress)
	if !found {
		return fmt.Errorf("wallet not found: %s", msg.WalletAddress)
	}

	// Validate recovery session
	if err := k.VerifyRecoverySession(ctx, wallet.Creator, msg.RecoveryProof); err != nil {
		return fmt.Errorf("failed to verify recovery session: %v", err)
	}

	// Delete the old wallet
	if err := k.DeleteWallet(ctx, msg.WalletAddress); err != nil {
		return fmt.Errorf("failed to delete old wallet: %v", err)
	}

	// Create new wallet with updated public key
	now := ctx.BlockTime()
	newWallet := &types.Wallet{
		Id:            msg.WalletAddress,
		Creator:       msg.Creator,
		PublicKey:     msg.NewPubKey,
		WalletAddress: msg.WalletAddress,
		ChainId:      wallet.ChainId,
		Status:       types.WalletStatus_WALLET_STATUS_ACTIVE,
		KeyVersion:   wallet.KeyVersion + 1,
		CreatedAt:    &now,
		UpdatedAt:    &now,
		LastUsed:     &now,
		UsageCount:   0,
	}

	// Save the new wallet
	if err := k.SaveWallet(ctx, newWallet); err != nil {
		return fmt.Errorf("failed to save new wallet: %v", err)
	}

	// Delete recovery info
	if err := k.DeleteRecoveryInfo(ctx, wallet.Creator); err != nil {
		return fmt.Errorf("failed to delete recovery info: %v", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWalletRecovered,
			sdk.NewAttribute(types.AttributeKeyWalletAddress, msg.WalletAddress),
			sdk.NewAttribute(types.AttributeKeyNewOwner, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyTimestamp, ctx.BlockTime().String()),
		),
	)

	return nil
}
