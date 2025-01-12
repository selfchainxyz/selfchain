package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"selfchain/x/identity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	SchemaPrefix = "schema:"
)

// CredentialSchema represents a schema for validating credentials
type CredentialSchema struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Properties  map[string]SchemaField `json:"properties"`
	Required    bool                   `json:"required"`
}

// SchemaField represents a field in the credential schema
type SchemaField struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// SetCredential stores a credential in the keeper
func (k Keeper) SetCredential(ctx sdk.Context, credential types.Credential) error {
	store := ctx.KVStore(k.storeKey)
	now := ctx.BlockTime()

	// Set issuance date if not set
	if credential.IssuanceDate == nil {
		credential.IssuanceDate = &now
	}

	// Set expiration date if not set
	if credential.ExpirationDate == nil {
		// Default expiration is 1 year from issuance
		expiry := now.Add(365 * 24 * time.Hour)
		credential.ExpirationDate = &expiry
	}

	// Set initial status
	credential.Status = types.CredentialStatus_ACTIVE

	key := append([]byte(types.CredentialPrefix), []byte(credential.Id)...)
	bz, err := k.cdc.Marshal(&credential)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal credential")
	}

	store.Set(key, bz)

	// Create audit log
	err = k.CreateAuditLog(ctx, types.AuditLogEntry{
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_UPDATED,
		Did:       credential.Subject,
		Action:    "update_credential",
		Actor:     credential.Issuer,
		Metadata:  map[string]string{"credential_id": credential.Id},
		Timestamp: &now,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetCredential retrieves a credential from the keeper
func (k Keeper) GetCredential(ctx sdk.Context, id string) (*types.Credential, error) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.CredentialPrefix), []byte(id)...)
	bz := store.Get(key)
	if bz == nil {
		return nil, status.Error(codes.NotFound, "credential not found")
	}

	var credential types.Credential
	err := k.cdc.Unmarshal(bz, &credential)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to unmarshal credential")
	}

	return &credential, nil
}

// DeleteCredential deletes a credential from the keeper
func (k Keeper) DeleteCredential(ctx sdk.Context, id string) error {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.CredentialPrefix), []byte(id)...)

	if !store.Has(key) {
		return status.Error(codes.NotFound, "credential not found")
	}

	store.Delete(key)

	// Create audit log
	now := ctx.BlockTime()
	err := k.CreateAuditLog(ctx, types.AuditLogEntry{
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_DELETED,
		Did:       id,
		Action:    "delete_credential",
		Actor:     id,
		Metadata:  map[string]string{"credential_id": id},
		Timestamp: &now,
	})
	if err != nil {
		return err
	}

	return nil
}

// IsCredentialValid checks if a credential is valid
func (k Keeper) IsCredentialValid(ctx sdk.Context, credential *types.Credential) bool {
	if credential == nil {
		return false
	}

	// Check if credential is expired
	if credential.ExpirationDate != nil {
		expiryTime := *credential.ExpirationDate
		if expiryTime.Before(ctx.BlockTime()) {
			return false
		}
	}

	// Check if credential is revoked
	if credential.Status == types.CredentialStatus_REVOKED {
		return false
	}

	return true
}

// VerifyCredential verifies a credential's validity
func (k Keeper) VerifyCredential(ctx sdk.Context, id string) error {
	credential, err := k.GetCredential(ctx, id)
	if err != nil {
		return err
	}

	if !k.IsCredentialValid(ctx, credential) {
		return sdkerrors.Wrap(types.ErrInvalidCredential, "credential is not valid")
	}

	// Create audit log
	now := ctx.BlockTime()
	err = k.CreateAuditLog(ctx, types.AuditLogEntry{
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_UPDATED,
		Did:       credential.Subject,
		Action:    "verify_credential",
		Actor:     credential.Issuer,
		Metadata:  map[string]string{"credential_id": credential.Id},
		Timestamp: &now,
	})
	if err != nil {
		return err
	}

	return nil
}

// GenerateAuditLogID generates a unique ID for an audit log entry
func (k Keeper) GenerateAuditLogID(credentialID string) string {
	hasher := sha256.New()
	hasher.Write([]byte(credentialID + time.Now().String()))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetCredentialsByDID retrieves all credentials for a given DID
func (k Keeper) GetCredentialsByDID(ctx sdk.Context, did string) ([]*types.Credential, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CredentialPrefix))
	defer iterator.Close()

	var credentials []*types.Credential
	for ; iterator.Valid(); iterator.Next() {
		var credential types.Credential
		err := k.cdc.Unmarshal(iterator.Value(), &credential)
		if err != nil {
			continue
		}

		if credential.Subject == did || credential.Issuer == did {
			credentials = append(credentials, &credential)
		}
	}

	return credentials, nil
}

// RevokeCredential revokes a credential
func (k Keeper) RevokeCredential(ctx sdk.Context, id string) error {
	store := ctx.KVStore(k.storeKey)

	// Get existing credential
	credential, err := k.GetCredential(ctx, id)
	if err != nil {
		return sdkerrors.Wrap(types.ErrCredentialNotFound, "credential does not exist")
	}

	// Check if already revoked
	if credential.Status == types.CredentialStatus_REVOKED {
		return sdkerrors.Wrap(types.ErrCredentialAlreadyRevoked, "credential is already revoked")
	}

	// Update status to revoked
	credential.Status = types.CredentialStatus_REVOKED

	// Store updated credential
	bz, err := k.cdc.Marshal(credential)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal credential")
	}
	store.Set([]byte(types.CredentialPrefix+credential.Id), bz)

	// Create audit log
	now := ctx.BlockTime()
	err = k.CreateAuditLog(ctx, types.AuditLogEntry{
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_REVOKED,
		Did:       credential.Subject,
		Action:    "revoke_credential",
		Actor:     credential.Issuer,
		Metadata:  map[string]string{"credential_id": credential.Id},
		Timestamp: &now,
	})
	if err != nil {
		return err
	}

	return nil
}

