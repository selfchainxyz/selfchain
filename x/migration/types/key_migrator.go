package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// MigratorKeyPrefix is the prefix to retrieve all Migrator
	MigratorKeyPrefix = "Migrator/value/"
)

// MigratorKey returns the store key to retrieve a Migrator from the index fields
func MigratorKey(
	migrator string,
) []byte {
	var key []byte

	migratorBytes := []byte(migrator)
	key = append(key, migratorBytes...)
	key = append(key, []byte("/")...)

	return key
}
