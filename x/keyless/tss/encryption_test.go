package tss

import (
	"context"
	"testing"
	"time"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/stretchr/testify/require"
	"selfchain/x/keyless/crypto"
)

func TestEncryptDecryptShare(t *testing.T) {
	// Generate a key share first
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	preParams, err := keygen.GeneratePreParams(time.Minute)
	require.NoError(t, err, "Failed to generate pre-parameters")

	result, err := GenerateKey(ctx, preParams, "test-chain-1")
	require.NoError(t, err, "Failed to generate key shares")
	require.NotNil(t, result.Party1Data, "Party1Data should not be nil")

	tests := []struct {
		name      string
		shareData *keygen.LocalPartySaveData
		wantErr   bool
	}{
		{
			name:      "Encrypt and decrypt Party1 data",
			shareData: result.Party1Data,
			wantErr:   false,
		},
		{
			name:      "Encrypt and decrypt Party2 data",
			shareData: result.Party2Data,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate encryption key
			key, err := crypto.NewEncryptionKey()
			require.NoError(t, err, "Failed to generate encryption key")

			// Encrypt the share
			encryptedShare, err := EncryptShare(key, tt.shareData)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptShare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NotEmpty(t, encryptedShare.EncryptedData, "Encrypted data should not be empty")

			// Decrypt the share
			decryptedShare, err := DecryptShare(key, encryptedShare)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptShare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify the decrypted share
			require.Equal(t, tt.shareData.ShareID, decryptedShare.ShareID, "ShareID mismatch")
			require.NotNil(t, decryptedShare.ECDSAPub, "ECDSAPub should not be nil")
			require.NotNil(t, decryptedShare.Xi, "Xi should not be nil")
		})
	}
}

func TestEncryptShareErrors(t *testing.T) {
	// Generate a key share first
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	preParams, err := keygen.GeneratePreParams(time.Minute)
	require.NoError(t, err, "Failed to generate pre-parameters")

	result, err := GenerateKey(ctx, preParams, "test-chain-1")
	require.NoError(t, err, "Failed to generate key shares")

	tests := []struct {
		name      string
		key       crypto.EncryptionKey
		shareData *keygen.LocalPartySaveData
		wantErr   bool
	}{
		{
			name:      "Invalid key length",
			key:       make([]byte, 16), // Too short for AES-256
			shareData: result.Party1Data,
			wantErr:   true,
		},
		{
			name:      "Nil share data",
			key:       make([]byte, 32),
			shareData: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EncryptShare(tt.key, tt.shareData)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptShare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecryptShareErrors(t *testing.T) {
	key, err := crypto.NewEncryptionKey()
	require.NoError(t, err, "Failed to generate encryption key")

	tests := []struct {
		name          string
		encryptedData string
		chainID       string
		wantErr       bool
	}{
		{
			name:          "Invalid encrypted data",
			encryptedData: "invalid-base64-data",
			chainID:       "test-chain-1",
			wantErr:       true,
		},
		{
			name:          "Empty encrypted data",
			encryptedData: "",
			chainID:       "test-chain-1",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptedShare := &EncryptedShare{
				EncryptedData: tt.encryptedData,
				ChainID:      tt.chainID,
			}

			_, err := DecryptShare(key, encryptedShare)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptShare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
