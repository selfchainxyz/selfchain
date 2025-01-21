package keeper

import (
	"fmt"
	"time"

	"selfchain/x/keyless/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	keyRotationPrefix = "key_rotation"
)

// GetKeyRotationStore returns the store for key rotation data
func (k Keeper) GetKeyRotationStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, []byte(keyRotationPrefix))
}

// InitKeyRotation initializes a key rotation operation for a wallet
func (k Keeper) InitKeyRotation(ctx sdk.Context, walletId string, newPubKey string, signature string) (*types.KeyRotation, error) {
	// Verify wallet exists
	wallet, found := k.GetWalletFromStore(ctx, walletId)
	if !found {
		return nil, fmt.Errorf("wallet not found: %s", walletId)
	}

	// Get current version
	version := uint64(wallet.KeyVersion) + 1

	// Create new key rotation
	now := time.Now()
	rotation := &types.KeyRotation{
		WalletId:  walletId,
		Version:   version,
		NewPubKey: newPubKey,
		Status:    "pending",
		CreatedAt: &now,
		Signature: signature,
	}

	// Store key rotation
	k.setKeyRotation(ctx, rotation)

	return rotation, nil
}

// CompleteKeyRotation completes a pending key rotation operation
func (k Keeper) CompleteKeyRotation(ctx sdk.Context, walletId string, version uint64) (*types.KeyRotation, error) {
	// Get key rotation
	rotation, found := k.GetKeyRotation(ctx, walletId, version)
	if !found {
		return nil, fmt.Errorf("key rotation not found: wallet=%s, version=%d", walletId, version)
	}

	// Verify status is pending
	if rotation.Status != "pending" {
		return nil, fmt.Errorf("key rotation is not pending: %s", rotation.Status)
	}

	// Update wallet with new public key
	wallet, found := k.GetWalletFromStore(ctx, walletId)
	if !found {
		return nil, fmt.Errorf("wallet not found: %s", walletId)
	}

	// Update wallet's public key and version
	wallet.PublicKey = rotation.NewPubKey
	wallet.KeyVersion = uint32(rotation.Version)
	k.SetWallet(ctx, &wallet)

	// Update rotation status
	now := time.Now()
	rotation.Status = "completed"
	rotation.CompletedAt = &now
	k.setKeyRotation(ctx, rotation)

	return rotation, nil
}

// CancelKeyRotation cancels a pending key rotation operation
func (k Keeper) CancelKeyRotation(ctx sdk.Context, walletId string, version uint64) (*types.KeyRotation, error) {
	// Get key rotation
	rotation, found := k.GetKeyRotation(ctx, walletId, version)
	if !found {
		return nil, fmt.Errorf("key rotation not found: wallet=%s, version=%d", walletId, version)
	}

	// Verify status is pending
	if rotation.Status != "pending" {
		return nil, fmt.Errorf("key rotation is not pending: %s", rotation.Status)
	}

	// Update rotation status
	now := time.Now()
	rotation.Status = "cancelled"
	rotation.CompletedAt = &now
	k.setKeyRotation(ctx, rotation)

	return rotation, nil
}

// GetKeyRotation gets a key rotation by wallet ID and version
func (k Keeper) GetKeyRotation(ctx sdk.Context, walletId string, version uint64) (*types.KeyRotation, bool) {
	store := k.GetKeyRotationStore(ctx)
	key := []byte(fmt.Sprintf("%s_%d", walletId, version))

	bz := store.Get(key)
	if bz == nil {
		return nil, false
	}

	var rotation types.KeyRotation
	k.cdc.MustUnmarshal(bz, &rotation)
	return &rotation, true
}

// GetKeyRotations gets all key rotations for a wallet
func (k Keeper) GetKeyRotations(ctx sdk.Context, walletId string) []types.KeyRotation {
	store := k.GetKeyRotationStore(ctx)
	prefix := []byte(walletId + "_")

	iterator := store.Iterator(prefix, nil)
	defer iterator.Close()

	var rotations []types.KeyRotation
	for ; iterator.Valid(); iterator.Next() {
		var rotation types.KeyRotation
		k.cdc.MustUnmarshal(iterator.Value(), &rotation)
		rotations = append(rotations, rotation)
	}

	return rotations
}

// setKeyRotation stores a key rotation in the KVStore
func (k Keeper) setKeyRotation(ctx sdk.Context, rotation *types.KeyRotation) {
	store := k.GetKeyRotationStore(ctx)
	key := []byte(fmt.Sprintf("%s_%d", rotation.WalletId, rotation.Version))
	bz := k.cdc.MustMarshal(rotation)
	store.Set(key, bz)
}
