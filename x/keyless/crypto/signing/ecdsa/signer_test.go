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

func TestECDSASigner_Sign(t *testing.T) {
	signer := NewECDSASigner(nil, nil)
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
	signer := NewECDSASigner(nil, nil)
	require.NotNil(t, signer)

	message := []byte("test message")
	signature, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	pubKey, err := signer.GetPublicKey(context.Background(), types.ECDSA)
	require.NoError(t, err)

	valid, err := signer.Verify(context.Background(), message, signature, pubKey)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestECDSASigner_VerifyWithDifferentFormats(t *testing.T) {
	signer := NewECDSASigner(nil, nil)
	require.NotNil(t, signer)

	message := []byte("test message")
	signature, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	// Test with compressed public key
	pubKey := signer.pubKey.SerializeCompressed()
	valid, err := signer.Verify(context.Background(), message, signature, pubKey)
	require.NoError(t, err)
	assert.True(t, valid)

	// Test with uncompressed public key
	pubKey = signer.pubKey.SerializeUncompressed()
	valid, err = signer.Verify(context.Background(), message, signature, pubKey)
	require.NoError(t, err)
	assert.True(t, valid)

	// Test with raw public key
	pubKey = signer.pubKey.SerializeUncompressed()[1:] // Remove the marker byte
	valid, err = signer.Verify(context.Background(), message, signature, pubKey)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestECDSASigner_InvalidSignature(t *testing.T) {
	signer := NewECDSASigner(nil, nil)
	require.NotNil(t, signer)

	message := []byte("test message")
	signature, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	pubKey, err := signer.GetPublicKey(context.Background(), types.ECDSA)
	require.NoError(t, err)

	// Modify the signature to make it invalid
	signature.Bytes[0] ^= 0xFF

	valid, err := signer.Verify(context.Background(), message, signature, pubKey)
	require.Error(t, err)
	assert.False(t, valid)
}

func TestECDSASigner_InvalidPublicKey(t *testing.T) {
	signer := NewECDSASigner(nil, nil)
	require.NotNil(t, signer)

	message := []byte("test message")
	signature, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	// Test with invalid public key length
	_, err = signer.Verify(context.Background(), message, signature, []byte{0x00})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid public key length")

	// Test with invalid public key format
	_, err = signer.Verify(context.Background(), message, signature, make([]byte, 33))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse public key")
}

func TestECDSASigner_KnownTestVector(t *testing.T) {
	// Create a known private key
	privKeyBytes, _ := hex.DecodeString("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	signer := &ECDSASigner{
		privKey: privKey,
		pubKey:  privKey.PubKey(),
	}

	message := []byte("test message")
	signature, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	// Verify the signature with the known public key
	pubKey := privKey.PubKey().SerializeCompressed()
	valid, err := signer.Verify(context.Background(), message, signature, pubKey)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestECDSASigner_SignAndVerifyLargeMessage(t *testing.T) {
	signer := NewECDSASigner(nil, nil)
	require.NotNil(t, signer)

	// Create a large message
	message := make([]byte, 1024*1024) // 1MB
	for i := range message {
		message[i] = byte(i % 256)
	}

	signature, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	pubKey, err := signer.GetPublicKey(context.Background(), types.ECDSA)
	require.NoError(t, err)

	valid, err := signer.Verify(context.Background(), message, signature, pubKey)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestECDSASigner_SignAndVerifyMultipleMessages(t *testing.T) {
	signer := NewECDSASigner(nil, nil)
	require.NotNil(t, signer)

	messages := [][]byte{
		[]byte("message 1"),
		[]byte("message 2"),
		[]byte("message 3"),
		[]byte("message 4"),
		[]byte("message 5"),
	}

	pubKey, err := signer.GetPublicKey(context.Background(), types.ECDSA)
	require.NoError(t, err)

	for _, message := range messages {
		signature, err := signer.Sign(context.Background(), message, types.ECDSA)
		require.NoError(t, err)

		valid, err := signer.Verify(context.Background(), message, signature, pubKey)
		require.NoError(t, err)
		assert.True(t, valid)

		// Verify signature fails with different message
		wrongMessage := append(message, []byte("wrong")...)
		valid, err = signer.Verify(context.Background(), wrongMessage, signature, pubKey)
		require.NoError(t, err)
		assert.False(t, valid)
	}
}
