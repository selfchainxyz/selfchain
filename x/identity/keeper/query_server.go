package keeper

import (
	"context"

	"selfchain/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// queryServer implements the QueryServer interface
type queryServer struct {
	k Keeper
}

// NewQueryServer creates a new query server instance
func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServer{k}
}

// DIDDocuments implements types.QueryServer
func (q queryServer) DIDDocuments(ctx context.Context, req *types.QueryDIDDocumentsRequest) (*types.QueryDIDDocumentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Did == "" {
		return nil, status.Error(codes.InvalidArgument, "DID cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	document, found := q.k.GetDIDDocument(sdkCtx, req.Did)
	if !found {
		return nil, status.Error(codes.NotFound, "DID document not found")
	}

	return &types.QueryDIDDocumentsResponse{
		Document: &document,
	}, nil
}

// Credential implements types.QueryServer
func (q queryServer) Credential(ctx context.Context, req *types.QueryCredentialRequest) (*types.QueryCredentialResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	credential, err := q.k.GetCredential(sdkCtx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryCredentialResponse{
		Credential: credential,
	}, nil
}

// CredentialsByDID implements types.QueryServer
func (q queryServer) CredentialsByDID(ctx context.Context, req *types.QueryCredentialsByDIDRequest) (*types.QueryCredentialsByDIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Did == "" {
		return nil, status.Error(codes.InvalidArgument, "DID cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	credentials, err := q.k.GetCredentialsByDID(sdkCtx, req.Did)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCredentialsByDIDResponse{
		Credentials: credentials,
	}, nil
}

// SocialIdentities implements types.QueryServer
func (q queryServer) SocialIdentities(ctx context.Context, req *types.QuerySocialIdentitiesRequest) (*types.QuerySocialIdentitiesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	identities, err := q.k.GetSocialIdentities(sdkCtx, req.Did)
	if err != nil {
		return nil, err
	}

	return &types.QuerySocialIdentitiesResponse{
		Identities: identities,
	}, nil
}

// SocialIdentity implements types.QueryServer
func (q queryServer) SocialIdentity(ctx context.Context, req *types.QuerySocialIdentityRequest) (*types.QuerySocialIdentityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Did == "" || req.Provider == "" {
		return nil, status.Error(codes.InvalidArgument, "DID and provider cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	identity, found := q.k.GetSocialIdentity(sdkCtx, req.Did, req.Provider)
	if !found {
		return nil, status.Error(codes.NotFound, "social identity not found")
	}

	return &types.QuerySocialIdentityResponse{
		SocialId: identity.ProviderId,
	}, nil
}

// MFAConfig implements types.QueryServer
func (q queryServer) MFAConfig(ctx context.Context, req *types.QueryMFAConfigRequest) (*types.QueryMFAConfigResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	config, found := q.k.GetMFAConfig(sdkCtx, req.Did)
	if !found {
		return nil, status.Error(codes.NotFound, "MFA config not found")
	}

	return &types.QueryMFAConfigResponse{
		Config: config,
	}, nil
}

// MFAChallenge implements types.QueryServer
func (q queryServer) MFAChallenge(ctx context.Context, req *types.QueryMFAChallengeRequest) (*types.QueryMFAChallengeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	challenge, err := q.k.CreateMFAChallenge(sdkCtx, req.Did, req.Method)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMFAChallengeResponse{
		Challenge: challenge,
	}, nil
}

// LinkedDID implements types.QueryServer
func (q queryServer) LinkedDID(ctx context.Context, req *types.QueryLinkedDIDRequest) (*types.QueryLinkedDIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Provider == "" {
		return nil, status.Error(codes.InvalidArgument, "provider cannot be empty")
	}

	if req.SocialId == "" {
		return nil, status.Error(codes.InvalidArgument, "social ID cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	did, found := q.k.GetLinkedDID(sdkCtx, req.Provider, req.SocialId)
	if !found {
		return nil, status.Error(codes.NotFound, "no DID found for the given social identity")
	}

	return &types.QueryLinkedDIDResponse{
		Did: did,
	}, nil
}

// Params implements types.QueryServer
func (q queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	return &types.QueryParamsResponse{Params: q.k.GetParams(sdkCtx)}, nil
}

// SocialIdentityBySocialID implements types.QueryServer
func (q queryServer) SocialIdentityBySocialID(ctx context.Context, req *types.QuerySocialIdentityBySocialIDRequest) (*types.QuerySocialIdentityBySocialIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Provider == "" {
		return nil, status.Error(codes.InvalidArgument, "provider cannot be empty")
	}

	if req.SocialId == "" {
		return nil, status.Error(codes.InvalidArgument, "social ID cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	identity, found := q.k.GetSocialIdentityBySocialID(sdkCtx, req.Provider, req.SocialId)
	if !found {
		return nil, status.Error(codes.NotFound, "social identity not found")
	}

	return &types.QuerySocialIdentityBySocialIDResponse{
		Identity: identity,
	}, nil
}
