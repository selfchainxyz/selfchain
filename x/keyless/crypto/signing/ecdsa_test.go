package signing

import (
	"context"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/require"
	"selfchain/x/keyless/crypto/signing/ecdsa"
	"selfchain/x/keyless/crypto/signing/format"
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

		// Test with nil keys
		nilSigner := ecdsa.NewECDSASigner(nil, nil)
		require.NotNil(t, nilSigner)
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

		// Test with empty message
		_, err = signer.Sign(context.Background(), nil, types.ECDSA)
		require.Error(t, err)

		// Test with empty signature
		valid, err = signer.Verify(context.Background(), message, nil, pubKey.SerializeCompressed())
		require.Error(t, err)
		require.False(t, valid)

		// Test with invalid signature
		invalidSig := &types.SignatureResult{
			R:     big.NewInt(0),
			S:     big.NewInt(0),
			Bytes: []byte("invalid"),
		}
		valid, err = signer.Verify(context.Background(), message, invalidSig, pubKey.SerializeCompressed())
		require.NoError(t, err)
		require.False(t, valid)

		// Test with invalid public key format
		valid, err = signer.Verify(context.Background(), message, sig, []byte("invalid"))
		require.Error(t, err)
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

func TestECDSASigner_GetPublicKey(t *testing.T) {
	// Generate test keys
	privKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	pubKey := privKey.PubKey()

	t.Run("Test with valid public key", func(t *testing.T) {
		signer := ecdsa.NewECDSASigner(privKey, pubKey)
		require.NotNil(t, signer)

		// Test with ECDSA algorithm
		retrievedPubKey, err := signer.GetPublicKey(context.Background(), types.ECDSA)
		require.NoError(t, err)
		require.Equal(t, pubKey.SerializeCompressed(), retrievedPubKey)

		// Test with unsupported algorithm
		_, err = signer.GetPublicKey(context.Background(), types.EdDSA)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported algorithm")
	})

	t.Run("Test with nil public key", func(t *testing.T) {
		signer := ecdsa.NewECDSASigner(privKey, nil)
		require.NotNil(t, signer)

		// Should fail to get public key
		_, err := signer.GetPublicKey(context.Background(), types.ECDSA)
		require.Error(t, err)
		require.Contains(t, err.Error(), "public key not initialized")
	})
}

func TestPadTo32(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "empty input",
			input:    []byte{},
			expected: make([]byte, 32),
		},
		{
			name:     "input shorter than 32",
			input:    []byte{1, 2, 3},
			expected: append(make([]byte, 29), []byte{1, 2, 3}...),
		},
		{
			name:     "input exactly 32",
			input:    make([]byte, 32),
			expected: make([]byte, 32),
		},
		{
			name:     "input longer than 32",
			input:    make([]byte, 40),
			expected: make([]byte, 32),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padTo32(tt.input)
			require.Equal(t, 32, len(result))
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatDERSignature(t *testing.T) {
	// Generate test keys
	privKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	pubKey := privKey.PubKey()

	signer := ecdsa.NewECDSASigner(privKey, pubKey)
	require.NotNil(t, signer)

	// Test signing and formatting
	message := []byte("test message")
	typeSig, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	// Convert types.SignatureResult to format.SignatureResult
	sig := &format.SignatureResult{
		R:     typeSig.R,
		S:     typeSig.S,
		V:     typeSig.V,
		Bytes: typeSig.Bytes,
	}

	// Format signature to DER
	derSig, err := format.FormatDERSignature(sig)
	require.NoError(t, err)
	require.NotNil(t, derSig)

	// Test with nil signature
	_, err = format.FormatDERSignature(nil)
	require.Error(t, err)

	// Test UnmarshalDERSignature
	r, s, err := format.UnmarshalDERSignature(derSig)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.NotNil(t, s)

	// Test with nil DER signature
	_, _, err = format.UnmarshalDERSignature(nil)
	require.Error(t, err)

	// Test with invalid DER signature
	_, _, err = format.UnmarshalDERSignature([]byte("invalid"))
	require.Error(t, err)
}
