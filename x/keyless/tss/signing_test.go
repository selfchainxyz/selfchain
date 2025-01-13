package tss

import (
	"context"
	"crypto/sha256"
	"testing"
	"time"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/stretchr/testify/assert"
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

	keygenResult, err := GenerateKey(ctx, preParams)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
		return
	}
	if !assert.NotNil(t, keygenResult) {
		return
	}

	// Test message signing
	tests := []struct {
		name            string
		msg            []byte
		personalSaveData *keygen.LocalPartySaveData
		remoteSaveData   *keygen.LocalPartySaveData
		timeout         time.Duration
		wantErr         bool
	}{
		{
			name:            "successful signing",
			msg:            []byte("test message"),
			personalSaveData: keygenResult.Party1Data,
			remoteSaveData:   keygenResult.Party2Data,
			timeout:         60 * time.Second,
			wantErr:         false,
		},
		{
			name:            "timeout during signing",
			msg:            []byte("test message"),
			personalSaveData: keygenResult.Party1Data,
			remoteSaveData:   keygenResult.Party2Data,
			timeout:         1 * time.Millisecond,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			// Hash the message before signing
			msgHash := sha256.Sum256(tt.msg)

			result, err := SignMessage(ctx, msgHash[:], tt.personalSaveData, tt.remoteSaveData)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			if !assert.NoError(t, err) {
				return
			}
			if !assert.NotNil(t, result) {
				return
			}

			// Verify signature components
			assert.NotNil(t, result.R)
			assert.NotNil(t, result.S)
			assert.True(t, result.R.BitLen() > 0)
			assert.True(t, result.S.BitLen() > 0)

			// TODO: Add ECDSA signature verification when we have proper key reconstruction
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

	keygenResult, err := GenerateKey(ctx, preParams)
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

	keygenResult, err := GenerateKey(ctx, preParams)
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

	// TODO: Add ECDSA signature verification when we have proper key reconstruction
}
