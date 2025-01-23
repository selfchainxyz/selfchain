package signing

import (
	"context"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
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
	Sign(ctx context.Context, msg []byte, algorithm SigningAlgorithm) (*SignatureResult, error)

	// Verify verifies a signature
	Verify(ctx context.Context, msg []byte, sig *SignatureResult, pubKey *btcec.PublicKey) (bool, error)

	// GetPublicKey returns the public key for the signing service
	GetPublicKey(ctx context.Context, algorithm SigningAlgorithm) (*btcec.PublicKey, error)
}
