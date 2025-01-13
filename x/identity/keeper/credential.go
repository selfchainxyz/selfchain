package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"selfchain/x/identity/types"
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

// CreateCredential creates a new credential
func (k Keeper) CreateCredential(ctx context.Context, credential *types.Credential) error {
	if credential == nil {
		return sdkerrors.Register(types.ModuleName, 1100, "credential cannot be nil")
	}

	if err := credential.ValidateBasic(); err != nil {
		return sdkerrors.Wrapf(err, "invalid credential")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(sdkCtx.KVStore(k.storeKey), []byte(types.CredentialPrefix))

	// Check if credential with same ID already exists
	if store.Has([]byte(credential.Id)) {
		return sdkerrors.Register(types.ModuleName, 1101, "credential already exists")
	}

	// Set issuance date if not set
	if credential.IssuanceDate == 0 {
		credential.IssuanceDate = time.Now().Unix()
	}

	// Set expiration date if provided
	if credential.ExpirationDate != 0 {
		// Ensure expiration date is in the future
		if time.Unix(credential.ExpirationDate, 0).Before(time.Now()) {
			return sdkerrors.Register(types.ModuleName, 1102, "expiration date must be in the future")
		}
	}

	// Set initial status to active
	credential.Status = string(types.CredentialStatusActive)

	// Create audit log
	now := time.Now().Unix()
	auditLog := types.AuditLogEntry{
		Id:        k.GenerateAuditLogID(credential.Id),
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_CREATED,
		Did:       credential.Id,
		Action:    "create_credential",
		Actor:     credential.Issuer,
		Timestamp: now,
		Severity:  types.SecuritySeverity_SEVERITY_INFO,
	}

	if err := k.StoreAuditLog(sdkCtx, auditLog); err != nil {
		return sdkerrors.Wrapf(err, "failed to create audit log")
	}

	// Store the credential
	bz, err := k.cdc.Marshal(credential)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to marshal credential")
	}

	store.Set([]byte(credential.Id), bz)
	return nil
}

// GetCredential returns a credential by ID
func (k Keeper) GetCredential(ctx context.Context, id string) (*types.Credential, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(sdkCtx.KVStore(k.storeKey), []byte(types.CredentialPrefix))
	bz := store.Get([]byte(id))
	if bz == nil {
		return nil, sdkerrors.Register(types.ModuleName, 1103, "credential not found")
	}

	var credential types.Credential
	if err := k.cdc.Unmarshal(bz, &credential); err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to unmarshal credential")
	}

	return &credential, nil
}

