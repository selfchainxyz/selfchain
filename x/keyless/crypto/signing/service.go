package signing

import (
	"context"
	"selfchain/x/keyless/crypto/signing/types"
)

// SigningAlgorithm represents supported signing algorithms
type SigningAlgorithm string

const (
	ECDSA SigningAlgorithm = "ECDSA"
	EdDSA SigningAlgorithm = "EdDSA"
)

// SignatureResult contains the final signature data
type SignatureResult struct {
	R        *big.Int
	S        *big.Int
	V        uint8  // For ECDSA recovery
	Bytes    []byte // Raw signature bytes
	Recovery byte   // Recovery ID for ECDSA
}

// SigningService defines the interface for cryptographic signing operations
type SigningService interface {
	// Sign signs a message using the specified algorithm
	Sign(ctx context.Context, msg []byte, algorithm types.SigningAlgorithm) (*types.SignatureResult, error)

	// Verify verifies a signature
	Verify(ctx context.Context, msg []byte, sig *types.SignatureResult, pubKeyBytes []byte) (bool, error)

	// GetPublicKey returns the public key for the signing service
	GetPublicKey(ctx context.Context, algorithm types.SigningAlgorithm) ([]byte, error)
}
