package types

import (
	sdkmath "cosmossdk.io/math"
)

// This is 1 Token based on 18 decimal points.
// TODO: we might want to move that to the store so we can change that value???
func GetMinMigrationAmount() sdkmath.Uint {
	return sdkmath.NewUint(1000000000000000000)
}

const DENOM = "ufront"

type Token int64

// Token enum
const (
	Front    Token = 0
	Hotcross Token = 1
)

// Ratios
const (
	FRONT_RATIO    = 10 // 10%
	HOTCROSS_RATIO = 5
)
