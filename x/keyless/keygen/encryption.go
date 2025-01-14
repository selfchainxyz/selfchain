package keygen

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"

	"selfchain/x/keyless/types"
	"selfchain/x/keyless/crypto"
)

// EncryptionManager handles key encryption
type EncryptionManager struct {
	masterKey []byte
}

// NewEncryptionManager creates a new encryption manager
func NewEncryptionManager() *EncryptionManager {
	// In production, this should be loaded from a secure source or HSM
	return &EncryptionManager{
		masterKey: make([]byte, 32), // Placeholder for demo
	}
}

// EncryptShare encrypts a party's save data
func (em *EncryptionManager) EncryptShare(data *keygen.LocalPartySaveData) (*types.EncryptedShare, error) {
	// 1. Generate unique encryption key
	key, err := crypto.NewEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}

	// 2. Marshal the save data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal share data: %w", err)
	}

	// 3. Encrypt the JSON data
	encryptedData, err := crypto.Encrypt(key, jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt share data: %w", err)
	}

	// 4. Create encrypted share with metadata
	return &types.EncryptedShare{
		EncryptedData: encryptedData,
		KeyId:        uuid.New().String(),
		Version:      1,
		CreatedAt:    time.Now(),
	}, nil
}

// DecryptShare decrypts an encrypted share
func (em *EncryptionManager) DecryptShare(share *types.EncryptedShare) (*keygen.LocalPartySaveData, error) {
	// 1. Get encryption key for the share
	key, err := em.getKeyForShare(share.KeyId)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	// 2. Decrypt the data
	decryptedData, err := crypto.Decrypt(key, share.EncryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt share data: %w", err)
	}

	// 3. Unmarshal the JSON data
	var saveData keygen.LocalPartySaveData
	if err := json.Unmarshal(decryptedData, &saveData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal share data: %w", err)
	}

	return &saveData, nil
}

// getKeyForShare retrieves the encryption key for a share
func (em *EncryptionManager) getKeyForShare(keyID string) (crypto.EncryptionKey, error) {
	// In production, this should retrieve the key from a secure key management service
	// For now, return a placeholder key
	return make([]byte, 32), nil
}
