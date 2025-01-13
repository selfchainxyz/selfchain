package tss

import (
	"context"
	"testing"
	"time"
	"math/big"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		wantErr bool
	}{
		{
			name:    "Success case",
			timeout: 60 * time.Second,
			wantErr: false,
		},
		{
			name:    "Timeout case",
			timeout: 1 * time.Millisecond,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			preParams, err := keygen.GeneratePreParams(time.Minute)
			if err != nil {
				t.Fatalf("Failed to generate pre-parameters: %v", err)
				return
			}
			assert.NotNil(t, preParams)

			result, err := GenerateKey(ctx, preParams)
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

			// Verify party 1 data
			assert.NotNil(t, result.Party1Data)
			assert.NotNil(t, result.Party1Data.ECDSAPub)
			assert.NotNil(t, result.Party1Data.ShareID)
			assert.Equal(t, 0, result.Party1Data.ShareID.Cmp(big.NewInt(1)))

			// Verify party 2 data
			assert.NotNil(t, result.Party2Data)
			assert.NotNil(t, result.Party2Data.ECDSAPub)
			assert.NotNil(t, result.Party2Data.ShareID)
			assert.Equal(t, 0, result.Party2Data.ShareID.Cmp(big.NewInt(2)))
		})
	}
}

func TestGenerateKey_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	preParams, err := keygen.GeneratePreParams(time.Minute)
	if err != nil {
		t.Fatalf("Failed to generate pre-parameters: %v", err)
		return
	}
	assert.NotNil(t, preParams)

	// Cancel the context immediately after starting key generation
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	result, err := GenerateKey(ctx, preParams)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, context.Canceled, ctx.Err())
}
