package networks

import (
	"fmt"
	"sync"
)

// CurveType represents the elliptic curve type used by the network
type CurveType string

const (
	// Supported curve types
	Secp256k1  CurveType = "secp256k1"
	Ed25519    CurveType = "ed25519"
	Curve25519 CurveType = "curve25519"
	P256       CurveType = "p256"
	BLS12_381  CurveType = "bls12_381"
	Stark256   CurveType = "stark256"
)

// NetworkType represents the blockchain network type
type NetworkType string

const (
	// Major network types
	Bitcoin  NetworkType = "bitcoin"
	Ethereum NetworkType = "ethereum"
	Cosmos   NetworkType = "cosmos"
	Solana   NetworkType = "solana"
	Polkadot NetworkType = "polkadot"
	Cardano  NetworkType = "cardano"
	Algorand NetworkType = "algorand"
	Tron     NetworkType = "tron"
	Near     NetworkType = "near"
	Stellar  NetworkType = "stellar"
	Aptos    NetworkType = "aptos"
	Sui      NetworkType = "sui"
)

// SigningConfig contains network-specific signing parameters
type SigningConfig struct {
	// Bitcoin-specific
	P2PKHPrefix  uint8  // Pay to Public Key Hash prefix
	P2SHPrefix   uint8  // Pay to Script Hash prefix
	HRP          string // Human Readable Part for bech32 addresses
	Base58Hasher string // Base58 hashing algorithm

	// Ethereum-specific
	ChainID  string // EVM chain ID
	GasToken string // Native gas token symbol

	// Cosmos-specific
	AddressPrefix string // Bech32 address prefix
	PubKeyPrefix  string // Public key prefix

	// Additional parameters can be added for other chains
}

// NetworkInfo contains detailed information about a blockchain network
type NetworkInfo struct {
	NetworkType      NetworkType
	ChainID          string
	SigningAlgorithm SigningAlgorithm
	Curve            CurveType
	AddressPrefix    string
	CoinType         uint32
	Decimals         uint8
	SymbolName       string
	DisplayName      string
	SigningConfig    *SigningConfig // Network-specific signing parameters
}

// NetworkRegistry manages the registry of supported networks
type NetworkRegistry struct {
	networks map[string]*NetworkInfo
	mu       sync.RWMutex
}

// NewNetworkRegistry creates a new network registry
func NewNetworkRegistry() *NetworkRegistry {
	registry := &NetworkRegistry{
		networks: make(map[string]*NetworkInfo),
	}
	registry.initializeDefaultNetworks()
	return registry
}

// initializeDefaultNetworks sets up the default supported networks
func (r *NetworkRegistry) initializeDefaultNetworks() {
	defaultNetworks := []*NetworkInfo{
		{
			NetworkType:      Bitcoin,
			ChainID:          "mainnet",
			SigningAlgorithm: ECDSA,
			Curve:            Secp256k1,
			CoinType:         0,
			Decimals:         8,
			SymbolName:       "BTC",
			DisplayName:      "Bitcoin",
			SigningConfig: &SigningConfig{
				P2PKHPrefix:  0x00,
				P2SHPrefix:   0x05,
				HRP:          "bc",
				Base58Hasher: "sha256",
			},
		},
		{
			NetworkType:      Ethereum,
			ChainID:          "1",
			SigningAlgorithm: ECDSA,
			Curve:            Secp256k1,
			CoinType:         60,
			Decimals:         18,
			SymbolName:       "ETH",
			DisplayName:      "Ethereum",
			SigningConfig: &SigningConfig{
				ChainID:  "1",
				GasToken: "ETH",
			},
		},
		{
			NetworkType:      Cosmos,
			ChainID:          "cosmoshub-4",
			SigningAlgorithm: ECDSA,
			Curve:            Secp256k1,
			AddressPrefix:    "cosmos",
			CoinType:         118,
			Decimals:         6,
			SymbolName:       "ATOM",
			DisplayName:      "Cosmos Hub",
			SigningConfig: &SigningConfig{
				AddressPrefix: "cosmos",
				PubKeyPrefix:  "cosmospub",
			},
		},
		{
			NetworkType:      Solana,
			ChainID:          "mainnet-beta",
			SigningAlgorithm: EdDSA,
			Curve:            Ed25519,
			CoinType:         501,
			Decimals:         9,
			SymbolName:       "SOL",
			DisplayName:      "Solana",
		},
		// Add more networks as needed
	}

	for _, network := range defaultNetworks {
		r.RegisterNetwork(network)
	}
}

// RegisterNetwork adds a new network to the registry
func (r *NetworkRegistry) RegisterNetwork(info *NetworkInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%s-%s", info.NetworkType, info.ChainID)
	r.networks[key] = info
	return nil
}

// GetNetwork retrieves network information by type and chain ID
func (r *NetworkRegistry) GetNetwork(networkType NetworkType, chainID string) (*NetworkInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := fmt.Sprintf("%s-%s", networkType, chainID)
	if info, exists := r.networks[key]; exists {
		return info, nil
	}
	return nil, fmt.Errorf("network not found: %s-%s", networkType, chainID)
}

// ListNetworks returns all registered networks
func (r *NetworkRegistry) ListNetworks() []*NetworkInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	networks := make([]*NetworkInfo, 0, len(r.networks))
	for _, info := range r.networks {
		networks = append(networks, info)
	}
	return networks
}

// IsSupportedNetwork checks if a network is supported
func (r *NetworkRegistry) IsSupportedNetwork(networkType NetworkType, chainID string) bool {
	_, err := r.GetNetwork(networkType, chainID)
	return err == nil
}
