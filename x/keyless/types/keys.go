package types

const (
	// ModuleName defines the module name
	ModuleName = "keyless"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_keyless"

	// Key prefixes
	WalletKey    = "Wallet-"
	KeyShareKey  = "keyshare"
	ParamsKey    = "params"
	PartyDataKey = "PartyData-"
)

// KeyPrefix returns the KVStore key prefix for the given key type
func KeyPrefix(p string) []byte {
	return []byte(p)
}
