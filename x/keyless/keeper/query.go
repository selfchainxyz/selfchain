package keeper

import (
	"selfchain/x/keyless/types"
)

var _ types.QueryServer = Keeper{}
