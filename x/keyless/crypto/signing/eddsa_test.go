package signing

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEdDSASigner(t *testing.T) {
	// Generate test key pair
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	t.Run("Test_Sign", func(t *testing.T) {
		// Create signer with test key
		signer := &EdDSASigner{
			privKey: privKey,
			pubKey:  pubKey,
		}
		require.NotNil(t, signer)

		// Test signing
		message := []byte("test message")
		sig, err := signer.Sign(context.Background(), message, EdDSA)
		require.NoError(t, err)
		require.NotNil(t, sig)

		// Verify signature using standard ed25519 package
		require.True(t, ed25519.Verify(pubKey, message, sig.Bytes))
	})

	t.Run("Test_Verify", func(t *testing.T) {
		// Create signer with test key
		signer := &EdDSASigner{
			privKey: privKey,
			pubKey:  pubKey,
		}
		require.NotNil(t, signer)

		// Test message
		message := []byte("test message")

		// Sign message using standard ed25519 package
		signature := ed25519.Sign(privKey, message)

		// Create SignatureResult
		sigResult := &SignatureResult{
			Bytes: signature,
		}

		// Verify using our signer
		valid, err := signer.Verify(context.Background(), message, sigResult, pubKey)
		require.NoError(t, err)
		require.True(t, valid)

		// Test invalid signature
		invalidSig := make([]byte, len(signature))
		copy(invalidSig, signature)
		invalidSig[0] ^= 0xFF // Flip some bits
		sigResult.Bytes = invalidSig
		valid, err = signer.Verify(context.Background(), message, sigResult, pubKey)
		require.NoError(t, err)
		require.False(t, valid)

		// Test with nil signature
		valid, err = signer.Verify(context.Background(), message, nil, pubKey)
		require.Error(t, err)
		require.False(t, valid)

		// Test with nil public key
		valid, err = signer.Verify(context.Background(), message, sigResult, nil)
		require.Error(t, err)
		require.False(t, valid)
	})

	t.Run("Test_GetPublicKey", func(t *testing.T) {
		// Create signer with test key
		signer := &EdDSASigner{
			privKey: privKey,
			pubKey:  pubKey,
		}
		require.NotNil(t, signer)

		// Get public key
		pubKeyBytes, err := signer.GetPublicKey(context.Background(), EdDSA)
		require.NoError(t, err)
		require.NotNil(t, pubKeyBytes)
		require.Equal(t, ed25519.PublicKeySize, len(pubKeyBytes))
		require.Equal(t, []byte(pubKey), pubKeyBytes)
	})

	t.Run("Test_with_nil_party_data", func(t *testing.T) {
		// Create signer with nil party data
		signer := NewEdDSASigner(nil, nil)
		require.NotNil(t, signer)

		// Get public key
		pubKeyBytes, err := signer.GetPublicKey(context.Background(), EdDSA)
		require.NoError(t, err)
		require.NotNil(t, pubKeyBytes)
		require.Equal(t, ed25519.PublicKeySize, len(pubKeyBytes))
	})
}
