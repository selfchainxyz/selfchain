package networks

import (
	"fmt"
)

// SigningAlgorithm represents a supported signing algorithm
type SigningAlgorithm string

const (
	// ECDSA signing algorithm
	ECDSA SigningAlgorithm = "ECDSA"
	// EdDSA signing algorithm
	EdDSA SigningAlgorithm = "EdDSA"
	// Schnorr signing algorithm
	Schnorr SigningAlgorithm = "Schnorr"
	// BLS signing algorithm
	BLS SigningAlgorithm = "BLS"
)

// NetworkConfig represents the configuration for a specific network
type NetworkConfig struct {
	ChainID          string
	SigningAlgorithm SigningAlgorithm
	AddressPrefix    string
	CoinType         uint32
}

// ValidateSigningAlgorithm checks if the signing algorithm is supported
func ValidateSigningAlgorithm(algo SigningAlgorithm) error {
	switch algo {
	case ECDSA, EdDSA:
		return nil
	case Schnorr, BLS:
		return fmt.Errorf("signing algorithm %s not yet implemented", algo)
	default:
		return fmt.Errorf("unsupported signing algorithm: %s", algo)
	}
}

// GetNetworkConfig returns the configuration for a specific chain ID
func GetNetworkConfig(chainID string) (*NetworkConfig, error) {
	// TODO: Implement network registry
	// For now, return default config for testing
	return &NetworkConfig{
		ChainID:          chainID,
		SigningAlgorithm: ECDSA,
		AddressPrefix:    "self",
		CoinType:         118,
	}, nil
}
