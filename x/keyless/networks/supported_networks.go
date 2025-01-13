package networks

// Network IDs as defined in TrustWallet's wallet-core
const (
	// Bitcoin and Bitcoin-like
	BitcoinID         = "bitcoin"
	LitecoinID        = "litecoin"
	DogecoinID        = "doge"
	BitcoinCashID     = "bitcoincash"
	BitcoinGoldID     = "bitcoingold"
	GroestlcoinID     = "groestlcoin"
	
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
	CosmosID          = "cosmos"
	OsmosisID         = "osmosis"
	AxelarID          = "axelar"
	CoreumID          = "coreum"
	CryptoOrgID       = "cryptoorg"
	
	// Other major chains
	SolanaID          = "solana"
	TronID            = "tron"
	CardanoID         = "cardano"
	AlgorandID        = "algorand"
	AptosID           = "aptos"
	SuiID             = "sui"
	FilecoinID        = "filecoin"
	HederaID          = "hedera"
	NearID            = "near"
	StellarID         = "stellar"
	ZilliqaID         = "zilliqa"
)

// Initialize supported networks with their configurations
func (r *NetworkRegistry) initializeSupportedNetworks() {
	networks := []*NetworkInfo{
		// Bitcoin and Bitcoin-like networks
		{
			NetworkType:      Bitcoin,
			ChainID:         "1",
			SigningAlgorithm: ECDSA,
			Curve:           Secp256k1,
			CoinType:        0,
			Decimals:        8,
			SymbolName:      "BTC",
			DisplayName:     "Bitcoin",
		},
		{
			NetworkType:      Bitcoin,
			ChainID:         "2",
			SigningAlgorithm: ECDSA,
			Curve:           Secp256k1,
			CoinType:        2,
			Decimals:        8,
			SymbolName:      "LTC",
			DisplayName:     "Litecoin",
		},
		
		// Ethereum and EVM chains
		{
			NetworkType:      Ethereum,
			ChainID:         "1",
			SigningAlgorithm: ECDSA,
			Curve:           Secp256k1,
			CoinType:        60,
			Decimals:        18,
			SymbolName:      "ETH",
			DisplayName:     "Ethereum",
		},
		{
			NetworkType:      Ethereum,
			ChainID:         "8453",
			SigningAlgorithm: ECDSA,
			Curve:           Secp256k1,
			CoinType:        8453,
			Decimals:        18,
			SymbolName:      "ETH",
			DisplayName:     "Base",
		},
		{
			NetworkType:      Ethereum,
			ChainID:         "10",
			SigningAlgorithm: ECDSA,
			Curve:           Secp256k1,
			CoinType:        10000070,
			Decimals:        18,
			SymbolName:      "ETH",
			DisplayName:     "Optimism",
		},
		
		// Cosmos ecosystem
		{
			NetworkType:      Cosmos,
			ChainID:         "cosmoshub-4",
			SigningAlgorithm: ECDSA,
			Curve:           Secp256k1,
			AddressPrefix:   "cosmos",
			CoinType:        118,
			Decimals:        6,
			SymbolName:      "ATOM",
			DisplayName:     "Cosmos Hub",
		},
		{
			NetworkType:      Cosmos,
			ChainID:         "osmosis-1",
			SigningAlgorithm: ECDSA,
			Curve:           Secp256k1,
			AddressPrefix:   "osmo",
			CoinType:        118,
			Decimals:        6,
			SymbolName:      "OSMO",
			DisplayName:     "Osmosis",
		},
		
		// Other major chains
		{
			NetworkType:      Solana,
			ChainID:         "mainnet-beta",
			SigningAlgorithm: EdDSA,
			Curve:           Ed25519,
			CoinType:        501,
			Decimals:        9,
			SymbolName:      "SOL",
			DisplayName:     "Solana",
		},
		{
			NetworkType:      Cardano,
			ChainID:         "mainnet",
			SigningAlgorithm: EdDSA,
			Curve:           Ed25519,
			CoinType:        1815,
			Decimals:        6,
			SymbolName:      "ADA",
			DisplayName:     "Cardano",
		},
		{
			NetworkType:      Aptos,
			ChainID:         "1",
			SigningAlgorithm: EdDSA,
			Curve:           Ed25519,
			CoinType:        637,
			Decimals:        8,
			SymbolName:      "APT",
			DisplayName:     "Aptos",
		},
		{
			NetworkType:      Sui,
			ChainID:         "mainnet",
			SigningAlgorithm: EdDSA,
			Curve:           Ed25519,
			CoinType:        784,
			Decimals:        9,
			SymbolName:      "SUI",
			DisplayName:     "Sui",
		},
	}

	// Register all networks
	for _, network := range networks {
		r.RegisterNetwork(network)
	}
}
