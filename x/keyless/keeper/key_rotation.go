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

// GetKeyRotation gets a key rotation record
func (k Keeper) GetKeyRotation(ctx sdk.Context, walletAddress string, version uint64) (*types.KeyRotation, error) {
	key := types.GetKeyRotationKey(walletAddress, version)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("key rotation not found")
	}

	var rotation types.KeyRotation
	if err := k.cdc.Unmarshal(bz, &rotation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal key rotation: %v", err)
	}
	return &rotation, nil
}

// SaveKeyRotation saves a key rotation record
func (k Keeper) SaveKeyRotation(ctx sdk.Context, rotation *types.KeyRotation) error {
	if rotation == nil {
		return fmt.Errorf("key rotation cannot be nil")
	}

	key := types.GetKeyRotationKey(rotation.WalletAddress, rotation.Version)
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(rotation)
	if err != nil {
		return fmt.Errorf("failed to marshal key rotation: %v", err)
	}
	store.Set(key, bz)
	return nil
}

// InitiateKeyRotation starts the key rotation process for a wallet
func (k Keeper) InitiateKeyRotation(ctx sdk.Context, msg *types.MsgInitiateKeyRotation) error {
	// Get wallet
	wallet, found := k.GetWallet(ctx, msg.WalletAddress)
	if !found {
		return fmt.Errorf("wallet not found: %s", msg.WalletAddress)
	}

	// Check if there's already an ongoing key rotation
	status, err := k.GetKeyRotationStatus(ctx, msg.WalletAddress)
	if err != nil {
		return err
	}
	if status.Status == types.KeyRotationStatus_KEY_ROTATION_STATUS_IN_PROGRESS {
		return fmt.Errorf("key rotation already in progress for wallet: %s", msg.WalletAddress)
	}

	// Create new key rotation record
	nextVersion := uint64(wallet.KeyVersion) + 1
	rotation := &types.KeyRotation{
		WalletAddress: msg.WalletAddress,
		OldPubKey:     wallet.PublicKey,
		NewPubKey:     msg.NewPubKey,
		Version:       nextVersion,
		Status:        types.KeyRotationStatus_KEY_ROTATION_STATUS_IN_PROGRESS,
		Error:         "",
		Signature:     msg.Signature,
	}

	// Store key rotation record
	if err := k.SaveKeyRotation(ctx, rotation); err != nil {
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
	wallet, found := k.GetWallet(ctx, msg.WalletAddress)
	if !found {
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
	k.SetWallet(ctx, wallet)

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
	if err := k.SaveKeyRotation(ctx, rotation); err != nil {
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

// GetAllKeyRotations gets all key rotations for a wallet
func (k Keeper) GetAllKeyRotations(ctx sdk.Context, walletAddress string) ([]*types.KeyRotation, error) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte(types.KeyRotationKey))

	var rotations []*types.KeyRotation
	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var rotation types.KeyRotation
		err := k.cdc.Unmarshal(iterator.Value(), &rotation)
		if err != nil {
			return nil, err
		}
		if rotation.WalletAddress == walletAddress {
			rotations = append(rotations, &rotation)
		}
	}

	return rotations, nil
}

// GetKeyRotationsForWallet returns all key rotations for a wallet
func (k Keeper) GetKeyRotationsForWallet(ctx sdk.Context, walletAddress string) []*types.KeyRotation {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte(types.KeyRotationKey))

	var rotations []*types.KeyRotation
	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var rotation types.KeyRotation
		err := k.cdc.Unmarshal(iterator.Value(), &rotation)
		if err != nil {
			continue
		}
		if rotation.WalletAddress == walletAddress {
			rotations = append(rotations, &rotation)
		}
	}

	return rotations
}
