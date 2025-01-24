package types

import (
	"math/big"
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
