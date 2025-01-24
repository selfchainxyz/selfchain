package signing

import (
	"context"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/require"
	"selfchain/x/keyless/crypto/signing/ecdsa"
	"selfchain/x/keyless/crypto/signing/types"
)

func TestECDSASigner(t *testing.T) {
	// Generate test keys
	privKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	pubKey := privKey.PubKey()

	t.Run("Test NewECDSASigner", func(t *testing.T) {
		signer := ecdsa.NewECDSASigner(privKey, pubKey)
		require.NotNil(t, signer)
		
		// Test public key retrieval
		retrievedPubKey, err := signer.GetPublicKey(context.Background(), types.ECDSA)
		require.NoError(t, err)
		require.Equal(t, pubKey.SerializeCompressed(), retrievedPubKey)
	})

	t.Run("Test Sign and Verify", func(t *testing.T) {
		signer := ecdsa.NewECDSASigner(privKey, pubKey)
		require.NotNil(t, signer)

		// Test signing
		message := []byte("test message")
		sig, err := signer.Sign(context.Background(), message, types.ECDSA)
		require.NoError(t, err)
		require.NotNil(t, sig)

		// Test verification
		valid, err := signer.Verify(context.Background(), message, sig, pubKey.SerializeCompressed())
		require.NoError(t, err)
		require.True(t, valid)

		// Test invalid message
		invalidMessage := []byte("wrong message")
		valid, err = signer.Verify(context.Background(), invalidMessage, sig, pubKey.SerializeCompressed())
		require.NoError(t, err)
		require.False(t, valid)
	})

	t.Run("Test Invalid Keys", func(t *testing.T) {
		// Test with nil private key
		signer := ecdsa.NewECDSASigner(nil, pubKey)
		require.NotNil(t, signer)

		// Should fail to sign without private key
		message := []byte("test message")
		_, err := signer.Sign(context.Background(), message, types.ECDSA)
		require.Error(t, err)

		// Test with nil public key
		signer = ecdsa.NewECDSASigner(privKey, nil)
		require.NotNil(t, signer)

		// Should fail to verify without public key
		sig, err := signer.Sign(context.Background(), message, types.ECDSA)
		require.NoError(t, err)
		_, err = signer.Verify(context.Background(), message, sig, nil)
		require.Error(t, err)
	})
}
