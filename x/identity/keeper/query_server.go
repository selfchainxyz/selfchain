package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"selfchain/x/identity/types"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{Keeper: k}
}

// DIDDocument returns a DID document by DID
func (k Keeper) DIDDocument(goCtx context.Context, req *types.QueryDIDDocumentRequest) (*types.QueryDIDDocumentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	doc, found := k.GetDIDDocument(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "DID document not found")
	}

	return &types.QueryDIDDocumentResponse{Document: doc}, nil
}

// Credential returns a credential by ID
func (k Keeper) Credential(goCtx context.Context, req *types.QueryCredentialRequest) (*types.QueryCredentialResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	cred, found := k.GetCredential(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "credential not found")
	}

	return &types.QueryCredentialResponse{Credential: cred}, nil
}

// SocialIdentity returns a social identity by DID and provider
func (k Keeper) SocialIdentity(goCtx context.Context, req *types.QuerySocialIdentityRequest) (*types.QuerySocialIdentityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	identity, found := k.GetSocialIdentityByDIDAndProvider(ctx, req.Did, req.Provider)
	if !found {
		return nil, status.Error(codes.NotFound, "social identity not found")
	}

	return &types.QuerySocialIdentityResponse{Identity: identity}, nil
}

// LinkedDID returns the DID linked to a social identity
func (k Keeper) LinkedDID(goCtx context.Context, req *types.QueryLinkedDIDRequest) (*types.QueryLinkedDIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
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

	return &types.QueryLinkedDIDResponse{Did: identity.Did}, nil
}

// VerifyCredential verifies a credential
func (k Keeper) VerifyCredential(goCtx context.Context, req *types.QueryVerifyCredentialRequest) (*types.QueryVerifyCredentialResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	cred, found := k.GetCredential(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "credential not found")
	}

	// Verify the credential
	_, found = k.GetDIDDocument(ctx, cred.Issuer)
	if !found {
		return nil, status.Error(codes.NotFound, "issuer DID document not found")
	}

	// TODO: Implement actual credential verification logic
	status := types.VerificationStatus{
		Valid:  true,
		Reason: "credential verification not implemented yet",
	}

	return &types.QueryVerifyCredentialResponse{Status: status}, nil
}

// MFAConfig returns the MFA configuration for a DID
func (k Keeper) MFAConfig(goCtx context.Context, req *types.QueryMFAConfigRequest) (*types.QueryMFAConfigResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	config, found := k.GetMFAConfig(ctx, req.Did)
	if !found {
		return nil, status.Error(codes.NotFound, "MFA config not found")
	}

	return &types.QueryMFAConfigResponse{Config: config}, nil
}

// MFAChallenge returns an MFA challenge by ID
func (k Keeper) MFAChallenge(goCtx context.Context, req *types.QueryMFAChallengeRequest) (*types.QueryMFAChallengeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	challenge, found := k.GetMFAChallenge(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "MFA challenge not found")
	}

	return &types.QueryMFAChallengeResponse{Challenge: challenge}, nil
}

// RecoverySession returns a recovery session by ID
func (k Keeper) RecoverySession(goCtx context.Context, req *types.QueryRecoverySessionRequest) (*types.QueryRecoverySessionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	session, found := k.GetRecoverySession(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "recovery session not found")
	}

	return &types.QueryRecoverySessionResponse{Session: session}, nil
}
