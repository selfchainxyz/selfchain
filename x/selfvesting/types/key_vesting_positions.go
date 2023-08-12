package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// VestingPositionsKeyPrefix is the prefix to retrieve all VestingPositions
	VestingPositionsKeyPrefix = "VestingPositions/value/"
)

// VestingPositionsKey returns the store key to retrieve a VestingPositions from the index fields
func VestingPositionsKey(
	beneficiary string,
) []byte {
	var key []byte

	beneficiaryBytes := []byte(beneficiary)
	key = append(key, beneficiaryBytes...)
	key = append(key, []byte("/")...)

	return key
}
