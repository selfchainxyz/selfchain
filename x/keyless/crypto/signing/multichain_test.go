package signing

import (
	"context"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/btcsuite/btcd/btcec/v2"

	"selfchain/x/keyless/networks"
	"selfchain/x/keyless/tss"
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

	keygenResult, err := tss.GenerateKey(ctx, preParams, "test-key")
	require.NoError(t, err)
	require.NotNil(t, keygenResult)

	// Create signer factory
	signerFactory := NewSignerFactory(networks.NewNetworkRegistry())

	// Get public key from party data
	pubKeyPoint := keygenResult.Party1Data.ECDSAPub
	curve := btcec.S256()
	pubKeyBytes := elliptic.Marshal(curve, pubKeyPoint.X(), pubKeyPoint.Y())
	require.NotNil(t, pubKeyBytes)

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
			signResult, err := tss.SignMessage(ctx, messageHash[:], keygenResult.Party1Data, keygenResult.Party2Data)
			require.NoError(t, err)
			require.NotNil(t, signResult)

			// Format the signature according to the network
			signature, err := signerFactory.Sign(ctx, tt.networkID, messageHash[:], metadata, signResult)
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

	keygenResult, err := tss.GenerateKey(ctx, preParams, "test-key")
	require.NoError(t, err)
	require.NotNil(t, keygenResult)

	// Get public key from party data
	pubKeyPoint := keygenResult.Party1Data.ECDSAPub
	curve := btcec.S256()
	pubKeyBytes := elliptic.Marshal(curve, pubKeyPoint.X(), pubKeyPoint.Y())
	require.NotNil(t, pubKeyBytes)

	// Create signer factory
	signerFactory := NewSignerFactory(networks.NewNetworkRegistry())

	// Create test cases for concurrent signing
	networks := []string{"bitcoin:mainnet", "ethereum:1", "cosmos:cosmoshub-4"}
	numRequests := 5 // Number of concurrent requests per network

	// Create channels for results
	type result struct {
		networkID string
		signature []byte
		err       error
	}
	results := make(chan result, len(networks)*numRequests)

	// Start concurrent signing requests
	for _, networkID := range networks {
		for i := 0; i < numRequests; i++ {
			go func(networkID string, i int) {
				message := []byte(fmt.Sprintf("test message %d for %s", i, networkID))
				messageHash := sha256.Sum256(message)

				metadata := map[string]interface{}{
					"network_id": networkID,
					"public_key": hex.EncodeToString(pubKeyBytes),
				}

				// Sign the message using TSS
				signResult, err := tss.SignMessage(ctx, messageHash[:], keygenResult.Party1Data, keygenResult.Party2Data)
				if err != nil {
					results <- result{networkID: networkID, err: err}
					return
				}

				// Format the signature according to the network
				signature, err := signerFactory.Sign(ctx, networkID, messageHash[:], metadata, signResult)
				results <- result{networkID: networkID, signature: signature, err: err}
			}(networkID, i)
		}
	}

	// Collect and verify results
	successCount := make(map[string]int)
	timeout := time.After(2 * time.Minute)

	for i := 0; i < len(networks)*numRequests; i++ {
		select {
		case r := <-results:
			if r.err != nil {
				t.Errorf("Error signing for network %s: %v", r.networkID, r.err)
				continue
			}
			successCount[r.networkID]++

		case <-timeout:
			t.Fatal("Test timed out waiting for signing results")
		}
	}

	// Verify that we got the expected number of successful signatures for each network
	for _, networkID := range networks {
		assert.Equal(t, numRequests, successCount[networkID],
			"Expected %d successful signatures for network %s, got %d",
			numRequests, networkID, successCount[networkID])
	}
}
