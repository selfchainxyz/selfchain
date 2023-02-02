package keeper

import (
	"frontier/x/migration/types"
)

var _ types.QueryServer = Keeper{}
