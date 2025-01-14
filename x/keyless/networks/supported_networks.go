package networks

import (
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
		// Add more networks as needed
	}
}

// GetDefaultNetworkParams returns network parameters for the specified network ID
func GetDefaultNetworkParams(networkID string) *types.NetworkParams {
	networks := DefaultNetworks()
	for _, network := range networks {
		if network.ChainId == networkID {
			return network
		}
	}
	return nil
}
