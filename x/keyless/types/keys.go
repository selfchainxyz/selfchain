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
)

var (
	// WalletKeyPrefix is the prefix for storing wallets
	WalletKeyPrefix = []byte{0x01}
)

// WalletKey returns the store key to retrieve a Wallet from the index fields
func WalletKey(address string) []byte {
	return append(WalletKeyPrefix, []byte(address)...)
}

func KeyPrefix(p string) []byte {
	return []byte(p)
}
