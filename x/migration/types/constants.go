package types

import (
	sdkmath "cosmossdk.io/math"
)

// 1 Self will be credited to user account so future transaction to release from vesting position are possible
func GetInstantlyReleasedAmount() sdkmath.Uint {
	return sdkmath.NewUint(1000000)
}

const DENOM = "uself"

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

// // Vesting info
// const (
// 	SECONDS_IN_DAY   = 60 * 60 * 24
// 	VESTING_DURATION = SECONDS_IN_DAY * 30
// 	VESTING_CLIFF    = SECONDS_IN_DAY * 7
// )
