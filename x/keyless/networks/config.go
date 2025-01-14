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
	registry := NewNetworkRegistry()
	
	// Try to find the network in registry
	var networkType NetworkType
	switch {
	case chainID == "1" || chainID == "5":
		networkType = Ethereum
	case chainID == "bitcoin" || chainID == "bitcoin-testnet":
		networkType = Bitcoin
	default:
		networkType = Cosmos // Default to Cosmos for other chains
	}
	
	networkInfo, err := registry.GetNetwork(networkType, chainID)
	if err != nil {
		return nil, err
	}
	
	return &NetworkConfig{
		ChainID:          networkInfo.ChainID,
		SigningAlgorithm: networkInfo.SigningAlgorithm,
		AddressPrefix:    networkInfo.AddressPrefix,
		CoinType:         networkInfo.CoinType,
	}, nil
}
