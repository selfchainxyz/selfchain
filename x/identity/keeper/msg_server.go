package keeper

import (
	"context"
	"fmt"
	"time"

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

// CreateDID creates a new DID document
func (k msgServer) CreateDID(goCtx context.Context, msg *types.MsgCreateDID) (*types.MsgCreateDIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if DID already exists
	if k.HasDIDDocument(ctx, msg.Did) {
		return nil, fmt.Errorf("DID already exists: %s", msg.Did)
	}

	// Convert verification methods
	verificationMethods := make([]types.VerificationMethod, len(msg.VerificationMethods))
	for i, vm := range msg.VerificationMethods {
		verificationMethods[i] = *vm
	}

	now := time.Now().UTC()

	// Create DID document
	doc := types.DIDDocument{
		Id:                  msg.Did,
		VerificationMethods: verificationMethods,
		Authentication:      msg.Authentication,
		Created:            &now,
		Updated:            &now,
	}

	// Store DID document
	k.SetDIDDocument(ctx, msg.Did, doc)

	return &types.MsgCreateDIDResponse{}, nil
}

// UpdateDID updates an existing DID document
func (k msgServer) UpdateDID(goCtx context.Context, msg *types.MsgUpdateDID) (*types.MsgUpdateDIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if DID exists
	doc, found := k.GetDIDDocument(ctx, msg.Did)
	if !found {
		return nil, fmt.Errorf("DID not found: %s", msg.Did)
	}

	// Convert verification methods
	verificationMethods := make([]types.VerificationMethod, len(msg.VerificationMethods))
	for i, vm := range msg.VerificationMethods {
		verificationMethods[i] = *vm
	}

	now := time.Now().UTC()

	// Update DID document
	doc.VerificationMethods = verificationMethods
	doc.Authentication = msg.Authentication
	doc.Updated = &now

	// Store updated DID document
	k.SetDIDDocument(ctx, msg.Did, doc)

	return &types.MsgUpdateDIDResponse{}, nil
}

// CreateCredential creates a new credential
func (k msgServer) CreateCredential(goCtx context.Context, msg *types.MsgCreateCredential) (*types.MsgCreateCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if credential already exists
	if k.HasCredential(ctx, msg.Id) {
		return nil, fmt.Errorf("credential already exists: %s", msg.Id)
	}

	now := time.Now().UTC()

	// Create credential
	cred := types.Credential{
		Id:        msg.Id,
		Issuer:    msg.Issuer,
		Subject:   msg.Subject,
		Claims:    msg.Claims,
		SchemaId:  msg.SchemaId,
		Created:   &now,
		Revoked:   false,
	}

	// Store credential
	k.SetCredential(ctx, cred)

	return &types.MsgCreateCredentialResponse{}, nil
}

// RevokeCredential revokes a credential
func (k msgServer) RevokeCredential(goCtx context.Context, msg *types.MsgRevokeCredential) (*types.MsgRevokeCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if credential exists
	cred, found := k.GetCredential(ctx, msg.Id)
	if !found {
		return nil, fmt.Errorf("credential not found: %s", msg.Id)
	}

	// Check if issuer matches
	if cred.Issuer != msg.Issuer {
		return nil, fmt.Errorf("unauthorized: only issuer can revoke credential")
	}

	// Revoke credential
	cred.Revoked = true
	k.SetCredential(ctx, cred)

	return &types.MsgRevokeCredentialResponse{}, nil
}

// LinkSocialIdentity links a social identity to a DID
func (k msgServer) LinkSocialIdentity(goCtx context.Context, msg *types.MsgLinkSocialIdentity) (*types.MsgLinkSocialIdentityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify OAuth token
	socialId, err := k.VerifyOAuthToken(ctx, msg.Provider, msg.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to verify OAuth token: %w", err)
	}

	now := time.Now().UTC()

	// Create social identity
	identity := types.SocialIdentity{
		Did:      msg.Did,
		Provider: msg.Provider,
		Id:       socialId,
		Verified: true,
		Created:  &now,
	}

	// Store social identity
	if err := k.StoreSocialIdentity(ctx, identity); err != nil {
		return nil, fmt.Errorf("failed to store social identity: %w", err)
	}

	return &types.MsgLinkSocialIdentityResponse{}, nil
}

var _ types.MsgServer = msgServer{}
