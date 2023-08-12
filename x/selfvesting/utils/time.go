package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BlockTime(ctx sdk.Context) uint64 {
	return uint64(ctx.BlockHeader().Time.Unix())
}
