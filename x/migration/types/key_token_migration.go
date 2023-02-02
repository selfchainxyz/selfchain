package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// TokenMigrationKeyPrefix is the prefix to retrieve all TokenMigration
	TokenMigrationKeyPrefix = "TokenMigration/value/"
)

// TokenMigrationKey returns the store key to retrieve a TokenMigration from the index fields
func TokenMigrationKey(
	msgHash string,
) []byte {
	var key []byte

	msgHashBytes := []byte(msgHash)
	key = append(key, msgHashBytes...)
	key = append(key, []byte("/")...)

	return key
}
