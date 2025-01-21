package common

// NetworkType represents different types of networks
type NetworkType string

const (
	Bitcoin  NetworkType = "bitcoin"
	Ethereum NetworkType = "ethereum"
	Cosmos   NetworkType = "cosmos"
	Solana   NetworkType = "solana"
	Cardano  NetworkType = "cardano"
	Aptos    NetworkType = "aptos"
	Sui      NetworkType = "sui"
)

// SigningAlgorithm represents supported signing algorithms
type SigningAlgorithm string

const (
	// ECDSA variants
	SigningAlgoSecp256k1 SigningAlgorithm = "secp256k1"
	SigningAlgoP256      SigningAlgorithm = "p256"
	
	// EdDSA variants
	SigningAlgoEd25519   SigningAlgorithm = "ed25519"
	
	// Future algorithms
	SigningAlgoBLS       SigningAlgorithm = "bls"
	SigningAlgoSchnorr   SigningAlgorithm = "schnorr"
)

// SignatureFormat represents different signature formats
type SignatureFormat string

const (
	// DER ASN.1 format (Bitcoin, traditional ECDSA)
	SigFormatDER SignatureFormat = "der"
	
	// RSV format (Ethereum, EVM chains)
	SigFormatRSV SignatureFormat = "rsv"
	
	// RS format (Cosmos chains)
	SigFormatRS SignatureFormat = "rs"
	
	// Pure Ed25519 format (Solana, Cardano, etc.)
	SigFormatEd25519 SignatureFormat = "ed25519"
)

// HashAlgorithm represents supported hashing algorithms
type HashAlgorithm string

const (
	HashAlgoSha256    HashAlgorithm = "sha256"
	HashAlgoKeccak256 HashAlgorithm = "keccak256"
	HashAlgoBlake2b   HashAlgorithm = "blake2b"
)
