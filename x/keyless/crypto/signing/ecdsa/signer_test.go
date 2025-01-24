package ecdsa

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"selfchain/x/keyless/crypto/signing/types"
)

func setupTestKeys(t *testing.T) (*btcec.PrivateKey, *btcec.PublicKey) {
	privKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	return privKey, privKey.PubKey()
}

func TestECDSASigner_Sign(t *testing.T) {
	privKey, pubKey := setupTestKeys(t)
	signer := NewECDSASigner(privKey, pubKey)
	require.NotNil(t, signer)

	message := []byte("test message")
	result, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.R)
	require.NotNil(t, result.S)
	require.NotEmpty(t, result.Bytes)
}

func TestECDSASigner_Verify(t *testing.T) {
	privKey, pubKey := setupTestKeys(t)
	signer := NewECDSASigner(privKey, pubKey)
	require.NotNil(t, signer)

	message := []byte("test message")
	signature, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	pubKeyBytes := pubKey.SerializeCompressed()
	valid, err := signer.Verify(context.Background(), message, signature, pubKeyBytes)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestECDSASigner_VerifyWithDifferentFormats(t *testing.T) {
	privKey, pubKey := setupTestKeys(t)
	signer := NewECDSASigner(privKey, pubKey)
	require.NotNil(t, signer)

	message := []byte("test message")
	signature, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	// Test with compressed public key
	compressedPubKey := pubKey.SerializeCompressed()
	valid, err := signer.Verify(context.Background(), message, signature, compressedPubKey)
	require.NoError(t, err)
	assert.True(t, valid)

	// Test with uncompressed public key
	uncompressedPubKey := pubKey.SerializeUncompressed()
	valid, err = signer.Verify(context.Background(), message, signature, uncompressedPubKey)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestECDSASigner_InvalidSignature(t *testing.T) {
	privKey, pubKey := setupTestKeys(t)
	signer := NewECDSASigner(privKey, pubKey)
	require.NotNil(t, signer)

	message := []byte("test message")
	signature, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	// Modify the signature to make it invalid
	signature.R.Add(signature.R, btcec.S256().P)

	valid, err := signer.Verify(context.Background(), message, signature, pubKey.SerializeCompressed())
	require.NoError(t, err)
	assert.False(t, valid)
}

func TestECDSASigner_InvalidPublicKey(t *testing.T) {
	privKey, pubKey := setupTestKeys(t)
	signer := NewECDSASigner(privKey, pubKey)
	require.NotNil(t, signer)

	message := []byte("test message")
	signature, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	// Test with invalid public key
	invalidPubKey := []byte("invalid public key")
	_, err = signer.Verify(context.Background(), message, signature, invalidPubKey)
	require.Error(t, err)
}

func TestECDSASigner_KnownTestVector(t *testing.T) {
	// Known test vector from Bitcoin's test cases
	privKeyHex := "0000000000000000000000000000000000000000000000000000000000000001"
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	require.NoError(t, err)

	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)
	require.NotNil(t, privKey)

	signer := NewECDSASigner(privKey, privKey.PubKey())
	require.NotNil(t, signer)

	message := []byte("test message")
	signature, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)
	require.NotNil(t, signature)
}

func TestECDSASigner_SignAndVerifyLargeMessage(t *testing.T) {
	privKey, pubKey := setupTestKeys(t)
	signer := NewECDSASigner(privKey, pubKey)
	require.NotNil(t, signer)

	// Create a large message
	message := make([]byte, 1024*1024) // 1MB
	for i := range message {
		message[i] = byte(i % 256)
	}

	signature, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	valid, err := signer.Verify(context.Background(), message, signature, pubKey.SerializeCompressed())
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestECDSASigner_SignAndVerifyMultipleMessages(t *testing.T) {
	privKey, pubKey := setupTestKeys(t)
	signer := NewECDSASigner(privKey, pubKey)
	require.NotNil(t, signer)

	messages := [][]byte{
		[]byte("message 1"),
		[]byte("message 2"),
		[]byte("message 3"),
		[]byte("message 4"),
		[]byte("message 5"),
	}

	for _, msg := range messages {
		signature, err := signer.Sign(context.Background(), msg, types.ECDSA)
		require.NoError(t, err)

		valid, err := signer.Verify(context.Background(), msg, signature, pubKey.SerializeCompressed())
		require.NoError(t, err)
		assert.True(t, valid)

		// Verify with a different message should fail
		differentMsg := append(msg, []byte("modified")...)
		valid, err = signer.Verify(context.Background(), differentMsg, signature, pubKey.SerializeCompressed())
		require.NoError(t, err)
		assert.False(t, valid)
	}
}
