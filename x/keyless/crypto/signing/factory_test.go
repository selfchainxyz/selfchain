package signing

import (
	"context"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/require"
	"selfchain/x/keyless/crypto/signing/types"
)

func TestSignerFactory(t *testing.T) {
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

		// Test invalid message
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
		invalidPrivKey := []byte("invalid private key")
		_, err := factory.CreateSigner(context.Background(), types.ECDSA, invalidPrivKey, nil)
		require.Error(t, err)

		// Test with invalid public key
		invalidPubKey := []byte("invalid public key")
		_, err = factory.CreateSigner(context.Background(), types.ECDSA, nil, invalidPubKey)
		require.Error(t, err)

		// Test with no keys
		_, err = factory.CreateSigner(context.Background(), types.ECDSA, nil, nil)
		require.Error(t, err)
	})
}

// setupTestKeys is a helper function to generate test keys
func setupTestKeys(t *testing.T) (*btcec.PrivateKey, *btcec.PublicKey) {
	privKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	pubKey := privKey.PubKey()
	return privKey, pubKey
}

// Test_SignerFactory_Sign tests the Sign method
func Test_SignerFactory_Sign(t *testing.T) {
	ctx := context.Background()
	factory := NewSignerFactory()
	privKey, pubKey := setupTestKeys(t)

	tests := []struct {
		name      string
		message   []byte
		signer    types.SigningService
		wantError bool
	}{
		{
			name:      "successful signing",
			message:   []byte("test message"),
			signer:    setupSigner(t, privKey, pubKey),
			wantError: false,
		},
		{
			name:      "nil signer",
			message:   []byte("test message"),
			signer:    nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sig, err := factory.Sign(ctx, tt.message, types.ECDSA, tt.signer)
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, sig)
			} else {
				require.NoError(t, err)
				require.NotNil(t, sig)
				require.NotNil(t, sig.R)
				require.NotNil(t, sig.S)
				require.NotNil(t, sig.Bytes)
			}
		})
	}
}

// Test_SignerFactory_Verify tests the Verify method
func Test_SignerFactory_Verify(t *testing.T) {
	ctx := context.Background()
	factory := NewSignerFactory()
	privKey, pubKey := setupTestKeys(t)
	signer := setupSigner(t, privKey, pubKey)
	message := []byte("test message")

	// Generate a valid signature
	sig, err := signer.Sign(ctx, message, types.ECDSA)
	require.NoError(t, err)

	tests := []struct {
		name      string
		message   []byte
		sig       *types.SignatureResult
		pubKey    []byte
		signer    types.SigningService
		want      bool
		wantError bool
	}{
		{
			name:      "valid signature",
			message:   message,
			sig:       sig,
			pubKey:    pubKey.SerializeCompressed(),
			signer:    signer,
			want:      true,
			wantError: false,
		},
		{
			name:      "nil signer",
			message:   message,
			sig:       sig,
			pubKey:    pubKey.SerializeCompressed(),
			signer:    nil,
			want:      false,
			wantError: true,
		},
		{
			name:      "invalid message",
			message:   []byte("wrong message"),
			sig:       sig,
			pubKey:    pubKey.SerializeCompressed(),
			signer:    signer,
			want:      false,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := factory.Verify(ctx, tt.message, tt.sig, tt.pubKey, types.ECDSA, tt.signer)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

// Test_SignerFactory_FormatSignature tests the FormatSignature method
func Test_SignerFactory_FormatSignature(t *testing.T) {
	ctx := context.Background()
	factory := NewSignerFactory()
	privKey, pubKey := setupTestKeys(t)
	signer := setupSigner(t, privKey, pubKey)
	message := []byte("test message")

	// Generate a valid signature
	sig, err := signer.Sign(ctx, message, types.ECDSA)
	require.NoError(t, err)

	tests := []struct {
		name      string
		sig       *types.SignatureResult
		wantError bool
	}{
		{
			name:      "valid signature",
			sig:       sig,
			wantError: false,
		},
		{
			name:      "nil signature",
			sig:       nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted, err := factory.FormatSignature(ctx, tt.sig, types.ECDSA)
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, formatted)
			} else {
				require.NoError(t, err)
				require.NotNil(t, formatted)
			}
		})
	}
}

// Test_SignerFactory_UnformatSignature tests the UnformatSignature method
func Test_SignerFactory_UnformatSignature(t *testing.T) {
	ctx := context.Background()
	factory := NewSignerFactory()
	privKey, pubKey := setupTestKeys(t)
	signer := setupSigner(t, privKey, pubKey)
	message := []byte("test message")

	// Generate a valid signature
	sig, err := signer.Sign(ctx, message, types.ECDSA)
	require.NoError(t, err)

	// Format the signature
	formatted, err := factory.FormatSignature(ctx, sig, types.ECDSA)
	require.NoError(t, err)

	tests := []struct {
		name      string
		sigBytes  []byte
		wantError bool
	}{
		{
			name:      "valid signature bytes",
			sigBytes:  formatted,
			wantError: false,
		},
		{
			name:      "empty signature bytes",
			sigBytes:  []byte{},
			wantError: true,
		},
		{
			name:      "invalid signature bytes",
			sigBytes:  []byte("invalid signature"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unformatted, err := factory.UnformatSignature(ctx, tt.sigBytes, types.ECDSA)
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, unformatted)
			} else {
				require.NoError(t, err)
				require.NotNil(t, unformatted)
				require.NotNil(t, unformatted.R)
				require.NotNil(t, unformatted.S)
				require.NotNil(t, unformatted.Bytes)
			}
		})
	}
}

// setupSigner is a helper function to create a signer for testing
func setupSigner(t *testing.T, privKey *btcec.PrivateKey, pubKey *btcec.PublicKey) types.SigningService {
	factory := NewSignerFactory()
	signer, err := factory.createECDSASigner(privKey.Serialize(), pubKey.SerializeCompressed())
	require.NoError(t, err)
	return signer
}
