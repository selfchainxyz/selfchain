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

	// WalletKey defines the key for storing wallets
	WalletKey = "Wallet-value-"

	// ParamsKey defines the store key for module parameters
	ParamsKey = "params"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