// HasCredential checks if a credential exists
func (k Keeper) HasCredential(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.CredentialPrefix), []byte(id)...)
	return store.Has(key)
}

// GetAllCredentials returns all credentials
func (k Keeper) GetAllCredentials(ctx sdk.Context) ([]types.Credential, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CredentialPrefix))
	defer iterator.Close()

	var credentials []types.Credential
	for ; iterator.Valid(); iterator.Next() {
		var credential types.Credential
		err := k.cdc.Unmarshal(iterator.Value(), &credential)
		if err != nil {
			continue
		}
		credentials = append(credentials, credential)
	}

	return credentials, nil
}

// VerifyPresentation verifies a credential presentation
func (k Keeper) VerifyPresentation(ctx sdk.Context, presentation *types.CredentialPresentation, verifierDID string) (bool, error) {
	// Verify the original credential exists and is valid
	err := k.VerifyCredential(ctx, presentation.CredentialId)
	if err != nil {
		return false, err
	}

	// Verify presentation proof
	valid, err := k.verifyPresentationProof(ctx, presentation)
	if err != nil {
		return false, fmt.Errorf("failed to verify presentation proof: %v", err)
	}

	// Create audit log for presentation verification
	now := ctx.BlockTime()
	err = k.CreateAuditLog(ctx, types.AuditLogEntry{
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_UPDATED,
		Did:       presentation.CredentialId,
		Action:    "verify_presentation",
		Actor:     verifierDID,
		Metadata:  map[string]string{"credential_id": presentation.CredentialId},
		Timestamp: &now,
	})
	if err != nil {
		return false, err
	}

	return valid, nil
}

// createPresentationProof creates a proof for a credential presentation
func (k Keeper) createPresentationProof(ctx sdk.Context, presentation *types.CredentialPresentation) (*types.CredentialProof, error) {
	now := time.Now()
	proof := &types.CredentialProof{
		Type:               "Ed25519Signature2020",
		Created:            &now,
		VerificationMethod: "did:selfchain:key1",
		ProofPurpose:       "assertionMethod",
	}
	return proof, nil
}

// verifyPresentationProof verifies a presentation's proof
func (k Keeper) verifyPresentationProof(ctx sdk.Context, presentation *types.CredentialPresentation) (bool, error) {
	// TODO: Implement actual presentation proof verification
	// This would involve:
	// 1. Verifying the proof matches the disclosed claims
	// 2. Checking the proof was created by the credential subject
	// 3. Verifying the signature
	return true, nil
}

// GetClaimHash returns a hash of a claim value
func (k Keeper) GetClaimHash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

// VerifyClaimHash verifies if a claim value matches its hash
func (k Keeper) VerifyClaimHash(value string, hash string) bool {
	computedHash := k.GetClaimHash(value)
	return computedHash == hash
}

// CreateAuditLog creates an audit log entry
func (k Keeper) CreateAuditLog(ctx sdk.Context, entry types.AuditLogEntry) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&entry)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal audit log entry")
	}

	store.Set([]byte(entry.Id), bz)
	return nil
}

// GetAuditLogs retrieves audit logs for a credential
func (k Keeper) GetAuditLogs(ctx sdk.Context, id string) ([]types.AuditLogEntry, error) {
	store := ctx.KVStore(k.storeKey)
	prefix := append([]byte(types.AuditLogPrefix), []byte(id)...)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var logs []types.AuditLogEntry
	for ; iterator.Valid(); iterator.Next() {
		var log types.AuditLogEntry
		k.cdc.MustUnmarshal(iterator.Value(), &log)
		if log.Metadata["credential_id"] == id {
			logs = append(logs, log)
		}
	}
	return logs, nil
}

// CreatePresentation creates a selective disclosure presentation of a credential
func (k Keeper) CreatePresentation(ctx sdk.Context, credentialID string, claimsToDisclose []string) (*types.CredentialPresentation, error) {
	found := k.HasCredential(ctx, credentialID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrCredentialNotFound, credentialID)
	}

	credential, err := k.GetCredential(ctx, credentialID)
	if err != nil {
		return nil, err
	}

	now := ctx.BlockTime()
	// Create presentation with only disclosed claims
	presentation := &types.CredentialPresentation{
		CredentialId:     credentialID,
		Verifier:         "",
		DisclosedClaims:  claimsToDisclose,
		PresentationDate: &now,
	}

	// Create proof for presentation
	proof, err := k.createPresentationProof(ctx, presentation)
	if err != nil {
		return nil, err
	}
	presentation.Proof = proof

	// Log presentation creation
	now = ctx.BlockTime()
	err = k.CreateAuditLog(ctx, types.AuditLogEntry{
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_CREATED,
		Did:       credential.Subject,
		Action:    "create_presentation",
		Actor:     credential.Subject,
		Metadata:  map[string]string{"credential_id": credential.Id},
		Timestamp: &now,
	})
	if err != nil {
		return nil, err
	}

	return presentation, nil
}

// Credential implements the Query/Credential gRPC method
func (k Keeper) Credential(ctx context.Context, req *types.QueryCredentialRequest) (*types.QueryCredentialResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	credential, err := k.GetCredential(sdkCtx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryCredentialResponse{
		Credential: credential,
	}, nil
}
