package eddsa

import (
	"crypto/ed25519"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/require"
)

func TestConvertECDSAToEdDSA(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *btcec.PublicKey
		wantErr bool
	}{
		{
			name: "Valid conversion",
			setup: func() *btcec.PublicKey {
				privKey, err := btcec.NewPrivateKey()
				require.NoError(t, err)
				return privKey.PubKey()
			},
			wantErr: false,
		},
		{
			name: "Nil public key",
			setup: func() *btcec.PublicKey {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ecdsaPub := tt.setup()
			result, err := ConvertECDSAToEdDSA(ecdsaPub)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			require.NotNil(t, result.EdDSAPublicKey)
			require.Equal(t, ed25519.PublicKeySize, len(result.EdDSAPublicKey))
			require.NotNil(t, result.ConversionProof)
		})
	}
}

func TestConvertSignatureToEdDSA(t *testing.T) {
	maxInt := new(big.Int).Lsh(big.NewInt(1), 256)

	tests := []struct {
		name    string
		r       *big.Int
		s       *big.Int
		wantErr bool
	}{
		{
			name:    "Valid signature components",
			r:       new(big.Int).SetInt64(123),
			s:       new(big.Int).SetInt64(456),
			wantErr: false,
		},
		{
			name:    "Nil R component",
			r:       nil,
			s:       new(big.Int).SetInt64(456),
			wantErr: true,
		},
		{
			name:    "Nil S component",
			r:       new(big.Int).SetInt64(123),
			s:       nil,
			wantErr: true,
		},
		{
			name:    "R component too large",
			r:       maxInt,
			s:       new(big.Int).SetInt64(456),
			wantErr: true,
		},
		{
			name:    "S component too large",
			r:       new(big.Int).SetInt64(123),
			s:       maxInt,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sig, err := ConvertSignatureToEdDSA(tt.r, tt.s)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, ed25519.SignatureSize, len(sig))
		})
	}
}

func TestVerifyEdDSASignature(t *testing.T) {
	// Generate Ed25519 keypair
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	message := []byte("test message")
	signature := ed25519.Sign(priv, message)

	tests := []struct {
		name      string
		pubKey    ed25519.PublicKey
		message   []byte
		signature []byte
		want      bool
	}{
		{
			name:      "Valid signature",
			pubKey:    pub,
			message:   message,
			signature: signature,
			want:      true,
		},
		{
			name:      "Invalid message",
			pubKey:    pub,
			message:   []byte("wrong message"),
			signature: signature,
			want:      false,
		},
		{
			name:      "Invalid signature",
			pubKey:    pub,
			message:   message,
			signature: make([]byte, ed25519.SignatureSize),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VerifyEdDSASignature(tt.pubKey, tt.message, tt.signature)
			require.Equal(t, tt.want, result)
		})
	}
}

func TestEndToEndConversion(t *testing.T) {
	// Generate ECDSA key
	ecdsaPriv, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	ecdsaPub := ecdsaPriv.PubKey()

	// Convert to EdDSA
	conversion, err := ConvertECDSAToEdDSA(ecdsaPub)
	require.NoError(t, err)

	// Create test signature components
	r := new(big.Int).SetInt64(123)
	s := new(big.Int).SetInt64(456)

	// Convert signature
	eddsaSig, err := ConvertSignatureToEdDSA(r, s)
	require.NoError(t, err)

	// Verify signature size
	require.Equal(t, ed25519.SignatureSize, len(eddsaSig))

	// Verify conversion parameters
	require.Equal(t, ecdsaPub, conversion.ECDSAPublicKey)
	require.Equal(t, ed25519.PublicKeySize, len(conversion.EdDSAPublicKey))
	require.NotNil(t, conversion.ConversionProof)
}
