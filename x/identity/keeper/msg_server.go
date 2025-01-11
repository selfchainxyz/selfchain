package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateDIDDocument creates a new DID document
func (k msgServer) CreateDIDDocument(goCtx context.Context, msg *types.MsgCreateDIDDocument) (*types.MsgCreateDIDDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the DID already exists
	if k.HasDIDDocument(ctx, msg.Document.Id) {
		return nil, fmt.Errorf("DID document already exists: %s", msg.Document.Id)
	}

	// Store the DID document
	id := k.StoreDIDDocument(ctx, msg.Document)

	return &types.MsgCreateDIDDocumentResponse{
		Id: id,
	}, nil
}

// UpdateDIDDocument updates an existing DID document
func (k msgServer) UpdateDIDDocument(goCtx context.Context, msg *types.MsgUpdateDIDDocument) (*types.MsgUpdateDIDDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the DID document exists
	if !k.HasDIDDocument(ctx, msg.Document.Id) {
		return nil, fmt.Errorf("DID document does not exist: %s", msg.Document.Id)
	}

	// Update the DID document
	k.StoreDIDDocument(ctx, msg.Document)

	return &types.MsgUpdateDIDDocumentResponse{
		Success: true,
	}, nil
}

// IssueCredential issues a new verifiable credential
func (k msgServer) IssueCredential(goCtx context.Context, msg *types.MsgIssueCredential) (*types.MsgIssueCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the credential already exists
	if k.HasCredential(ctx, msg.Credential.Id) {
		return nil, fmt.Errorf("credential already exists: %s", msg.Credential.Id)
	}

	// Store the credential
	err := k.SetCredential(ctx, msg.Credential)
	if err != nil {
		return nil, fmt.Errorf("failed to store credential: %w", err)
	}

	return &types.MsgIssueCredentialResponse{
		Id: msg.Credential.Id,
	}, nil
}

// RevokeCredential revokes an existing credential
func (k msgServer) RevokeCredential(goCtx context.Context, msg *types.MsgRevokeCredential) (*types.MsgRevokeCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the credential exists
	if !k.HasCredential(ctx, msg.CredentialId) {
		return nil, fmt.Errorf("credential does not exist: %s", msg.CredentialId)
	}

	// Get the credential
	cred, found := k.GetCredential(ctx, msg.CredentialId)
	if !found {
		return nil, fmt.Errorf("credential not found: %s", msg.CredentialId)
	}

	// Update credential status
	cred.Revoked = true
	cred.IssuedAt = ctx.BlockTime()

	// Store updated credential
	err := k.SetCredential(ctx, cred)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke credential: %w", err)
	}

	return &types.MsgRevokeCredentialResponse{
		Success: true,
	}, nil
}

// LinkSocialIdentity links a social identity to a DID
func (k msgServer) LinkSocialIdentity(goCtx context.Context, msg *types.MsgLinkSocialIdentity) (*types.MsgLinkSocialIdentityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the OAuth token
	socialID, err := k.VerifyOAuthToken(ctx, msg.Provider, msg.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to verify OAuth token: %w", err)
	}

	// Create and store the social identity
	identity := types.SocialIdentity{
		Id:         socialID,
		Did:        msg.Did,
		VerifiedAt: ctx.BlockTime().Unix(),
	}

	// Store the social identity
	k.StoreSocialIdentity(ctx, identity)

	return &types.MsgLinkSocialIdentityResponse{
		Success: true,
	}, nil
}
