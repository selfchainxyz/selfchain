package keeper

import (
	"selfchain/x/migration/types"
)

var _ types.QueryServer = Keeper{}
