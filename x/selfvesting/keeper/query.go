package keeper

import (
	"selfchain/x/selfvesting/types"
)

var _ types.QueryServer = Keeper{}
