package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"selfchain/x/identity/types"
)

type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
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
func (k queryServer) DIDDocument(goCtx context.Context, req *types.QueryDIDDocumentRequest) (*types.QueryDIDDocumentResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "DID ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	doc, found := k.GetDIDDocument(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "DID document not found")
	}

	return &types.QueryDIDDocumentResponse{
		Document: doc,
	}, nil
}

// Credential returns a credential by ID
func (k queryServer) Credential(goCtx context.Context, req *types.QueryCredentialRequest) (*types.QueryCredentialResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "credential ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	cred, found := k.GetCredential(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "credential not found")
	}

	return &types.QueryCredentialResponse{
		Credential: cred,
	}, nil
}

// SocialIdentity returns a social identity by DID and provider
func (k queryServer) SocialIdentity(goCtx context.Context, req *types.QuerySocialIdentityRequest) (*types.QuerySocialIdentityResponse, error) {
	if req.Did == "" {
		return nil, status.Error(codes.InvalidArgument, "DID cannot be empty")
	}
	if req.Provider == "" {
		return nil, status.Error(codes.InvalidArgument, "provider cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	identity, found := k.GetSocialIdentityByDID(ctx, req.Did, req.Provider)
	if !found {
		return nil, status.Error(codes.NotFound, "social identity not found")
	}

	return &types.QuerySocialIdentityResponse{
		Identity: &identity,
	}, nil
}

// LinkedDID returns the DID linked to a social identity
func (k queryServer) LinkedDID(goCtx context.Context, req *types.QueryLinkedDIDRequest) (*types.QueryLinkedDIDResponse, error) {
	if req.Provider == "" {
		return nil, status.Error(codes.InvalidArgument, "provider cannot be empty")
	}
	if req.SocialId == "" {
		return nil, status.Error(codes.InvalidArgument, "social ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	identity, found := k.GetSocialIdentityBySocialID(ctx, req.Provider, req.SocialId)
	if !found {
		return nil, status.Error(codes.NotFound, "social identity not found")
	}

	return &types.QueryLinkedDIDResponse{
		Did: identity.Did,
	}, nil
}

// VerifyCredential verifies a credential
func (k queryServer) VerifyCredential(goCtx context.Context, req *types.QueryVerifyCredentialRequest) (*types.QueryVerifyCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	status := types.VerificationStatus{
		Valid:        true,
		ErrorMessage: "",
	}

	// Verify credential exists
	cred, found := k.GetCredential(ctx, req.Id)
	if !found {
		status.Valid = false
		status.ErrorMessage = "credential not found"
		return &types.QueryVerifyCredentialResponse{
			Status: status,
		}, nil
	}

	// Verify credential is not revoked
	if cred.Revoked {
		status.Valid = false
		status.ErrorMessage = "credential is revoked"
		return &types.QueryVerifyCredentialResponse{
			Status: status,
		}, nil
	}

	// Verify issuer exists
	if !k.HasDIDDocument(ctx, cred.Issuer) {
		status.Valid = false
		status.ErrorMessage = "issuer DID not found"
		return &types.QueryVerifyCredentialResponse{
			Status: status,
		}, nil
	}

	// Verify subject exists
	if !k.HasDIDDocument(ctx, cred.Subject) {
		status.Valid = false
		status.ErrorMessage = "subject DID not found"
		return &types.QueryVerifyCredentialResponse{
			Status: status,
		}, nil
	}

	// Verify schema
	if err := k.ValidateCredential(ctx, cred); err != nil {
		status.Valid = false
		status.ErrorMessage = fmt.Sprintf("schema validation failed: %v", err)
		return &types.QueryVerifyCredentialResponse{
			Status: status,
		}, nil
	}

	return &types.QueryVerifyCredentialResponse{
		Status: status,
	}, nil
}

var _ types.QueryServer = queryServer{}
