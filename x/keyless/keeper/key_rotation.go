package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

// GetKeyRotationStatus retrieves the key rotation status for a given wallet
func (k Keeper) GetKeyRotationStatus(ctx sdk.Context, walletAddress string) (*types.KeyRotationStatusInfo, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.KeyRotationStatusPrefix))
	key := types.KeyRotationStatusKey(walletAddress)
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("key rotation status not found for wallet %s", walletAddress)
	}

	var status types.KeyRotationStatusInfo
	k.cdc.MustUnmarshal(bz, &status)
	return &status, nil
}

// SetKeyRotationStatus stores the key rotation status for a wallet
func (k Keeper) SetKeyRotationStatus(ctx sdk.Context, status *types.KeyRotationStatusInfo) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.KeyRotationStatusPrefix))
	key := types.KeyRotationStatusKey(status.WalletAddress)
	bz := k.cdc.MustMarshal(status)
	store.Set(key, bz)
	return nil
}

// DeleteKeyRotationStatus deletes the key rotation status for a wallet
func (k Keeper) DeleteKeyRotationStatus(ctx sdk.Context, walletAddress string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.KeyRotationStatusPrefix))
	key := types.KeyRotationStatusKey(walletAddress)
	store.Delete(key)
}

// GetKeyRotation retrieves a key rotation record by wallet ID and version
func (k Keeper) GetKeyRotation(ctx sdk.Context, walletId string, version uint64) (*types.KeyRotation, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.KeyRotationKey))
	key := types.GetKeyRotationKey(walletId, version)
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("key rotation not found for wallet %s version %d", walletId, version)
	}

	var rotation types.KeyRotation
	k.cdc.MustUnmarshal(bz, &rotation)
	return &rotation, nil
}

// SetKeyRotation stores a key rotation record
func (k Keeper) SetKeyRotation(ctx sdk.Context, rotation *types.KeyRotation) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.KeyRotationKey))
	key := types.GetKeyRotationKey(rotation.WalletId, rotation.Version)
	bz := k.cdc.MustMarshal(rotation)
	store.Set(key, bz)
	return nil
}

// InitiateKeyRotation starts the key rotation process for a wallet
func (k Keeper) InitiateKeyRotation(ctx sdk.Context, msg *types.MsgInitiateKeyRotation) error {
	// Get wallet
	wallet, err := k.GetWallet(ctx, msg.WalletAddress)
	if err != nil {
		return fmt.Errorf("wallet not found: %s", msg.WalletAddress)
	}

	// Check if there's already an ongoing key rotation
	status, err := k.GetKeyRotationStatus(ctx, msg.WalletAddress)
	if err == nil && status.Status == types.KeyRotationStatus_KEY_ROTATION_STATUS_IN_PROGRESS {
		return fmt.Errorf("key rotation already in progress for wallet: %s", msg.WalletAddress)
	}

	// Create new key rotation record
	nextVersion := uint64(wallet.KeyVersion) + 1
	rotation := &types.KeyRotation{
		WalletId:  msg.WalletAddress,
		OldPubKey: wallet.PublicKey,
		NewPubKey: msg.NewPubKey,
		Version:   nextVersion,
		Status:    types.KeyRotationStatus_KEY_ROTATION_STATUS_IN_PROGRESS,
		Error:     "",
		Signature: msg.Signature,
	}

	// Store key rotation record
	if err := k.SetKeyRotation(ctx, rotation); err != nil {
		return err
	}

	// Update key rotation status
	status = &types.KeyRotationStatusInfo{
		WalletAddress: msg.WalletAddress,
		NewPublicKey:  msg.NewPubKey,
		Status:        types.KeyRotationStatus_KEY_ROTATION_STATUS_IN_PROGRESS,
		Version:       nextVersion,
	}
	if err := k.SetKeyRotationStatus(ctx, status); err != nil {
		return err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeKeyRotationInitiated,
			sdk.NewAttribute(types.AttributeKeyWalletAddress, msg.WalletAddress),
			sdk.NewAttribute(types.AttributeKeyVersion, fmt.Sprintf("%d", rotation.Version)),
			sdk.NewAttribute(types.AttributeKeyNewPubKey, msg.NewPubKey),
		),
	)

	return nil
}

// CompleteKeyRotation completes the key rotation process
func (k Keeper) CompleteKeyRotation(ctx sdk.Context, msg *types.MsgCompleteKeyRotation) error {
	// Get wallet
	wallet, err := k.GetWallet(ctx, msg.WalletAddress)
	if err != nil {
		return fmt.Errorf("wallet not found: %s", msg.WalletAddress)
	}

	// Get key rotation record
	nextVersion := uint64(wallet.KeyVersion) + 1
	rotation, err := k.GetKeyRotation(ctx, msg.WalletAddress, nextVersion)
	if err != nil {
		return fmt.Errorf("key rotation not found for wallet: %s", msg.WalletAddress)
	}

	// Verify signature
	if msg.Signature != rotation.Signature {
		return fmt.Errorf("invalid signature")
	}

	// Update wallet with new public key
	wallet.PublicKey = rotation.NewPubKey
	wallet.KeyVersion = uint32(rotation.Version)
	if err := k.SaveWallet(ctx, wallet); err != nil {
		return err
	}

	// Update key rotation status
	status := &types.KeyRotationStatusInfo{
		WalletAddress: msg.WalletAddress,
		NewPublicKey:  rotation.NewPubKey,
		Status:        types.KeyRotationStatus_KEY_ROTATION_STATUS_COMPLETED,
		Version:       rotation.Version,
	}
	if err := k.SetKeyRotationStatus(ctx, status); err != nil {
		return err
	}

	// Update key rotation record
	rotation.Status = types.KeyRotationStatus_KEY_ROTATION_STATUS_COMPLETED
	if err := k.SetKeyRotation(ctx, rotation); err != nil {
		return err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeKeyRotationCompleted,
			sdk.NewAttribute(types.AttributeKeyWalletAddress, msg.WalletAddress),
			sdk.NewAttribute(types.AttributeKeyVersion, fmt.Sprintf("%d", rotation.Version)),
			sdk.NewAttribute(types.AttributeKeyNewPubKey, rotation.NewPubKey),
		),
	)

	return nil
}

// GetKeyRotationsForWallet returns all key rotations for a wallet
func (k Keeper) GetKeyRotationsForWallet(ctx sdk.Context, walletId string) []*types.KeyRotation {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.KeyRotationKey))
	iterator := store.Iterator([]byte(walletId), nil)
	defer iterator.Close()

	var rotations []*types.KeyRotation
	for ; iterator.Valid(); iterator.Next() {
		var rotation types.KeyRotation
		k.cdc.MustUnmarshal(iterator.Value(), &rotation)
		if rotation.WalletId == walletId {
			rotations = append(rotations, &rotation)
		}
	}

	return rotations
}
