package keeper

import (
	"context"
	"strings"
	"time"

	"selfchain/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (k msgServer) CreateDID(goCtx context.Context, msg *types.MsgCreateDID) (*types.MsgCreateDIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if DID already exists
	if k.HasDIDDocument(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(types.ErrDIDAlreadyExists, "DID already exists")
	}

	// Create DID document
	now := time.Now()
	didDoc := types.DIDDocument{
		Id:                 msg.Id,
		Controller:         []string{msg.Creator},
		VerificationMethod: msg.VerificationMethod,
		Service:            msg.Service,
		Created:            &now,
		Updated:            &now,
	}

	// Store DID document
	if err := k.SetDIDDocument(ctx, msg.Id, didDoc); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to store DID document")
	}

	return &types.MsgCreateDIDResponse{}, nil
}

func (k msgServer) DeleteDID(goCtx context.Context, msg *types.MsgDeleteDID) (*types.MsgDeleteDIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if DID exists
	didDoc, found := k.GetDIDDocument(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrDIDNotFound, "DID not found")
	}

	// Check if the sender is the controller
	isController := false
	for _, controller := range didDoc.Controller {
		if controller == msg.Creator {
			isController = true
			break
		}
	}
	if !isController {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not the DID controller")
	}

	k.DeleteDIDDocument(ctx, msg.Id)

	return &types.MsgDeleteDIDResponse{}, nil
}

func (k msgServer) CreateCredential(goCtx context.Context, msg *types.MsgCreateCredential) (*types.MsgCreateCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if credential already exists
	if k.HasCredential(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(types.ErrInvalidCredentialID, "credential already exists")
	}

	// Check if issuer DID exists
	if !k.HasDIDDocument(ctx, msg.Issuer) {
		return nil, sdkerrors.Wrap(types.ErrDIDNotFound, "issuer DID not found")
	}

	// Check if subject DID exists
	if !k.HasDIDDocument(ctx, msg.Subject) {
		return nil, sdkerrors.Wrap(types.ErrDIDNotFound, "subject DID not found")
	}

	// Convert []string claims to map[string]string
	claimsMap := make(map[string]string)
	for _, claim := range msg.Claims {
		parts := strings.SplitN(claim, ":", 2)
		if len(parts) == 2 {
			claimsMap[parts[0]] = parts[1]
		}
	}

	// Create credential
	now := time.Now().Unix()
	credential := types.Credential{
		Id:           msg.Id,
		Type:         msg.Type,
		Issuer:       msg.Issuer,
		Subject:      msg.Subject,
		Claims:       claimsMap,
		IssuanceDate: now,
		Status:       string(types.CredentialStatusActive),
	}

	// Store credential
	if err := k.Keeper.CreateCredential(ctx, &credential); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to store credential")
	}

	return &types.MsgCreateCredentialResponse{}, nil
}

func (k msgServer) UpdateCredential(goCtx context.Context, msg *types.MsgUpdateCredential) (*types.MsgUpdateCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if credential exists
	credential, err := k.GetCredential(ctx, msg.Id)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCredentialNotFound, "credential not found")
	}

	// Check if the sender is the issuer
	if credential.Issuer != msg.Creator {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not the credential issuer")
	}

	// Convert []string claims to map[string]string
	claimsMap := make(map[string]string)
	for _, claim := range msg.Claims {
		parts := strings.SplitN(claim, ":", 2)
		if len(parts) == 2 {
			claimsMap[parts[0]] = parts[1]
		}
	}

	// Update credential
	credential.Claims = claimsMap
	now := time.Now().Unix()
	credential.IssuanceDate = now

	// Store updated credential
	if err := k.Keeper.CreateCredential(ctx, credential); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to update credential")
	}

	return &types.MsgUpdateCredentialResponse{}, nil
}

func (k msgServer) RevokeCredential(goCtx context.Context, msg *types.MsgRevokeCredential) (*types.MsgRevokeCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if credential exists
	credential, err := k.GetCredential(ctx, msg.Id)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCredentialNotFound, "credential not found")
	}

	// Check if the sender is the issuer
	if credential.Issuer != msg.Creator {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not the credential issuer")
	}

	// Update credential status
	credential.Status = string(types.CredentialStatusRevoked)
	now := time.Now().Unix()
	credential.IssuanceDate = now

	// Store updated credential
	if err := k.Keeper.CreateCredential(ctx, credential); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to revoke credential")
	}

	return &types.MsgRevokeCredentialResponse{}, nil
}

func (k msgServer) LinkSocialIdentity(goCtx context.Context, msg *types.MsgLinkSocialIdentity) (*types.MsgLinkSocialIdentityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if DID exists
	if !k.HasDIDDocument(ctx, msg.Creator) {
		return nil, sdkerrors.Wrap(types.ErrDIDNotFound, "DID not found")
	}

	// Create social identity
	createdAt := time.Now()
	socialIdentity := types.SocialIdentity{
		Did:       msg.Creator,
		Provider:  msg.Provider,
		Profile:   make(map[string]string),
		CreatedAt: &createdAt,
	}

	// Store social identity
	if err := k.StoreSocialIdentity(ctx, socialIdentity); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to store social identity")
	}

	return &types.MsgLinkSocialIdentityResponse{}, nil
}

