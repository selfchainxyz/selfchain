package keygen

import (
	"context"
	"fmt"
	"time"

	"selfchain/x/keyless/types"
	"selfchain/x/keyless/storage"
)

// KeyRotationManager handles periodic key rotation
type KeyRotationManager struct {
	keyGen  *KeyGenManager
	storage *storage.Storage
}

// NewKeyRotationManager creates a new key rotation manager
func NewKeyRotationManager(keyGen *KeyGenManager, storage *storage.Storage) *KeyRotationManager {
	return &KeyRotationManager{
		keyGen:  keyGen,
		storage: storage,
	}
}

// RotateKeyShares performs secure key rotation
func (kr *KeyRotationManager) RotateKeyShares(ctx context.Context, walletAddress string) error {
	// 1. Load existing shares
	share1, err := kr.storage.GetPartyShare(ctx, fmt.Sprintf("%s_share_1", walletAddress))
	if err != nil {
		return fmt.Errorf("failed to get share 1: %w", err)
	}

	share2, err := kr.storage.GetPartyShare(ctx, fmt.Sprintf("%s_share_2", walletAddress))
	if err != nil {
		return fmt.Errorf("failed to get share 2: %w", err)
	}

	// 2. Generate new key shares
	req := &types.KeyGenRequest{
		WalletAddress:  walletAddress,
		ChainId:       "temp", // This should be retrieved from wallet metadata
		SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_STANDARD,
	}

	resp, err := kr.keyGen.GenerateKeyShares(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to generate new shares: %w", err)
	}

	// 3. Validate new shares
	if err := kr.validateNewShares(ctx, resp); err != nil {
		return fmt.Errorf("new shares validation failed: %w", err)
	}

	// 4. Backup old shares
	if err := kr.backupOldShares(ctx, walletAddress, share1, share2); err != nil {
		return fmt.Errorf("failed to backup old shares: %w", err)
	}

	// 5. Update metadata
	resp.Metadata.LastRotated = time.Now()

	return nil
}

func (kr *KeyRotationManager) validateNewShares(ctx context.Context, resp *types.KeyGenResponse) error {
	// Implement share validation logic
	// This should verify that the new shares can successfully sign transactions
	return nil
}

func (kr *KeyRotationManager) backupOldShares(ctx context.Context, walletAddress string, share1, share2 *types.EncryptedShare) error {
	// Implement backup logic
	// This should store old shares in a secure backup location
	return nil
}

// ShouldRotate checks if key rotation is needed
func (kr *KeyRotationManager) ShouldRotate(metadata *types.KeyMetadata) bool {
	if metadata == nil {
		return false
	}

	// Check time since last rotation
	if time.Since(metadata.LastRotated) > 30*24*time.Hour { // 30 days
		return true
	}

	// Check usage count
	if metadata.UsageCount > 1000 { // Example threshold
		return true
	}

	return false
}
