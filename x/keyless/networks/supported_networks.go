package networks

import (
	"strings"

	"selfchain/x/keyless/types"
)

// DefaultNetworks returns a list of default supported networks
func DefaultNetworks() []*types.NetworkParams {
	return []*types.NetworkParams{
		{
			NetworkType:      string(Bitcoin),
			ChainId:         "bitcoin-mainnet",
			SigningAlgorithm: string(ECDSA),
			CurveType:       string(Secp256k1),
			AddressPrefix:   "",
			CoinType:        0,
			Decimals:        8,
			SymbolName:      "BTC",
			DisplayName:     "Bitcoin",
			SigningConfig: &types.SigningConfig{
				P2PkhPrefix:    0x00,
				P2ShPrefix:     0x05,
				Base58Hasher:   "sha256d",
				ChainId:        "bitcoin-mainnet",
				GasToken:       "BTC",
				AddressPrefix:  "",
			},
		},
		{
			NetworkType:      string(Ethereum),
			ChainId:         "1",
			SigningAlgorithm: string(ECDSA),
			CurveType:       string(Secp256k1),
			AddressPrefix:   "0x",
			CoinType:        60,
			Decimals:        18,
			SymbolName:      "ETH",
			DisplayName:     "Ethereum",
			SigningConfig: &types.SigningConfig{
				ChainId:       "1",
				GasToken:      "ETH",
				AddressPrefix: "0x",
			},
		},
		{
			NetworkType:      string(Cosmos),
			ChainId:         "cosmoshub-4",
			SigningAlgorithm: string(ECDSA),
			CurveType:       string(Secp256k1),
			AddressPrefix:   "cosmos",
			CoinType:        118,
			Decimals:        6,
			SymbolName:      "ATOM",
			DisplayName:     "Cosmos Hub",
			SigningConfig: &types.SigningConfig{
				ChainId:       "cosmoshub-4",
				GasToken:      "ATOM",
				AddressPrefix: "cosmos",
			},
		},
		// Add more networks as needed
	}
}

// GetDefaultNetworkParams returns network parameters for the specified network ID
func GetDefaultNetworkParams(networkID string) *types.NetworkParams {
	// Split network ID into network type and chain ID
	parts := strings.Split(networkID, ":")
	if len(parts) != 2 {
		return nil
	}
	networkType, chainID := parts[0], parts[1]

	// Handle special cases
	switch {
	case networkType == "bitcoin" && chainID == "mainnet":
		chainID = "bitcoin-mainnet"
	}

	networks := DefaultNetworks()
	for _, network := range networks {
		// Match both network type and chain ID
		if strings.EqualFold(network.NetworkType, networkType) && strings.EqualFold(network.ChainId, chainID) {
			return network
		}
	}
	return nil
}
