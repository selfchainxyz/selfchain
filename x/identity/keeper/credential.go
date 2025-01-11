package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

const (
	CredentialPrefix = "credential:"
	SchemaPrefix    = "schema:"
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

// SetCredential stores a credential
func (k Keeper) SetCredential(ctx sdk.Context, credential types.Credential) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredentialKey))
	b := k.cdc.MustMarshal(&credential)
	store.Set([]byte(credential.Id), b)
}

// GetCredential returns a credential
func (k Keeper) GetCredential(ctx sdk.Context, id string) (val types.Credential, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredentialKey))
	b := store.Get([]byte(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// HasCredential checks if a credential exists
func (k Keeper) HasCredential(ctx sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredentialKey))
	return store.Has([]byte(id))
}

// GetAllCredentials returns all credentials
func (k Keeper) GetAllCredentials(ctx sdk.Context) (list []types.Credential) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredentialKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Credential
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// ValidateCredential validates a credential against its schema
func (k Keeper) ValidateCredential(ctx sdk.Context, credential types.Credential) error {
	// Check if schema exists
	schema, found := k.GetCredentialSchema(ctx, credential.SchemaId)
	if !found {
		return fmt.Errorf("schema not found: %s", credential.SchemaId)
	}

	// Check if issuer exists
	if !k.HasDIDDocument(ctx, credential.Issuer) {
		return fmt.Errorf("issuer DID not found: %s", credential.Issuer)
	}

	// Check if subject exists
	if !k.HasDIDDocument(ctx, credential.Subject) {
		return fmt.Errorf("subject DID not found: %s", credential.Subject)
	}

	// Validate claims against schema
	for field := range credential.Claims {
		if _, ok := schema.Properties[field]; !ok {
			return fmt.Errorf("field not defined in schema: %s", field)
		}
	}

	// Check for required fields
	if schema.Required {
		for field := range schema.Properties {
			if _, ok := credential.Claims[field]; !ok {
				return fmt.Errorf("required field missing: %s", field)
			}
		}
	}

	return nil
}

// GetCredentialSchema returns a credential schema
func (k Keeper) GetCredentialSchema(ctx sdk.Context, id string) (val types.CredentialSchema, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredentialSchemaKey))
	b := store.Get([]byte(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetCredentialSchema stores a credential schema
func (k Keeper) SetCredentialSchema(ctx sdk.Context, schema types.CredentialSchema) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredentialSchemaKey))
	b := k.cdc.MustMarshal(&schema)
	store.Set([]byte(schema.Id), b)
}

// GetAllCredentialSchemas returns all credential schemas
func (k Keeper) GetAllCredentialSchemas(ctx sdk.Context) (list []types.CredentialSchema) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CredentialSchemaKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CredentialSchema
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
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
