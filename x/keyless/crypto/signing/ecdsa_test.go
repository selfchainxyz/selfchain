package signing

import (
	"context"
	"bytes"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/require"
	"selfchain/x/keyless/crypto/signing/format"
	"selfchain/x/keyless/crypto/signing/ecdsa"
	"selfchain/x/keyless/crypto/signing/types"
)

func TestECDSAPackage(t *testing.T) {
	// Generate test keys
	privKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	pubKey := privKey.PubKey()

	t.Run("Test NewECDSASigner", func(t *testing.T) {
		signer := ecdsa.NewECDSASigner(privKey, pubKey)
		require.NotNil(t, signer)

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

		// Test invalid algorithm
		_, err = signer.Sign(context.Background(), message, types.EdDSA)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported algorithm")

		// Test nil private key
		nilSigner := ecdsa.NewECDSASigner(nil, pubKey)
		_, err = nilSigner.Sign(context.Background(), message, types.ECDSA)
		require.Error(t, err)
		require.Contains(t, err.Error(), "private key not initialized")

		// Test empty message
		_, err = signer.Sign(context.Background(), []byte{}, types.ECDSA)
		require.Error(t, err)
		require.Contains(t, err.Error(), "empty message")
	})

	t.Run("Test GetPublicKey", func(t *testing.T) {
		signer := ecdsa.NewECDSASigner(privKey, pubKey)
		require.NotNil(t, signer)

		// Test with ECDSA algorithm
		retrievedPubKey, err := signer.GetPublicKey(context.Background(), types.ECDSA)
		require.NoError(t, err)
		require.Equal(t, pubKey.SerializeCompressed(), retrievedPubKey)

		// Test with invalid algorithm
		_, err = signer.GetPublicKey(context.Background(), types.EdDSA)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported algorithm")

		// Test with nil public key
		nilSigner := ecdsa.NewECDSASigner(privKey, nil)
		_, err = nilSigner.GetPublicKey(context.Background(), types.ECDSA)
		require.Error(t, err)
		require.Contains(t, err.Error(), "public key not initialized")
	})
}

func TestSignatureVerification(t *testing.T) {
	// Generate test keys
	privKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	pubKey := privKey.PubKey()
	signer := ecdsa.NewECDSASigner(privKey, pubKey)

	t.Run("Test basic signature verification", func(t *testing.T) {
		message := []byte("test message")
		sig, err := signer.Sign(context.Background(), message, types.ECDSA)
		require.NoError(t, err)
		require.NotNil(t, sig)

		// Verify valid signature
		valid, err := signer.Verify(context.Background(), message, sig, pubKey.SerializeCompressed())
		require.NoError(t, err)
		require.True(t, valid)

		// Verify with modified message
		modifiedMsg := []byte("modified message")
		valid, err = signer.Verify(context.Background(), modifiedMsg, sig, pubKey.SerializeCompressed())
		require.NoError(t, err)
		require.False(t, valid)
	})

	t.Run("Test invalid signature verification", func(t *testing.T) {
		message := []byte("test message")
		sig, err := signer.Sign(context.Background(), message, types.ECDSA)
		require.NoError(t, err)

		// Test with nil signature
		_, err = signer.Verify(context.Background(), message, nil, pubKey.SerializeCompressed())
		require.Error(t, err)
		require.Contains(t, err.Error(), "signature is nil")

		// Test with zero values
		zeroSig := &types.SignatureResult{
			R: big.NewInt(0),
			S: big.NewInt(0),
		}
		valid, err := signer.Verify(context.Background(), message, zeroSig, pubKey.SerializeCompressed())
		require.NoError(t, err)
		require.False(t, valid)

		// Test with very large values
		largeSig := &types.SignatureResult{
			R: new(big.Int).Lsh(big.NewInt(1), 257),
			S: new(big.Int).Lsh(big.NewInt(1), 257),
		}
		valid, err = signer.Verify(context.Background(), message, largeSig, pubKey.SerializeCompressed())
		require.NoError(t, err)
		require.False(t, valid)

		// Test with invalid public key
		_, err = signer.Verify(context.Background(), message, sig, []byte{0x00})
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse public key")
	})

	t.Run("Test signature malleability", func(t *testing.T) {
		message := []byte("test message")
		sig, err := signer.Sign(context.Background(), message, types.ECDSA)
		require.NoError(t, err)

		// Modify S value to test malleability
		curve := btcec.S256()
		n := curve.N
		modifiedS := new(big.Int).Sub(n, sig.S)
		
		modifiedSig := &types.SignatureResult{
			R: sig.R,
			S: modifiedS,
			Bytes: sig.Bytes,
		}

		// Both signatures should verify
		valid, err := signer.Verify(context.Background(), message, sig, pubKey.SerializeCompressed())
		require.NoError(t, err)
		require.True(t, valid)

		valid, err = signer.Verify(context.Background(), message, modifiedSig, pubKey.SerializeCompressed())
		require.NoError(t, err)
		require.True(t, valid)
	})
}

