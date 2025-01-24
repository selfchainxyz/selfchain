package signing

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/stretchr/testify/require"

	"selfchain/x/keyless/crypto/signing/types"
)

type testCase struct {
	name       string
	networkID  string
	message    string
	algorithm  types.SigningAlgorithm
	wantErr    bool
	verifyFunc func(t *testing.T, signature *types.SignatureResult, message []byte) bool
}

func setupTSSKeyPair(t *testing.T) (*btcec.PrivateKey, *btcec.PublicKey) {
	// Generate a test private key
	privKey, err := btcec.NewPrivateKey()
	require.NoError(t, err)
	
	// Get the corresponding public key
	pubKey := privKey.PubKey()
	
	return privKey, pubKey
}

func TestMultiChainSigning(t *testing.T) {
	privKey, pubKey := setupTSSKeyPair(t)

	// Create signer factory
	signerFactory := NewSignerFactory()

	// Get key bytes
	pubKeyBytes := pubKey.SerializeCompressed()
	privKeyBytes := privKey.Serialize()

	// Create signer with both private and public keys
	signer, err := signerFactory.CreateSigner(context.Background(), types.ECDSA, privKeyBytes, pubKeyBytes)
	require.NoError(t, err)
	require.NotNil(t, signer)

	tests := []testCase{
		{
			name:      "Bitcoin Mainnet",
			networkID: "bitcoin:1",
			message:   "test bitcoin message",
			algorithm: types.ECDSA,
			wantErr:   false,
			verifyFunc: func(t *testing.T, signature *types.SignatureResult, message []byte) bool {
				// Bitcoin specific verification
				messageHash := sha256.Sum256(message)

				// Convert big.Int to ModNScalar
				var r, s btcec.ModNScalar
				r.SetByteSlice(signature.R.Bytes())
				s.SetByteSlice(signature.S.Bytes())

				sig := ecdsa.NewSignature(&r, &s)
				return sig.Verify(messageHash[:], pubKey)
			},
		},
		{
			name:      "Ethereum Mainnet",
			networkID: "ethereum:1",
			message:   "test ethereum message",
			algorithm: types.ECDSA,
			wantErr:   false,
			verifyFunc: func(t *testing.T, signature *types.SignatureResult, message []byte) bool {
				// Ethereum specific verification - hash the message directly
				messageHash := sha256.Sum256(message)

				// Convert big.Int to ModNScalar
				var r, s btcec.ModNScalar
				r.SetByteSlice(signature.R.Bytes())
				s.SetByteSlice(signature.S.Bytes())

				sig := ecdsa.NewSignature(&r, &s)
				return sig.Verify(messageHash[:], pubKey)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sign message
			signature, err := signer.Sign(context.Background(), []byte(tt.message), tt.algorithm)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, signature)

			// Verify signature
			valid := tt.verifyFunc(t, signature, []byte(tt.message))
			require.True(t, valid)
		})
	}
}

func TestMultiChainConcurrentSigning(t *testing.T) {
	privKey, pubKey := setupTSSKeyPair(t)

	// Create signer factory
	signerFactory := NewSignerFactory()

	// Get key bytes
	pubKeyBytes := pubKey.SerializeCompressed()
	privKeyBytes := privKey.Serialize()

	// Create signer with both private and public keys
	signer, err := signerFactory.CreateSigner(context.Background(), types.ECDSA, privKeyBytes, pubKeyBytes)
	require.NoError(t, err)
	require.NotNil(t, signer)

	// Run concurrent signing tests
	numGoroutines := 10
	numSignatures := 5
	errChan := make(chan error, numGoroutines*numSignatures)
	doneChan := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(routineID int) {
			for j := 0; j < numSignatures; j++ {
				message := []byte(fmt.Sprintf("test message %d-%d", routineID, j))
				signature, err := signer.Sign(context.Background(), message, types.ECDSA)
				if err != nil {
					errChan <- err
					continue
				}

				// Verify signature
				messageHash := sha256.Sum256(message)

				// Convert big.Int to ModNScalar
				var r, s btcec.ModNScalar
				r.SetByteSlice(signature.R.Bytes())
				s.SetByteSlice(signature.S.Bytes())

				sig := ecdsa.NewSignature(&r, &s)
				if !sig.Verify(messageHash[:], pubKey) {
					errChan <- fmt.Errorf("signature verification failed for message %d-%d", routineID, j)
				}
			}
			doneChan <- true
		}(i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < numGoroutines; i++ {
		<-doneChan
	}

	// Check for any errors
	close(errChan)
	for err := range errChan {
		require.NoError(t, err)
	}
}
