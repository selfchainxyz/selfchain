package networks

import (
	"fmt"
)

// SigningAlgorithm represents the type of signing algorithm
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

// NetworkConfig represents the configuration for a blockchain network
type NetworkConfig struct {
	// ChainID is the unique identifier for the network
	ChainID string
	// Name is the human-readable name of the network
	Name string
	// Algorithm is the signing algorithm used by this network
	Algorithm SigningAlgorithm
	// Prefix is the address prefix used by this network
	Prefix string
	// CoinType is the BIP44 coin type
	CoinType uint32
	// Derivation is the derivation path template
	Derivation string
}

// NetworkRegistry maintains a registry of supported networks
type NetworkRegistry struct {
	networks map[string]NetworkConfig
}

// NewNetworkRegistry creates a new network registry
func NewNetworkRegistry() *NetworkRegistry {
	return &NetworkRegistry{
		networks: make(map[string]NetworkConfig),
	}
}

// RegisterNetwork adds a new network configuration to the registry
func (r *NetworkRegistry) RegisterNetwork(config NetworkConfig) error {
	if _, exists := r.networks[config.ChainID]; exists {
		return fmt.Errorf("network with chain ID %s already registered", config.ChainID)
	}
	r.networks[config.ChainID] = config
	return nil
}

// GetNetwork retrieves a network configuration by chain ID
func (r *NetworkRegistry) GetNetwork(chainID string) (NetworkConfig, error) {
	config, exists := r.networks[chainID]
	if !exists {
		return NetworkConfig{}, fmt.Errorf("network with chain ID %s not found", chainID)
	}
	return config, nil
}

// DefaultRegistry creates a registry with default network configurations
func DefaultRegistry() *NetworkRegistry {
	registry := NewNetworkRegistry()

	// Register Cosmos SDK based chains
	registry.RegisterNetwork(NetworkConfig{
		ChainID:    "selfchain_1234-1",
		Name:       "Self Chain",
		Algorithm:  ECDSA,
		Prefix:     "self",
		CoinType:   118,
		Derivation: "m/44'/118'/0'/0/0",
	})

	// Register Ethereum
	registry.RegisterNetwork(NetworkConfig{
		ChainID:    "1",
		Name:       "Ethereum Mainnet",
		Algorithm:  ECDSA,
		Prefix:     "0x",
		CoinType:   60,
		Derivation: "m/44'/60'/0'/0/0",
	})

	return registry
}