func TestSignatureEdgeCases(t *testing.T) {
	// Generate test keys
	privKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	pubKey := privKey.PubKey()
	signer := ecdsa.NewECDSASigner(privKey, pubKey)

	t.Run("Test with max message size", func(t *testing.T) {
		// Create a large message
		largeMessage := make([]byte, 1024*1024) // 1MB
		for i := range largeMessage {
			largeMessage[i] = byte(i % 256)
		}

		sig, err := signer.Sign(context.Background(), largeMessage, types.ECDSA)
		require.NoError(t, err)
		require.NotNil(t, sig)

		valid, err := signer.Verify(context.Background(), largeMessage, sig, pubKey.SerializeCompressed())
		require.NoError(t, err)
		require.True(t, valid)
	})

	t.Run("Test with special characters", func(t *testing.T) {
		specialMessage := []byte("!@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`")
		sig, err := signer.Sign(context.Background(), specialMessage, types.ECDSA)
		require.NoError(t, err)
		require.NotNil(t, sig)

		valid, err := signer.Verify(context.Background(), specialMessage, sig, pubKey.SerializeCompressed())
		require.NoError(t, err)
		require.True(t, valid)
	})

	t.Run("Test with repeated signing", func(t *testing.T) {
		message := []byte("test message")
		// Sign the same message multiple times
		for i := 0; i < 10; i++ {
			sig, err := signer.Sign(context.Background(), message, types.ECDSA)
			require.NoError(t, err)
			require.NotNil(t, sig)

			valid, err := signer.Verify(context.Background(), message, sig, pubKey.SerializeCompressed())
			require.NoError(t, err)
			require.True(t, valid)
		}
	})
}

func TestDERSignatureFormatting(t *testing.T) {
	// Generate test keys and signature
	privKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	pubKey := privKey.PubKey()
	signer := ecdsa.NewECDSASigner(privKey, pubKey)
	
	message := []byte("test message")
	sig, err := signer.Sign(context.Background(), message, types.ECDSA)
	require.NoError(t, err)

	t.Run("Test DER formatting", func(t *testing.T) {
		// Test formatting to DER
		derSig := sig.Bytes
		require.NotNil(t, derSig)
		require.Equal(t, byte(0x30), derSig[0]) // DER sequence identifier

		// Test unmarshalling DER and verify
		r, s, err := format.UnmarshalDERSignature(derSig)
		require.NoError(t, err)
		require.NotNil(t, r)
		require.NotNil(t, s)

		// Verify that the unmarshalled signature is valid
		verifyResult, err := signer.Verify(context.Background(), message, &types.SignatureResult{
			R: r,
			S: s,
			Bytes: derSig,
		}, pubKey.SerializeCompressed())
		require.NoError(t, err)
		require.True(t, verifyResult)
	})

	t.Run("Test invalid DER unmarshalling", func(t *testing.T) {
		// Test with invalid DER format
		_, _, err := format.UnmarshalDERSignature([]byte{0x00})
		require.Error(t, err)

		// Test with empty input
		_, _, err = format.UnmarshalDERSignature([]byte{})
		require.Error(t, err)

		// Test with nil input
		_, _, err = format.UnmarshalDERSignature(nil)
		require.Error(t, err)

		// Test with invalid sequence
		_, _, err = format.UnmarshalDERSignature([]byte{0x31, 0x00}) // 0x31 is not the DER sequence identifier
		require.Error(t, err)

		// Test with truncated DER
		_, _, err = format.UnmarshalDERSignature([]byte{0x30, 0x06, 0x02, 0x01}) // Truncated DER
		require.Error(t, err)

		// Test with invalid length
		_, _, err = format.UnmarshalDERSignature([]byte{0x30, 0x80}) // Invalid length
		require.Error(t, err)
	})
}

