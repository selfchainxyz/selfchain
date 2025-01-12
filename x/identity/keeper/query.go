package keeper

import (
	"context"

	"selfchain/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CredentialsByDID returns all credentials for a given DID
func (k Keeper) CredentialsByDID(goCtx context.Context, req *types.QueryCredentialsByDIDRequest) (*types.QueryCredentialsByDIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get all credentials for the DID
	store := ctx.KVStore(k.storeKey)
	prefix := append([]byte(types.CredentialByDIDPrefix), []byte(req.Did)...)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var credentials []*types.Credential
	for ; iterator.Valid(); iterator.Next() {
		var credential types.Credential
		k.cdc.MustUnmarshal(iterator.Value(), &credential)
		credentials = append(credentials, &credential)
	}

	return &types.QueryCredentialsByDIDResponse{
		Credentials: credentials,
	}, nil
}

// Params returns the module parameters
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{
		Params: k.GetParams(ctx),
	}, nil
}

// DIDDocuments returns a DID document by DID
func (k Keeper) DIDDocuments(goCtx context.Context, req *types.QueryDIDDocumentsRequest) (*types.QueryDIDDocumentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	document, found := k.GetDIDDocument(ctx, req.Did)
	if !found {
		return nil, status.Error(codes.NotFound, "DID document not found")
	}

	return &types.QueryDIDDocumentsResponse{
		Document: &document,
	}, nil
}
