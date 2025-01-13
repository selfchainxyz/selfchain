package networks

import (
	"selfchain/x/keyless/types"
)

// Network IDs as defined in TrustWallet's wallet-core
const (
	// Bitcoin and Bitcoin-like
	BitcoinID     = "bitcoin"
	LitecoinID    = "litecoin"
	DogecoinID    = "doge"
	BitcoinCashID = "bitcoincash"
	BitcoinGoldID = "bitcoingold"
	GroestlcoinID = "groestlcoin"

	// Ethereum and EVM chains
	EthereumID        = "ethereum"
	EthereumClassicID = "classic"
	BaseID            = "base"
	OptimismID        = "optimism"
	ArbitrumID        = "arbitrum"
	ArbitrumNovaID    = "arbitrumnova"
	BSCID             = "bsc"
	PolygonID         = "polygon"
	AvalancheID       = "avalanchec"
	EvmosID           = "nativeevmos"
	ConfluxID         = "cfxevm"
	ZetaID            = "zetaevm"
	MerlinID          = "merlin"

	// Cosmos ecosystem
	CosmosID    = "cosmos"
	OsmosisID   = "osmosis"
	AxelarID    = "axelar"
	CoreumID    = "coreum"
	CryptoOrgID = "cryptoorg"

	// Other major chains
	SolanaID   = "solana"
	TronID     = "tron"
	CardanoID  = "cardano"
	AlgorandID = "algorand"
	AptosID    = "aptos"
	SuiID      = "sui"
	FilecoinID = "filecoin"
	HederaID   = "hedera"
	NearID     = "near"
	StellarID  = "stellar"
	ZilliqaID  = "zilliqa"
)

// Initialize supported networks with their configurations
func (r *NetworkRegistry) initializeSupportedNetworks() {
	networks := DefaultNetworks()

	// Register all networks
	for _, network := range networks {
		r.RegisterNetwork(network)
	}
}

// DefaultNetworks returns a list of default supported networks
func DefaultNetworks() []*types.NetworkInfo {
	return []*types.NetworkInfo{
		{
			NetworkType: types.Bitcoin,
			ChainID:     "mainnet",
			SigningConfig: &types.SigningConfig{
				SigningAlgorithm: types.SigningAlgoSecp256k1,
				CoinType:         0,
				Decimals:         8,
			},
		},
		{
			NetworkType: types.Ethereum,
			ChainID:     "1",
			SigningConfig: &types.SigningConfig{
				SigningAlgorithm: types.SigningAlgoSecp256k1,
				CoinType:         60,
				Decimals:         18,
			},
		},
		{
			NetworkType: types.Cosmos,
			ChainID:     "cosmoshub-4",
			SigningConfig: &types.SigningConfig{
				SigningAlgorithm: types.SigningAlgoSecp256k1,
				AddressPrefix:    "cosmos",
				CoinType:         118,
				Decimals:         6,
			},
		},
		{
			NetworkType: types.Solana,
			ChainID:     "mainnet-beta",
			SigningConfig: &types.SigningConfig{
				SigningAlgorithm: types.SigningAlgoEd25519,
				CoinType:         501,
				Decimals:         9,
			},
		},
		{
			NetworkType: types.Cardano,
			ChainID:     "mainnet",
			SigningConfig: &types.SigningConfig{
				SigningAlgorithm: types.SigningAlgoEd25519,
				CoinType:         1815,
				Decimals:         6,
			},
		},
		{
			NetworkType: types.Aptos,
			ChainID:     "mainnet",
			SigningConfig: &types.SigningConfig{
				SigningAlgorithm: types.SigningAlgoEd25519,
				CoinType:         637,
				Decimals:         8,
			},
		},
		{
			NetworkType: types.Sui,
			ChainID:     "mainnet",
			SigningConfig: &types.SigningConfig{
				SigningAlgorithm: types.SigningAlgoEd25519,
				CoinType:         784,
				Decimals:         9,
			},
		},
	}
}
