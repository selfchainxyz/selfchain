package keygen

import (
	"encoding/base64"
	"testing"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"selfchain/x/keyless/types"
)

func TestNewEncryptionManager(t *testing.T) {
	em := NewEncryptionManager()
	require.NotNil(t, em)
	require.NotNil(t, em.keyStore)
}

func TestEncryptionManager_EncryptDecryptShare(t *testing.T) {
	em := NewEncryptionManager()
	require.NotNil(t, em)

	// Create a mock LocalPartySaveData
	mockData := &keygen.LocalPartySaveData{
		LocalPreParams: keygen.LocalPreParams{},
	}

	// Test encryption
	encryptedShare, err := em.EncryptShare(mockData)
	require.NoError(t, err)
	require.NotNil(t, encryptedShare)
	assert.NotEmpty(t, encryptedShare.KeyId)
	assert.NotEmpty(t, encryptedShare.EncryptedData)

	// Verify key is stored
	em.mu.RLock()
	_, ok := em.keyStore[encryptedShare.KeyId]
	em.mu.RUnlock()
	assert.True(t, ok)

	// Test decryption
	decryptedData, err := em.DecryptShare(encryptedShare)
	require.NoError(t, err)
	require.NotNil(t, decryptedData)
	assert.Equal(t, mockData.LocalPreParams, decryptedData.LocalPreParams)
}

func TestEncryptionManager_EncryptShare_NilData(t *testing.T) {
	em := NewEncryptionManager()
	require.NotNil(t, em)

	// Test with nil data
	encryptedShare, err := em.EncryptShare(nil)
	assert.Error(t, err)
	assert.Nil(t, encryptedShare)
	assert.Contains(t, err.Error(), "nil data")
}

func TestEncryptionManager_DecryptShare_InvalidShare(t *testing.T) {
	em := NewEncryptionManager()
	require.NotNil(t, em)

	// Test with nil share
	decryptedData, err := em.DecryptShare(nil)
	assert.Error(t, err)
	assert.Nil(t, decryptedData)
	assert.Contains(t, err.Error(), "nil share")

	// Test with invalid key ID
	invalidShare := &types.EncryptedShare{
		KeyId:         "nonexistent_key",
		EncryptedData: base64.StdEncoding.EncodeToString([]byte("invalid_data")),
	}
	decryptedData, err = em.DecryptShare(invalidShare)
	assert.Error(t, err)
	assert.Nil(t, decryptedData)
	assert.Contains(t, err.Error(), "key not found")

	// Test with invalid encrypted data
	validShare, err := em.EncryptShare(&keygen.LocalPartySaveData{})
	require.NoError(t, err)
	validShare.EncryptedData = base64.StdEncoding.EncodeToString([]byte("invalid_data"))
	decryptedData, err = em.DecryptShare(validShare)
	assert.Error(t, err)
	assert.Nil(t, decryptedData)
	assert.Contains(t, err.Error(), "failed to decrypt")
}

func TestEncryptionManager_GetKeyForShare(t *testing.T) {
	em := NewEncryptionManager()
	require.NotNil(t, em)

	// Test with nonexistent key
	key, err := em.getKeyForShare("nonexistent_key")
	assert.Error(t, err)
	assert.Nil(t, key)
	assert.Contains(t, err.Error(), "key not found")

	// Test with valid key
	mockData := &keygen.LocalPartySaveData{}
	encryptedShare, err := em.EncryptShare(mockData)
	require.NoError(t, err)

	key, err = em.getKeyForShare(encryptedShare.KeyId)
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.Len(t, key, 32)
}