// ListCredentials returns all credentials
func (k Keeper) ListCredentials(ctx context.Context, pagination *query.PageRequest) ([]*types.Credential, *query.PageResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(sdkCtx.KVStore(k.storeKey), []byte(types.CredentialPrefix))

	var credentials []*types.Credential
	pageRes, err := query.Paginate(store, pagination, func(key []byte, value []byte) error {
		var credential types.Credential
		if err := k.cdc.Unmarshal(value, &credential); err != nil {
			return err
		}

		credentials = append(credentials, &credential)
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return credentials, pageRes, nil
}

// UpdateCredentialStatus updates the status of a credential
func (k Keeper) UpdateCredentialStatus(ctx context.Context, id string, newStatus string) error {
	credential, err := k.GetCredential(ctx, id)
	if err != nil {
		return err
	}

	// Check if the credential is expired
	now := time.Now()
	if credential.ExpirationDate != 0 && time.Unix(credential.ExpirationDate, 0).Before(now) {
		return sdkerrors.Register(types.ModuleName, 1106, "cannot update expired credential")
	}

	// Validate the new status
	if !types.CredentialStatus(newStatus).IsValid() {
		return sdkerrors.Register(types.ModuleName, 1107, fmt.Sprintf("invalid credential status: %s", newStatus))
	}

	// Cannot reactivate a revoked credential
	if credential.Status == string(types.CredentialStatusRevoked) && newStatus == string(types.CredentialStatusActive) {
		return sdkerrors.Register(types.ModuleName, 1108, "cannot reactivate a revoked credential")
	}

	// Create audit log
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	auditLog := types.AuditLogEntry{
		Id:        k.GenerateAuditLogID(credential.Id),
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_UPDATED,
		Did:       credential.Id,
		Action:    fmt.Sprintf("update_status_%s", newStatus),
		Actor:     credential.Issuer,
		Timestamp: time.Now().Unix(),
		Severity:  types.SecuritySeverity_SEVERITY_INFO,
	}

	if err := k.StoreAuditLog(sdkCtx, auditLog); err != nil {
		return sdkerrors.Wrapf(err, "failed to create audit log")
	}

	credential.Status = newStatus

	// Store the updated credential
	store := prefix.NewStore(sdkCtx.KVStore(k.storeKey), []byte(types.CredentialPrefix))
	bz, err := k.cdc.Marshal(credential)
	if err != nil {
		return sdkerrors.Wrapf(err, "failed to marshal credential")
	}

	store.Set([]byte(credential.Id), bz)
	return nil
}

// RevokeCredential revokes a credential
func (k Keeper) RevokeCredential(ctx context.Context, id string) error {
	credential, err := k.GetCredential(ctx, id)
	if err != nil {
		return err
	}

	// Check if already revoked
	if credential.Status == string(types.CredentialStatusRevoked) {
		return sdkerrors.Register(types.ModuleName, 1109, "credential already revoked")
	}

	// Create audit log
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	auditLog := types.AuditLogEntry{
		Id:        k.GenerateAuditLogID(credential.Id),
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_REVOKED,
		Did:       credential.Id,
		Action:    "revoke_credential",
		Actor:     credential.Issuer,
		Timestamp: time.Now().Unix(),
		Severity:  types.SecuritySeverity_SEVERITY_INFO,
	}

	if err := k.StoreAuditLog(sdkCtx, auditLog); err != nil {
		return sdkerrors.Wrapf(err, "failed to create audit log")
	}

	return k.UpdateCredentialStatus(ctx, id, string(types.CredentialStatusRevoked))
}

// VerifyCredential verifies a credential
func (k Keeper) VerifyCredential(ctx context.Context, id string) (bool, error) {
	credential, err := k.GetCredential(ctx, id)
	if err != nil {
		return false, err
	}

	// Check if revoked
	if credential.Status == string(types.CredentialStatusRevoked) {
		return false, sdkerrors.Register(types.ModuleName, 1110, "credential is revoked")
	}

	// Check if expired
	now := time.Now()
	if credential.ExpirationDate != 0 && time.Unix(credential.ExpirationDate, 0).Before(now) {
		return false, sdkerrors.Register(types.ModuleName, 1111, "credential is expired")
	}

	// Verify the credential proof if present
	if credential.Proof != nil {
		if err := credential.Proof.ValidateBasic(); err != nil {
			return false, sdkerrors.Wrapf(err, "invalid credential proof")
		}
	}

	// Create audit log
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	auditLog := types.AuditLogEntry{
		Id:        k.GenerateAuditLogID(credential.Id),
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_VERIFIED,
		Did:       credential.Id,
		Action:    "verify_credential",
		Actor:     credential.Issuer,
		Timestamp: time.Now().Unix(),
		Severity:  types.SecuritySeverity_SEVERITY_INFO,
	}

	if err := k.StoreAuditLog(sdkCtx, auditLog); err != nil {
		return false, sdkerrors.Wrapf(err, "failed to create audit log")
	}

	return true, nil
}

// HasCredential checks if a credential exists
func (k Keeper) HasCredential(ctx context.Context, id string) bool {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(sdkCtx.KVStore(k.storeKey), []byte(types.CredentialPrefix))
	return store.Has([]byte(id))
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

// GenerateAuditLogID generates a unique ID for an audit log entry
func (k Keeper) GenerateAuditLogID(credentialID string) string {
	hasher := sha256.New()
	hasher.Write([]byte(credentialID + time.Now().String()))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetCredentialsByDID retrieves all credentials for a given DID
func (k Keeper) GetCredentialsByDID(ctx context.Context, did string) ([]*types.Credential, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(sdkCtx.KVStore(k.storeKey), []byte(types.CredentialPrefix))
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

// GetAllCredentials returns all credentials
func (k Keeper) GetAllCredentials(ctx context.Context) ([]types.Credential, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(sdkCtx.KVStore(k.storeKey), []byte(types.CredentialPrefix))
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
func (k Keeper) VerifyPresentation(ctx context.Context, presentation *types.CredentialPresentation, verifierDID string) (bool, error) {
	// Verify the original credential exists and is valid
	valid, err := k.VerifyCredential(ctx, presentation.VerifiableCredential)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, fmt.Errorf("credential verification failed")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Verify presentation proof
	valid, err = k.verifyPresentationProof(sdkCtx, presentation)
	if err != nil {
		return false, fmt.Errorf("failed to verify presentation proof: %v", err)
	}

	// Create audit log for presentation verification
	now := time.Now().Unix()
	auditLog := types.AuditLogEntry{
		Id:        k.GenerateAuditLogID(presentation.VerifiableCredential),
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_VERIFIED,
		Did:       presentation.VerifiableCredential,
		Action:    "verify_presentation",
		Actor:     verifierDID,
		Timestamp: now,
		Severity:  types.SecuritySeverity_SEVERITY_INFO,
		Metadata: map[string]string{
			"credential_id": presentation.VerifiableCredential,
		},
	}
	if err := k.StoreAuditLog(sdkCtx, auditLog); err != nil {
		return false, err
	}

	return valid, nil
}

// createPresentationProof creates a proof for a credential presentation
func (k Keeper) createPresentationProof(ctx sdk.Context, presentation *types.CredentialPresentation) (*types.CredentialProof, error) {
	credential, err := k.GetCredential(ctx, presentation.VerifiableCredential)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	// Create a new proof
	now := time.Now().Unix()

	// Create the presentation proof
	proof := &types.CredentialProof{
		Type:               "ZKProof",
		Created:            now,
		VerificationMethod: credential.Subject, // Use subject DID as verification method
		ProofPurpose:       "credentialPresentation",
	}

	// Create audit log for proof generation
	auditLog := types.AuditLogEntry{
		Id:        k.GenerateAuditLogID(presentation.VerifiableCredential),
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_UPDATED,
		Did:       credential.Subject,
		Action:    "generate_presentation_proof",
		Actor:     credential.Subject,
		Timestamp: now,
		Severity:  types.SecuritySeverity_SEVERITY_INFO,
		Metadata: map[string]string{
			"credential_id": presentation.VerifiableCredential,
			"proof_type":    proof.Type,
		},
	}

	if err := k.StoreAuditLog(ctx, auditLog); err != nil {
		return nil, fmt.Errorf("failed to create audit log: %w", err)
	}

	return proof, nil
}

// verifyPresentationProof verifies a presentation's proof
func (k Keeper) verifyPresentationProof(ctx sdk.Context, presentation *types.CredentialPresentation) (bool, error) {
	// Get the proof from the presentation
	proof := presentation.Proof
	if proof == nil {
		return false, fmt.Errorf("presentation proof is missing")
	}

	// Get the original credential to verify against
	credential, err := k.GetCredential(ctx, presentation.VerifiableCredential)
	if err != nil {
		return false, fmt.Errorf("failed to get credential: %w", err)
	}

	// Verify the proof
	valid := true // Implement actual verification logic here

	// Create audit log for proof verification
	now := time.Now().Unix()
	auditLog := types.AuditLogEntry{
		Id:        k.GenerateAuditLogID(presentation.VerifiableCredential),
		Type:      types.AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_VERIFIED,
		Did:       credential.Subject,
		Action:    "verify_presentation_proof",
		Actor:     credential.Subject,
		Timestamp: now,
		Severity:  types.SecuritySeverity_SEVERITY_INFO,
		Metadata: map[string]string{
			"credential_id": presentation.VerifiableCredential,
			"proof_type":    proof.Type,
			"valid":         fmt.Sprintf("%t", valid),
		},
	}

	if err := k.StoreAuditLog(ctx, auditLog); err != nil {
		return false, fmt.Errorf("failed to create audit log: %w", err)
	}

	return valid, nil
}

// SetCredential stores a credential in the module's KV store
func (k Keeper) SetCredential(ctx sdk.Context, credential *types.Credential) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.CredentialPrefix))
	bz := k.cdc.MustMarshal(credential)
	store.Set([]byte(credential.Id), bz)
}
