package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// EncryptionKey represents a key used for encryption/decryption
type EncryptionKey []byte

// NewEncryptionKey generates a random 32-byte key for AES-256
func NewEncryptionKey() (EncryptionKey, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}
	return key, nil
}

// Encrypt encrypts data using AES-256-GCM
func Encrypt(key EncryptionKey, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Never use more than 2^32 random nonces with a given key
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Encrypt and append nonce
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	encrypted := append(nonce, ciphertext...)

	// Encode as base64 for storage
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt decrypts data using AES-256-GCM
func Decrypt(key EncryptionKey, encryptedStr string) ([]byte, error) {
	// Decode from base64
	encrypted, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	if len(encrypted) < 12 {
		return nil, fmt.Errorf("invalid ciphertext length")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := encrypted[:12]
	ciphertext := encrypted[12:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}
