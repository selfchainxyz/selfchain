package types

import uint256 "github.com/holiman/uint256"

// This is 1 Token based on 18 decimal points.
// TODO: we might want to move that to the store so we can change that value???
func getMinMigrationAmount() *uint256.Int {
	return uint256.NewInt(1000000000000000000)
}

type Token int64

// Token enum
const (
	Front    Token = 0
	Hotcross Token = 1
)
