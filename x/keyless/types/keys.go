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
	KeyRotationKey = "key_rotation"
	AuditEventKey = "audit_event"
	BatchSignStatusPrefix = "batch_sign_status"
	KeyRotationStatusPrefix = "key_rotation_status"
)

// KeyPrefix returns the KVStore key prefix for the given key type
func KeyPrefix(p string) []byte {
	return []byte(p)
}

// BatchSignStatusKey returns the store key for batch sign status
func BatchSignStatusKey(walletId string, batchId string) []byte {
	key := []byte(BatchSignStatusPrefix)
	key = append(key, []byte(walletId)...)
	key = append(key, []byte("/")...)
	key = append(key, []byte(batchId)...)
	return key
}

// KeyRotationStatusKey returns the store key for key rotation status
func KeyRotationStatusKey(walletId string) []byte {
	key := []byte(KeyRotationStatusPrefix)
	key = append(key, []byte(walletId)...)
	return key
}
