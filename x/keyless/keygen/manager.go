package keygen

import (
	"context"
	"fmt"
	"time"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"

	"selfchain/x/keyless/types"
	"selfchain/x/keyless/storage"
	"selfchain/x/keyless/tss"
)

// KeyGenManager handles the key generation lifecycle
type KeyGenManager struct {
	keeper       storage.Storage
	encryptionMgr *EncryptionManager
}

// NewKeyGenManager creates a new key generation manager
func NewKeyGenManager(keeper storage.Storage) *KeyGenManager {
	return &KeyGenManager{
		keeper:       keeper,
		encryptionMgr: NewEncryptionManager(),
	}
}

// GenerateKeyShares creates a new key pair with enhanced security
func (km *KeyGenManager) GenerateKeyShares(ctx context.Context, req *types.KeyGenRequest) (*types.KeyGenResponse, error) {
	// 1. Validate request
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Check if shares already exist
	key1 := fmt.Sprintf("%s_share_1", req.WalletAddress)
	key2 := fmt.Sprintf("%s_share_2", req.WalletAddress)
	
	share1, err := km.keeper.GetPartyShare(ctx, key1)
	if err == nil && share1 != nil {
		return nil, fmt.Errorf("key shares already exist for wallet %s", req.WalletAddress)
	}
	
	share2, err := km.keeper.GetPartyShare(ctx, key2)
	if err == nil && share2 != nil {
		return nil, fmt.Errorf("key shares already exist for wallet %s", req.WalletAddress)
	}

	// 2. Generate pre-parameters with timeout
	preParams, err := km.generateSecurePreParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("pre-params generation failed: %w", err)
	}

	// 3. Generate key shares using TSS
	keygenResult, err := tss.GenerateKey(ctx, preParams, req.ChainId)
	if err != nil {
		return nil, fmt.Errorf("key generation failed: %w", err)
	}

	// Add metadata
	now := time.Now()
	metadata := &types.KeyMetadata{
		CreatedAt:     now,
		LastRotated:   now,
		LastUsed:      now,
		UsageCount:    0,
		BackupStatus:  types.BackupStatus_BACKUP_STATUS_NONE,
		SecurityLevel: req.SecurityLevel,
	}

	// 4. Encrypt key shares
	encryptedShares, err := km.encryptKeyShares(ctx, keygenResult)
	if err != nil {
		return nil, fmt.Errorf("share encryption failed: %w", err)
	}

	// 5. Store encrypted shares
	if err := km.storeKeyShares(ctx, req.WalletAddress, encryptedShares); err != nil {
		return nil, fmt.Errorf("share storage failed: %w", err)
	}

	// 6. Store party data
	if err := km.keeper.SavePartyData(ctx, req.WalletAddress, keygenResult.Party1Data, keygenResult.Party2Data); err != nil {
		// Cleanup stored shares on error
		km.keeper.DeletePartyShare(ctx, key1)
		km.keeper.DeletePartyShare(ctx, key2)
		return nil, fmt.Errorf("party data storage failed: %w", err)
	}

	return &types.KeyGenResponse{
		WalletAddress: req.WalletAddress,
		PublicKey:    keygenResult.PublicKeyBytes,
		Metadata:     metadata,
	}, nil
}

func (km *KeyGenManager) generateSecurePreParams(ctx context.Context) (*keygen.LocalPreParams, error) {
	// Add timeout for pre-params generation
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	preParams, err := keygen.GeneratePreParams(1 * time.Minute)
	if err != nil {
		return nil, err
	}

	return preParams, nil
}

func (km *KeyGenManager) encryptKeyShares(ctx context.Context, result *tss.KeygenResult) ([]*types.EncryptedShare, error) {
	shares := make([]*types.EncryptedShare, 2)

	// Encrypt Party1 Data
	share1, err := km.encryptionMgr.EncryptShare(result.Party1Data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt party 1 share: %w", err)
	}
	shares[0] = share1

	// Encrypt Party2 Data
	share2, err := km.encryptionMgr.EncryptShare(result.Party2Data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt party 2 share: %w", err)
	}
	shares[1] = share2

	return shares, nil
}

func (km *KeyGenManager) storeKeyShares(ctx context.Context, walletAddress string, shares []*types.EncryptedShare) error {
	for i, share := range shares {
		key := fmt.Sprintf("%s_share_%d", walletAddress, i+1)
		if err := km.keeper.SavePartyShare(ctx, key, share); err != nil {
			// Cleanup any previously stored shares on error
			for j := 0; j < i; j++ {
				km.keeper.DeletePartyShare(ctx, fmt.Sprintf("%s_share_%d", walletAddress, j+1))
			}
			return fmt.Errorf("failed to store share %d: %w", i+1, err)
		}
	}
	return nil
}
