package types

import (
	sdkmath "cosmossdk.io/math"
)

// 1 Self will be credited to user account so future transaction to release from vesting position are possible
func GetInstantlyReleasedAmount() sdkmath.Uint {
	return sdkmath.NewUint(1000000)
}

const DENOM = "uslf"

type Token int64

// Token enum
const (
	Front    Token = 0
	Hotcross Token = 1
)

// Ratios
const (
	FRONT_RATIO    = 100 // 100%
)
