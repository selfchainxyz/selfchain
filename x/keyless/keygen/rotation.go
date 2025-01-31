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
	storage storage.Storage
}

// NewKeyRotationManager creates a new key rotation manager
func NewKeyRotationManager(keyGen *KeyGenManager, storage storage.Storage) *KeyRotationManager {
	return &KeyRotationManager{
		keyGen:  keyGen,
		storage: storage,
	}
}

// RotateKeyShares performs secure key rotation
func (kr *KeyRotationManager) RotateKeyShares(ctx context.Context, walletAddress string) error {
	// 1. Get wallet metadata
	metadata, err := kr.storage.GetWalletMetadata(ctx, walletAddress)
	if err != nil {
		return fmt.Errorf("failed to get wallet metadata: %w", err)
	}

	// 2. Delete old shares
	key1 := fmt.Sprintf("%s_share_1", walletAddress)
	key2 := fmt.Sprintf("%s_share_2", walletAddress)
	kr.storage.DeletePartyShare(ctx, key1)
	kr.storage.DeletePartyShare(ctx, key2)

	// 3. Generate new key shares
	req := &types.KeyGenRequest{
		WalletAddress:  walletAddress,
		ChainId:       metadata.ChainId,
		SecurityLevel: metadata.SecurityLevel,
	}

	resp, err := kr.keyGen.GenerateKeyShares(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to generate new shares: %w", err)
	}

	// 4. Validate new shares
	if err := kr.validateNewShares(ctx, resp); err != nil {
		return fmt.Errorf("new shares validation failed: %w", err)
	}

	return nil
}

func (kr *KeyRotationManager) validateNewShares(ctx context.Context, resp *types.KeyGenResponse) error {
	// For now, just verify that we have valid response data
	if resp == nil {
		return fmt.Errorf("nil response")
	}
	if resp.PublicKey == nil {
		return fmt.Errorf("missing public key")
	}
	if resp.Metadata == nil {
		return fmt.Errorf("missing metadata")
	}
	return nil
}

// ShouldRotate checks if key rotation is needed
func (kr *KeyRotationManager) ShouldRotate(metadata *types.KeyMetadata) bool {
	if metadata == nil {
		return false
	}

	var lastRotation time.Time
	if metadata.LastRotated.IsZero() {
		lastRotation = metadata.CreatedAt
	} else {
		lastRotation = metadata.LastRotated
	}

	timeSinceRotation := time.Since(lastRotation)

	// Check if key age exceeds rotation threshold
	switch metadata.SecurityLevel {
	case types.SecurityLevel_SECURITY_LEVEL_HIGH:
		return timeSinceRotation >= 30*24*time.Hour // 30 days
	case types.SecurityLevel_SECURITY_LEVEL_STANDARD:
		return timeSinceRotation >= 90*24*time.Hour // 90 days
	case types.SecurityLevel_SECURITY_LEVEL_ENTERPRISE:
		return timeSinceRotation >= 180*24*time.Hour // 180 days
	default:
		return false
	}
}