func TestGetPublicKey(t *testing.T) {
	// Generate test keys
	privKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	pubKey := privKey.PubKey()
	signer := ecdsa.NewECDSASigner(privKey, pubKey)

	t.Run("Test GetPublicKey with different contexts", func(t *testing.T) {
		// Test with background context
		pubKeyBytes, err := signer.GetPublicKey(context.Background(), types.ECDSA)
		require.NoError(t, err)
		require.Equal(t, pubKey.SerializeCompressed(), pubKeyBytes)

		// Test with nil context
		pubKeyBytes, err = signer.GetPublicKey(nil, types.ECDSA)
		require.NoError(t, err)
		require.Equal(t, pubKey.SerializeCompressed(), pubKeyBytes)

		// Test with canceled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		pubKeyBytes, err = signer.GetPublicKey(ctx, types.ECDSA)
		require.NoError(t, err)
		require.Equal(t, pubKey.SerializeCompressed(), pubKeyBytes)
	})

	t.Run("Test GetPublicKey error cases", func(t *testing.T) {
		// Test with nil public key
		nilSigner := ecdsa.NewECDSASigner(privKey, nil)
		require.NotNil(t, nilSigner)

		_, err := nilSigner.GetPublicKey(context.Background(), types.ECDSA)
		require.Error(t, err)
		require.Contains(t, err.Error(), "public key not initialized")

		// Test with unsupported algorithm
		_, err = signer.GetPublicKey(context.Background(), types.EdDSA)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported algorithm")
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
			input:    bytes.Repeat([]byte{1}, 32),
			expected: bytes.Repeat([]byte{1}, 32),
		},
		{
			name:     "input longer than 32",
			input:    bytes.Repeat([]byte{1}, 40),
			expected: bytes.Repeat([]byte{1}, 32),
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

	t.Run("Test DER formatting success", func(t *testing.T) {
		message := []byte("test message")
		sig, err := signer.Sign(context.Background(), message, types.ECDSA)
		require.NoError(t, err)

		// Convert to format.SignatureResult
		formatSig := &format.SignatureResult{
			R: sig.R,
			S: sig.S,
		}

		// Format signature to DER
		derSig, err := format.FormatDERSignature(formatSig)
		require.NoError(t, err)
		require.NotNil(t, derSig)
		require.Equal(t, byte(0x30), derSig[0]) // DER sequence identifier

		// Test with high R value
		highR := new(big.Int).SetBytes(bytes.Repeat([]byte{0xff}, 31))
		highRSig := &format.SignatureResult{
			R: highR,
			S: sig.S,
		}
		derSig, err = format.FormatDERSignature(highRSig)
		require.NoError(t, err)
		require.NotNil(t, derSig)

		// Test with high S value
		highS := new(big.Int).SetBytes(bytes.Repeat([]byte{0xff}, 31))
		highSSig := &format.SignatureResult{
			R: sig.R,
			S: highS,
		}
		derSig, err = format.FormatDERSignature(highSSig)
		require.NoError(t, err)
		require.NotNil(t, derSig)
	})

	t.Run("Test DER formatting errors", func(t *testing.T) {
		// Test with nil signature
		derSig, err := format.FormatDERSignature(nil)
		require.Error(t, err)
		require.Nil(t, derSig)
		require.Contains(t, err.Error(), "signature is nil")

		// Test with nil R value
		derSig, err = format.FormatDERSignature(&format.SignatureResult{
			R: nil,
			S: big.NewInt(1),
		})
		require.Error(t, err)
		require.Nil(t, derSig)
		require.Contains(t, err.Error(), "R value is nil")

		// Test with nil S value
		derSig, err = format.FormatDERSignature(&format.SignatureResult{
			R: big.NewInt(1),
			S: nil,
		})
		require.Error(t, err)
		require.Nil(t, derSig)
		require.Contains(t, err.Error(), "S value is nil")

		// Test with zero values
		derSig, err = format.FormatDERSignature(&format.SignatureResult{
			R: big.NewInt(0),
			S: big.NewInt(0),
		})
		require.Error(t, err)
		require.Nil(t, derSig)
		require.Contains(t, err.Error(), "invalid signature")

		// Test with very large values
		largeNum := new(big.Int).Lsh(big.NewInt(1), 257)
		derSig, err = format.FormatDERSignature(&format.SignatureResult{
			R: largeNum,
			S: largeNum,
		})
		require.Error(t, err)
		require.Nil(t, derSig)
		require.Contains(t, err.Error(), "invalid signature")
	})
}
