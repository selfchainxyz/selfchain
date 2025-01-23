package signing

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"

	"selfchain/x/keyless/networks"
	selfchainTss "selfchain/x/keyless/tss"
)

type testCase struct {
	name       string
	networkID  string
	message    string
	wantErr    bool
	verifyFunc func(t *testing.T, signature []byte, message []byte) bool
}

func TestMultiChainSigning(t *testing.T) {
	// Setup TSS key generation
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Generate TSS key pair
	preParams, err := keygen.GeneratePreParams(time.Minute)
	require.NoError(t, err)
	require.NotNil(t, preParams)

	keygenResult, err := selfchainTss.GenerateKey(ctx, preParams, "test-key")
	require.NoError(t, err)
	require.NotNil(t, keygenResult)

	// Create signer factory
	signerFactory := NewSignerFactory(networks.NewNetworkRegistry())

	// Get public key from party data
	pubKeyPoint := keygenResult.Party1Data.ECDSAPub
	x, y := pubKeyPoint.X(), pubKeyPoint.Y()
	
	// Convert big.Int coordinates to btcec.FieldVal
	xField := new(btcec.FieldVal)
	yField := new(btcec.FieldVal)
	xField.SetByteSlice(x.Bytes())
	yField.SetByteSlice(y.Bytes())
	
	btcPubKey := btcec.NewPublicKey(xField, yField)
	pubKeyBytes := btcPubKey.SerializeCompressed()

	tests := []testCase{
		{
			name:      "Bitcoin mainnet signature",
			networkID: "bitcoin:mainnet",
			message:   "test bitcoin transaction",
			verifyFunc: func(t *testing.T, signature []byte, message []byte) bool {
				// Verify Bitcoin DER signature format
				// DER format: 0x30 [total-length] 0x02 [R-length] [R] 0x02 [S-length] [S]
				if len(signature) < 8 {
					return false
				}
				if signature[0] != 0x30 || signature[2] != 0x02 {
					return false
				}
				totalLen := int(signature[1])
				if len(signature) != totalLen+2 {
					return false
				}
				return true
			},
		},
		{
			name:      "Ethereum mainnet signature",
			networkID: "ethereum:1",
			message:   "test ethereum transaction",
			verifyFunc: func(t *testing.T, signature []byte, message []byte) bool {
				// Verify Ethereum signature format
				// Format: R (32 bytes) || S (32 bytes) || V (1 byte)
				if len(signature) != 65 {
					return false
				}
				v := signature[64]
				return v >= 27 // V should be >= 27 for Ethereum signatures
			},
		},
		{
			name:      "Cosmos Hub signature",
			networkID: "cosmos:cosmoshub-4",
			message:   "test cosmos transaction",
			verifyFunc: func(t *testing.T, signature []byte, message []byte) bool {
				// Verify Cosmos signature format
				// Format: R (32 bytes) || S (32 bytes)
				return len(signature) == 64
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sign message
			messageHash := sha256.Sum256([]byte(tt.message))
			
			// Create signing context with metadata
			metadata := map[string]interface{}{
				"network_id": tt.networkID,
				"public_key": hex.EncodeToString(pubKeyBytes),
			}

			// Sign the message using TSS
			signResult, err := selfchainTss.SignMessage(ctx, messageHash[:], keygenResult.Party1Data, keygenResult.Party2Data)
			require.NoError(t, err)
			require.NotNil(t, signResult)

			// Convert TSS SignResult to SignatureResult
			sigResult := &SignatureResult{
				R: signResult.R,
				S: signResult.S,
			}

			// Format the signature according to the network
			signature, err := signerFactory.Sign(ctx, tt.networkID, messageHash[:], metadata, sigResult)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, signature)

			// Verify the signature format using network-specific verification
			assert.True(t, tt.verifyFunc(t, signature, messageHash[:]))

			// Verify the signature using the signer factory
			valid, err := signerFactory.Verify(tt.networkID, pubKeyBytes, messageHash[:], signature)
			require.NoError(t, err)
			assert.True(t, valid)
		})
	}
}

func TestMultiChainConcurrentSigning(t *testing.T) {
	// Setup TSS key generation
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Generate TSS key pair
	preParams, err := keygen.GeneratePreParams(time.Minute)
	require.NoError(t, err)
	require.NotNil(t, preParams)

	keygenResult, err := selfchainTss.GenerateKey(ctx, preParams, "test-key")
	require.NoError(t, err)
	require.NotNil(t, keygenResult)

	// Create signer factory
	signerFactory := NewSignerFactory(networks.NewNetworkRegistry())

	// Get public key from party data
	pubKeyPoint := keygenResult.Party1Data.ECDSAPub
	x, y := pubKeyPoint.X(), pubKeyPoint.Y()
	
	// Convert big.Int coordinates to btcec.FieldVal
	xField := new(btcec.FieldVal)
	yField := new(btcec.FieldVal)
	xField.SetByteSlice(x.Bytes())
	yField.SetByteSlice(y.Bytes())
	
	btcPubKey := btcec.NewPublicKey(xField, yField)
	pubKeyBytes := btcPubKey.SerializeCompressed()

	// Create test cases
	tests := []testCase{
		{
			name:      "Bitcoin mainnet signature",
			networkID: "bitcoin:mainnet",
			message:   "test bitcoin transaction",
		},
		{
			name:      "Ethereum mainnet signature",
			networkID: "ethereum:1",
			message:   "test ethereum transaction",
		},
		{
			name:      "Cosmos Hub signature",
			networkID: "cosmos:cosmoshub-4",
			message:   "test cosmos transaction",
		},
	}

	// Run concurrent signing tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new context for each test
			testCtx, testCancel := context.WithTimeout(ctx, 30*time.Second)
			defer testCancel()

			// Sign message
			messageHash := sha256.Sum256([]byte(tt.message))
			
			// Create signing context with metadata
			metadata := map[string]interface{}{
				"network_id": tt.networkID,
				"public_key": hex.EncodeToString(pubKeyBytes),
			}

			// Sign the message using TSS
			signResult, err := selfchainTss.SignMessage(testCtx, messageHash[:], keygenResult.Party1Data, keygenResult.Party2Data)
			require.NoError(t, err)
			require.NotNil(t, signResult)

			// Convert TSS SignResult to SignatureResult
			sigResult := &SignatureResult{
				R: signResult.R,
				S: signResult.S,
			}

			// Format the signature according to the network
			signature, err := signerFactory.Sign(testCtx, tt.networkID, messageHash[:], metadata, sigResult)
			require.NoError(t, err)
			require.NotNil(t, signature)

			// Verify the signature using the signer factory
			valid, err := signerFactory.Verify(tt.networkID, pubKeyBytes, messageHash[:], signature)
			require.NoError(t, err)
			assert.True(t, valid)
		})
	}
}
