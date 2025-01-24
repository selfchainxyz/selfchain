package signing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"selfchain/x/keyless/types"
)

func TestSignerFactory(t *testing.T) {
	// Create test network params
	networkParams := &types.NetworkParams{
		ChainId:          "bitcoin:1",
		SigningAlgorithm: string(ECDSA),
	}

	t.Run("Test_Sign", func(t *testing.T) {
		// Create signer with network params
		signer := &UniversalSigner{
			networkParams: networkParams,
			ecdsaSigner:  NewECDSASigner(nil, nil),
			eddsaSigner:  NewEdDSASigner(nil, nil),
		}
		require.NotNil(t, signer)

		// Create signing context
		ctx := &SigningContext{
			NetworkParams: networkParams,
			Message:      []byte("test message"),
			Party1Data:   nil,
			Party2Data:   nil,
		}

		// Test signing
		sig, err := signer.Sign(context.Background(), ctx)
		require.NoError(t, err)
		require.NotNil(t, sig)
		require.NotNil(t, sig.Bytes)

		// Test unsupported algorithm
		networkParams.SigningAlgorithm = "unsupported"
		sig, err = signer.Sign(context.Background(), ctx)
		require.Error(t, err)
		require.Nil(t, sig)
	})

	t.Run("Test_Verify", func(t *testing.T) {
		// Reset network params
		networkParams.SigningAlgorithm = string(ECDSA)

		// Create signer with network params
		signer := &UniversalSigner{
			networkParams: networkParams,
			ecdsaSigner:  NewECDSASigner(nil, nil),
			eddsaSigner:  NewEdDSASigner(nil, nil),
		}
		require.NotNil(t, signer)

		// Create signing context
		ctx := &SigningContext{
			NetworkParams: networkParams,
			Message:      []byte("test message"),
			Party1Data:   nil,
			Party2Data:   nil,
		}

		// Sign message
		sig, err := signer.Sign(context.Background(), ctx)
		require.NoError(t, err)
		require.NotNil(t, sig)

		// Get public key from signer
		pubKey := signer.ecdsaSigner.pubKey.SerializeCompressed()
		require.NotNil(t, pubKey)

		// Verify signature
		valid, err := signer.Verify(context.Background(), ctx.Message, sig, pubKey)
		require.NoError(t, err)
		require.True(t, valid)

		// Test invalid signature
		invalidSig := &SignatureResult{
			Bytes: []byte("invalid signature"),
		}
		valid, err = signer.Verify(context.Background(), ctx.Message, invalidSig, pubKey)
		require.NoError(t, err)
		require.False(t, valid)

		// Test nil signature
		valid, err = signer.Verify(context.Background(), ctx.Message, nil, pubKey)
		require.Error(t, err)
		require.False(t, valid)

		// Test nil public key
		valid, err = signer.Verify(context.Background(), ctx.Message, sig, nil)
		require.Error(t, err)
		require.False(t, valid)
	})
}