func (k msgServer) UnlinkSocialIdentity(goCtx context.Context, msg *types.MsgUnlinkSocialIdentity) (*types.MsgUnlinkSocialIdentityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if social identity exists
	socialIdentity, found := k.Keeper.GetSocialIdentity(ctx, msg.Creator, msg.Provider)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrInvalidSocialIdentity, "social identity not found")
	}

	// Check if the sender is the DID owner
	if socialIdentity.Did != msg.Creator {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not the social identity owner")
	}

	// Remove social identity
	if err := k.DeleteSocialIdentity(ctx, socialIdentity.Did); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to delete social identity")
	}

	return &types.MsgUnlinkSocialIdentityResponse{}, nil
}

func (k msgServer) ConfigureMFA(goCtx context.Context, msg *types.MsgConfigureMFA) (*types.MsgConfigureMFAResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if DID exists
	didDoc, found := k.GetDIDDocument(ctx, msg.Did)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrDIDNotFound, "DID not found")
	}

	// Check if the sender is the DID controller
	isController := false
	for _, controller := range didDoc.Controller {
		if controller == msg.Creator {
			isController = true
			break
		}
	}
	if !isController {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "not the DID controller")
	}

	// Convert string array to MFAMethod array
	methods := make([]*types.MFAMethod, len(msg.Methods))
	for i, method := range msg.Methods {
		methods[i] = &types.MFAMethod{
			Type: method,
		}
	}

	mfaConfig := types.MFAConfig{
		Did:     msg.Did,
		Methods: methods,
	}

	// Store MFA config
	if err := k.StoreMFAConfig(ctx, mfaConfig); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to store MFA config")
	}

	return &types.MsgConfigureMFAResponse{}, nil
}

func (k msgServer) VerifyMFA(goCtx context.Context, msg *types.MsgVerifyMFA) (*types.MsgVerifyMFAResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify MFA challenge
	if err := k.Keeper.VerifyMFAChallenge(ctx, msg.Did, msg.Method, msg.Code); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to verify MFA challenge")
	}

	return &types.MsgVerifyMFAResponse{}, nil
}

func (k msgServer) IssueCredential(goCtx context.Context, msg *types.MsgCreateCredential) (*types.MsgCreateCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if issuer exists
	_, issuerFound := k.Keeper.GetDIDDocument(ctx, msg.Issuer)
	if !issuerFound {
		return nil, sdkerrors.Wrap(types.ErrDIDNotFound, "issuer DID not found")
	}

	// Check if subject exists
	_, subjectFound := k.Keeper.GetDIDDocument(ctx, msg.Subject)
	if !subjectFound {
		return nil, sdkerrors.Wrap(types.ErrDIDNotFound, "subject DID not found")
	}

	// Create credential
	now := time.Now().Unix()
	credential := types.Credential{
		Id:           msg.Id,
		Type:         msg.Type,
		Issuer:       msg.Creator,
		Subject:      msg.Subject,
		Claims:       msg.Claims,
		IssuanceDate: now,
		Status:       string(types.CredentialStatusActive),
	}

	// Store credential
	if err := k.Keeper.CreateCredential(ctx, &credential); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to store credential")
	}

	return &types.MsgCreateCredentialResponse{
		Success: true,
	}, nil
}

func (k msgServer) RemoveMFA(goCtx context.Context, msg *types.MsgRemoveMFA) (*types.MsgRemoveMFAResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if DID exists
	_, found := k.GetDIDDocument(ctx, msg.Did)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrDIDNotFound, msg.Did)
	}

	// Verify ownership
	if msg.Creator != msg.Did {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "not the DID owner")
	}

	if err := k.RemoveMFAMethod(ctx, msg.Did, msg.Method); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to remove MFA method")
	}

	return &types.MsgRemoveMFAResponse{}, nil
}

func (k msgServer) UpdateDID(goCtx context.Context, msg *types.MsgUpdateDID) (*types.MsgUpdateDIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if DID document exists
	_, found := k.GetDIDDocument(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrDIDNotFound, "DID document not found")
	}

	// Create DID document from the message fields
	didDoc := types.DIDDocument{
		Id:                 msg.Id,
		Controller:         msg.Controller,
		VerificationMethod: msg.VerificationMethod,
		Service:            msg.Service,
	}

	// Update DID document
	if err := k.Keeper.SetDIDDocument(ctx, msg.Id, didDoc); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to update DID document")
	}

	return &types.MsgUpdateDIDResponse{
		Success: true,
	}, nil
}

func (k msgServer) AddMFA(goCtx context.Context, msg *types.MsgAddMFA) (*types.MsgAddMFAResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if DID exists
	_, found := k.GetDIDDocument(ctx, msg.Did)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrDIDNotFound, msg.Did)
	}

	// Verify ownership
	if msg.Creator != msg.Did {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "not the DID owner")
	}

	// Create new MFA method
	now := ctx.BlockTime()
	method := &types.MFAMethod{
		Type:      msg.Method,
		Secret:    msg.Secret,
		CreatedAt: &now,
		Status:    types.MFAMethodStatus_MFA_METHOD_STATUS_ACTIVE,
	}

	// Add method to config
	if err := k.Keeper.AddMFAMethod(ctx, msg.Did, *method); err != nil {
		return nil, err
	}

	return &types.MsgAddMFAResponse{}, nil
}
