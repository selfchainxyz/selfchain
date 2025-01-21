package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"
)

// Paginate is a helper function for pagination
func (k Keeper) Paginate(
	store prefix.Store,
	pageRequest *query.PageRequest,
	onResult func(key []byte, value []byte) error,
) (*query.PageResponse, error) {
	// Create a paginated iterator
	pageRes, err := query.Paginate(store, pageRequest, onResult)
	if err != nil {
		return nil, err
	}

	return pageRes, nil
}
