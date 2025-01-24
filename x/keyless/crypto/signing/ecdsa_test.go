package signing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestECDSASigner(t *testing.T) {
	// Create test party data
	party1Data := &struct{ Key string }{"party1"}
	party2Data := &struct{ Key string }{"party2"}

	t.Run("Test NewECDSASigner", func(t *testing.T) {
		signer := NewECDSASigner(party1Data, party2Data)
		require.NotNil(t, signer)
		require.Equal(t, party1Data, signer.party1Data)
		require.Equal(t, party2Data, signer.party2Data)
	})

	t.Run("Test Sign", func(t *testing.T) {
		ctx := context.Background()
		signer := NewECDSASigner(party1Data, party2Data)
		message := []byte("test message")

		// Test ECDSA signing
		signature, err := signer.Sign(ctx, message, ECDSA)
		require.NoError(t, err)
		require.NotNil(t, signature)
		require.NotNil(t, signature.R)
		require.NotNil(t, signature.S)

		// Test unsupported algorithm
		signature, err = signer.Sign(ctx, message, "unsupported")
		require.Error(t, err)
		require.Nil(t, signature)
	})

	t.Run("Test Verify", func(t *testing.T) {
		ctx := context.Background()
		signer := NewECDSASigner(party1Data, party2Data)
		message := []byte("test message")

		// Generate a test signature first
		signature, err := signer.Sign(ctx, message, ECDSA)
		require.NoError(t, err)
		require.NotNil(t, signature)

		// Get the public key
		pubKey, err := signer.GetPublicKey(ctx, ECDSA)
		require.NoError(t, err)
		require.NotNil(t, pubKey)

		// Test signature verification
		valid, err := signer.Verify(ctx, message, signature, pubKey)
		require.NoError(t, err)
		require.True(t, valid) // Should be valid since we're using our own signature

		// Test with modified message
		modifiedMsg := append([]byte(nil), message...)
		modifiedMsg[0] ^= 0xFF // Flip some bits
		valid, err = signer.Verify(ctx, modifiedMsg, signature, pubKey)
		require.NoError(t, err)
		require.False(t, valid)

		// Test with nil signature
		valid, err = signer.Verify(ctx, message, nil, pubKey)
		require.Error(t, err)
		require.False(t, valid)

		// Test with nil public key
		valid, err = signer.Verify(ctx, message, signature, nil)
		require.Error(t, err)
		require.False(t, valid)
	})

	t.Run("Test GetPublicKey", func(t *testing.T) {
		ctx := context.Background()
		signer := NewECDSASigner(party1Data, party2Data)

		// Test ECDSA public key retrieval
		pubKey, err := signer.GetPublicKey(ctx, ECDSA)
		require.NoError(t, err)
		require.NotNil(t, pubKey)

		// Test unsupported algorithm
		pubKey, err = signer.GetPublicKey(ctx, "unsupported")
		require.Error(t, err)
		require.Nil(t, pubKey)
	})

	t.Run("Test with nil party data", func(t *testing.T) {
		ctx := context.Background()
		signer := NewECDSASigner(nil, nil)
		message := []byte("test message")

		// Test Sign
		signature, err := signer.Sign(ctx, message, ECDSA)
		require.NoError(t, err) // Current implementation returns dummy values
		require.NotNil(t, signature)

		// Test GetPublicKey
		pubKey, err := signer.GetPublicKey(ctx, ECDSA)
		require.NoError(t, err) // Current implementation returns dummy values
		require.NotNil(t, pubKey)
	})
}
