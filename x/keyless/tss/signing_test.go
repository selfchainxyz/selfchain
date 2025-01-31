package tss

import (
	"context"
	"crypto/sha256"
	"testing"
	"time"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignMessage(t *testing.T) {
	// First generate a key pair
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Generate pre-parameters
	preParams, err := keygen.GeneratePreParams(time.Minute)
	if err != nil {
		t.Fatalf("Failed to generate pre-parameters: %v", err)
		return
	}
	assert.NotNil(t, preParams)

	keygenResult, err := GenerateKey(ctx, preParams, "test-chain-1")
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
		return
	}
	if !assert.NotNil(t, keygenResult) {
		return
	}

	// Test cases
	tests := []struct {
		name    string
		msg     []byte
		wantErr bool
	}{
		{
			name:    "Valid message",
			msg:     []byte("test message"),
			wantErr: false,
		},
		{
			name:    "Empty message",
			msg:     []byte{},
			wantErr: true,
		},
		{
			name:    "successful signing",
			msg:     []byte("test message"),
			wantErr: false,
		},
		{
			name:    "timeout during signing",
			msg:     []byte("test message"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			if tt.name == "timeout during signing" {
				ctx, cancel = context.WithTimeout(context.Background(), 1*time.Millisecond)
			}
			defer cancel()

			// Hash the message before signing
			msgHash := sha256.Sum256(tt.msg)

			result, err := SignMessage(ctx, msgHash[:], keygenResult.Party1Data, keygenResult.Party2Data)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if !assert.NotNil(t, result) {
				return
			}

			// Verify signature components
			assert.NotNil(t, result.R)
			assert.NotNil(t, result.S)
			assert.True(t, result.R.BitLen() > 0)
			assert.True(t, result.S.BitLen() > 0)

			// Verify ECDSA signature
			pubKey := keygenResult.PublicKey

			// Convert big.Int to btcec types
			rScalar := new(btcec.ModNScalar)
			rScalar.SetByteSlice(result.R.Bytes())
			sScalar := new(btcec.ModNScalar)
			sScalar.SetByteSlice(result.S.Bytes())

			// Create btcec signature
			signature := ecdsa.NewSignature(rScalar, sScalar)

			// Hash the message
			msgHash = sha256.Sum256(tt.msg)

			// Convert to btcec public key
			x := new(btcec.FieldVal)
			x.SetByteSlice(pubKey.X.Bytes())
			y := new(btcec.FieldVal)
			y.SetByteSlice(pubKey.Y.Bytes())
			btcecPubKey := btcec.NewPublicKey(x, y)

			// Verify the signature using btcec
			valid := signature.Verify(msgHash[:], btcecPubKey)
			assert.True(t, valid)
		})
	}
}

func TestSignMessageCancel(t *testing.T) {
	// First generate a key pair
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	preParams, err := keygen.GeneratePreParams(time.Minute)
	if err != nil {
		t.Fatalf("Failed to generate pre-parameters: %v", err)
		return
	}
	assert.NotNil(t, preParams)

	keygenResult, err := GenerateKey(ctx, preParams, "test-chain-1")
	cancel()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
		return
	}
	if !assert.NotNil(t, keygenResult) {
		return
	}

	// Test cancellation during signing
	ctx, cancel = context.WithCancel(context.Background())

	// Cancel the context immediately after starting signing
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	msgHash := sha256.Sum256([]byte("test message"))
	result, err := SignMessage(ctx, msgHash[:], keygenResult.Party1Data, keygenResult.Party2Data)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, context.Canceled, ctx.Err())
}

func TestEndToEnd(t *testing.T) {
	// Test the entire flow: key generation -> signing -> verification
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 1. Generate key pair
	preParams, err := keygen.GeneratePreParams(time.Minute)
	if err != nil {
		t.Fatalf("Failed to generate pre-parameters: %v", err)
		return
	}
	assert.NotNil(t, preParams)

	keygenResult, err := GenerateKey(ctx, preParams, "test-chain-1")
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
		return
	}
	if !assert.NotNil(t, keygenResult) {
		return
	}

	// 2. Sign a message
	msg := []byte("test message")
	msgHash := sha256.Sum256(msg)

	signResult, err := SignMessage(ctx, msgHash[:], keygenResult.Party1Data, keygenResult.Party2Data)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, signResult) {
		return
	}

	// 3. Verify signature components
	assert.NotNil(t, signResult.R)
	assert.NotNil(t, signResult.S)
	assert.True(t, signResult.R.BitLen() > 0)
	assert.True(t, signResult.S.BitLen() > 0)

	// Verify ECDSA signature
	pubKey := keygenResult.PublicKey

	// Convert big.Int to btcec types
	rScalar := new(btcec.ModNScalar)
	rScalar.SetByteSlice(signResult.R.Bytes())
	sScalar := new(btcec.ModNScalar)
	sScalar.SetByteSlice(signResult.S.Bytes())

	// Create btcec signature
	signature := ecdsa.NewSignature(rScalar, sScalar)

	// Hash the message
	msgHash = sha256.Sum256(msg)

	// Convert to btcec public key
	x := new(btcec.FieldVal)
	x.SetByteSlice(pubKey.X.Bytes())
	y := new(btcec.FieldVal)
	y.SetByteSlice(pubKey.Y.Bytes())
	btcecPubKey := btcec.NewPublicKey(x, y)

	// Verify the signature using btcec
	valid := signature.Verify(msgHash[:], btcecPubKey)
	require.True(t, valid, "signature verification failed")
}
