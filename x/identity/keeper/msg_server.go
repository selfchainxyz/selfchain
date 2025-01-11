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

var _ types.MsgServer = msgServer{}

// RegisterDID implements types.MsgServer
func (k msgServer) RegisterDID(goCtx context.Context, msg *types.MsgRegisterDID) (*types.MsgRegisterDIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.ValidateDIDDocument(msg.Document); err != nil {
		return nil, err
	}

	if err := k.SetDIDDocument(ctx, msg.Document.Id, msg.Document); err != nil {
		return nil, err
	}

	return &types.MsgRegisterDIDResponse{
		Id: msg.Document.Id,
	}, nil
}

// UpdateDID implements types.MsgServer
func (k msgServer) UpdateDID(goCtx context.Context, msg *types.MsgUpdateDID) (*types.MsgUpdateDIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if DID exists
	if _, found := k.GetDIDDocument(ctx, msg.Document.Id); !found {
		return nil, types.ErrDIDNotFound
	}

	if err := k.ValidateDIDDocument(msg.Document); err != nil {
		return nil, err
	}

	if err := k.SetDIDDocument(ctx, msg.Document.Id, msg.Document); err != nil {
		return nil, err
	}

	return &types.MsgUpdateDIDResponse{
		Id: msg.Document.Id,
	}, nil
}

// VerifyIdentity implements types.MsgServer
func (k msgServer) VerifyIdentity(goCtx context.Context, msg *types.MsgVerifyIdentity) (*types.MsgVerifyIdentityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the provider is allowed
	if !k.IsOAuthProviderAllowed(ctx, msg.Provider) {
		return nil, types.ErrUnauthorizedProvider
	}

	// TODO: Implement actual OAuth verification
	// For now, just create a verification record
	verificationID := fmt.Sprintf("verification_%s_%d", msg.Did, ctx.BlockTime().Unix())

	now := ctx.BlockTime()
	expiresAt := now.Add(time.Duration(k.GetVerificationTimeout(ctx)) * time.Hour)

	// Create verification status
	status := types.VerificationStatus{
		Did:        msg.Did,
		Provider:   msg.Provider,
		Verified:   true,
		VerifiedAt: now,
		ExpiresAt:  &expiresAt,
	}

	// Store verification status
	if err := k.SetVerificationStatus(ctx, msg.Did, status); err != nil {
		return nil, err
	}

	return &types.MsgVerifyIdentityResponse{
		VerificationId: verificationID,
	}, nil
}

// IssueCredential implements types.MsgServer
func (k msgServer) IssueCredential(goCtx context.Context, msg *types.MsgIssueCredential) (*types.MsgIssueCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if credential type is allowed
	if !k.IsCredentialTypeAllowed(ctx, msg.CredentialType) {
		return nil, types.ErrUnauthorizedCredentialType
	}

	// Check if subject DID exists
	if _, found := k.GetDIDDocument(ctx, msg.SubjectDid); !found {
		return nil, types.ErrDIDNotFound
	}

	// Check if max credentials per DID is reached
	subjectCredentials := k.GetCredentialsBySubject(ctx, msg.SubjectDid)
	if uint32(len(subjectCredentials)) >= k.GetMaxCredentialsPerDID(ctx) {
		return nil, types.ErrMaxCredentialsReached
	}

	// Create credential
	credentialID := fmt.Sprintf("credential_%s_%d", msg.SubjectDid, ctx.BlockTime().Unix())
	now := ctx.BlockTime()

	credential := types.Credential{
		Id:         credentialID,
		Type:       msg.CredentialType,
		IssuerDid:  msg.Creator, // Using creator address as issuer DID
		SubjectDid: msg.SubjectDid,
		Claims:     msg.Claims,
		IssuedAt:   now,
		ExpiresAt:  msg.Expiry,
		Revoked:    false,
	}

	// Store credential
	if err := k.SetCredential(ctx, credential); err != nil {
		return nil, err
	}

	return &types.MsgIssueCredentialResponse{
		CredentialId: credentialID,
	}, nil
}

// RevokeCredential implements types.MsgServer
func (k msgServer) RevokeCredential(goCtx context.Context, msg *types.MsgRevokeCredential) (*types.MsgRevokeCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if credential exists
	credential, found := k.GetCredential(ctx, msg.CredentialId)
	if !found {
		return nil, types.ErrCredentialNotFound
	}

	// Check if the sender is the issuer
	if credential.IssuerDid != msg.Creator {
		return nil, fmt.Errorf("only the issuer can revoke a credential")
	}

	// Revoke the credential
	credential.Revoked = true
	if err := k.SetCredential(ctx, credential); err != nil {
		return nil, err
	}

	return &types.MsgRevokeCredentialResponse{
		Success: true,
	}, nil
}
