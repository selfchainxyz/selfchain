package types

import (
	"context"
	"math/big"
)

// SigningAlgorithm represents supported signing algorithms
type SigningAlgorithm string

const (
	ECDSA SigningAlgorithm = "ECDSA"
	EdDSA SigningAlgorithm = "EdDSA"
)

// SignatureResult contains the signature components
type SignatureResult struct {
	R     *big.Int
	S     *big.Int
	V     uint8
	Bytes []byte
}

// NetworkParams contains network-specific parameters
type NetworkParams struct {
	NetworkType      string
	ChainID         string
	SigningAlgorithm string
	CurveType       string
	AddressPrefix   string
	CoinType        uint32
	Decimals        uint8
	SymbolName      string
	DisplayName     string
	SigningConfig   *SigningConfig
}

// SigningConfig contains network-specific signing configuration
type SigningConfig struct {
	ChainID       string
	GasToken      string
	AddressPrefix string
}

// SigningService defines the interface for signing operations
type SigningService interface {
	Sign(ctx context.Context, message []byte, algorithm SigningAlgorithm) (*SignatureResult, error)
	Verify(ctx context.Context, message []byte, signature *SignatureResult, pubKeyBytes []byte) (bool, error)
	GetPublicKey(ctx context.Context, algorithm SigningAlgorithm) ([]byte, error)
}
