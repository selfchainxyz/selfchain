package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"selfchain/x/identity/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{Keeper: k}
}

// Params returns the module parameters
func (k queryServer) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// DIDDocument returns a DID document by ID
func (k Keeper) DIDDocument(c context.Context, req *types.QueryDIDDocumentRequest) (*types.QueryDIDDocumentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	document, found := k.GetDIDDocument(ctx, req.Did)
	if !found {
		return nil, status.Error(codes.NotFound, "DID document not found")
	}

	return &types.QueryDIDDocumentResponse{Document: document}, nil
}

// Credential returns a credential by ID
func (k Keeper) Credential(c context.Context, req *types.QueryCredentialRequest) (*types.QueryCredentialResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	credential, found := k.GetCredential(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "credential not found")
	}

	return &types.QueryCredentialResponse{Credential: credential}, nil
}

// CredentialsBySubject returns all credentials for a subject DID
func (k Keeper) CredentialsBySubject(c context.Context, req *types.QueryCredentialsBySubjectRequest) (*types.QueryCredentialsBySubjectResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	allCredentials := k.GetCredentialsBySubject(ctx, req.SubjectDid)

	var credentials []types.Credential
	pageRes := &query.PageResponse{}

	if req.Pagination != nil {
		start := int(req.Pagination.Offset)
		end := int(req.Pagination.Offset + req.Pagination.Limit)

		if start < len(allCredentials) {
			if end > len(allCredentials) {
				end = len(allCredentials)
			}
			credentials = allCredentials[start:end]
		}

		pageRes.Total = uint64(len(allCredentials))
	} else {
		credentials = allCredentials
		pageRes.Total = uint64(len(credentials))
	}

	return &types.QueryCredentialsBySubjectResponse{
		Credentials: credentials,
		Pagination: pageRes,
	}, nil
}

// VerificationStatus returns the verification status for a DID
func (k Keeper) VerificationStatus(c context.Context, req *types.QueryVerificationStatusRequest) (*types.QueryVerificationStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	verificationStatus, found := k.GetVerificationStatus(ctx, req.Did)
	if !found {
		return nil, status.Error(codes.NotFound, "verification status not found")
	}

	return &types.QueryVerificationStatusResponse{Status: verificationStatus}, nil
}

var _ types.QueryServer = Keeper{}
