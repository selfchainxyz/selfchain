package types

const (
	// ModuleName defines the module name
	ModuleName = "migration"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_migration"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	AclKey = "Acl/value/"
)

const (
	ConfigKey = "Config/value/"
)
