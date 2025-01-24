package signing

import (
	"context"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/require"
	"selfchain/x/keyless/crypto/signing/types"
)

func TestSignerFactory(t *testing.T) {
	// Create a new signer factory
	factory := NewSignerFactory()
	require.NotNil(t, factory)

	// Generate a test private key
	privKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	require.NotNil(t, privKey)

	t.Run("Test_CreateSigner", func(t *testing.T) {
		// Create signer with private key
		signer, err := factory.CreateSigner(context.Background(), types.ECDSA, privKey.Serialize(), privKey.PubKey().SerializeCompressed())
		require.NoError(t, err)
		require.NotNil(t, signer)

		// Test signing
		message := []byte("test message")
		sig, err := signer.Sign(context.Background(), message, types.ECDSA)
		require.NoError(t, err)
		require.NotNil(t, sig)
		require.NotNil(t, sig.Bytes)

		// Test verification
		valid, err := signer.Verify(context.Background(), message, sig, privKey.PubKey().SerializeCompressed())
		require.NoError(t, err)
		require.True(t, valid)

		// Test invalid signature
		invalidMessage := []byte("wrong message")
		valid, err = signer.Verify(context.Background(), invalidMessage, sig, privKey.PubKey().SerializeCompressed())
		require.NoError(t, err)
		require.False(t, valid)

		// Test unsupported algorithm
		_, err = factory.CreateSigner(context.Background(), "unsupported", nil, nil)
		require.Error(t, err)
	})

	t.Run("Test_InvalidKeys", func(t *testing.T) {
		// Test with invalid private key
		_, err := factory.CreateSigner(context.Background(), types.ECDSA, []byte("invalid"), nil)
		require.Error(t, err)

		// Test with invalid public key
		_, err = factory.CreateSigner(context.Background(), types.ECDSA, nil, []byte("invalid"))
		require.Error(t, err)

		// Test with no keys
		_, err = factory.CreateSigner(context.Background(), types.ECDSA, nil, nil)
		require.Error(t, err)
	})
}
