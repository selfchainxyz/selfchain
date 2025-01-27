package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// Query types for permissions
type (
	QueryListPermissionsRequest struct {
		Pagination *query.PageRequest `json:"pagination,omitempty"`
	}

	QueryListPermissionsResponse struct {
		Permissions []*Permission       `json:"permissions"`
		Pagination  *query.PageResponse `json:"pagination,omitempty"`
	}
)
