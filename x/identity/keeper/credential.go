package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

const (
	CredentialPrefix = "credential:"
)

// HasCredential checks if a credential exists
func (k Keeper) HasCredential(ctx sdk.Context, credentialID string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(CredentialPrefix))
	return store.Has([]byte(credentialID))
}

// SetCredential stores a verifiable credential
func (k Keeper) SetCredential(ctx sdk.Context, credential types.Credential) error {
	if err := k.ValidateCredential(credential); err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(CredentialPrefix))
	credentialBytes := k.cdc.MustMarshal(&credential)
	store.Set([]byte(credential.Id), credentialBytes)

	// Also store by subject for efficient querying
	subjectStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("credential_by_subject/"))
	subjectKey := []byte(fmt.Sprintf("%s/%s", credential.SubjectDid, credential.Id))
	subjectStore.Set(subjectKey, []byte{1}) // Just store a flag, actual data is in main store

	return nil
}

// GetCredential returns a credential by ID
func (k Keeper) GetCredential(ctx sdk.Context, id string) (types.Credential, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(CredentialPrefix))
	credentialBytes := store.Get([]byte(id))
	if credentialBytes == nil {
		return types.Credential{}, false
	}

	var credential types.Credential
	k.cdc.MustUnmarshal(credentialBytes, &credential)
	return credential, true
}

// GetCredentialsBySubject returns all credentials for a subject DID
func (k Keeper) GetCredentialsBySubject(ctx sdk.Context, subjectDid string) []types.Credential {
	var credentials []types.Credential
	subjectStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("credential_by_subject/"))
	
	prefix := []byte(fmt.Sprintf("%s/", subjectDid))
	iterator := subjectStore.Iterator(prefix, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// Extract credential ID from the key
		key := iterator.Key()
		credentialID := string(key[len(prefix):])
		
		if credential, found := k.GetCredential(ctx, credentialID); found {
			credentials = append(credentials, credential)
		}
	}

	return credentials
}

// RevokeCredential revokes a credential
func (k Keeper) RevokeCredential(ctx sdk.Context, id string, issuerDid string) error {
	credential, found := k.GetCredential(ctx, id)
	if !found {
		return types.ErrCredentialNotFound
	}

	if credential.IssuerDid != issuerDid {
		return types.ErrUnauthorizedCredentialType
	}

	credential.Revoked = true
	credential.IssuedAt = ctx.BlockTime()

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(CredentialPrefix))
	credentialBytes := k.cdc.MustMarshal(&credential)
	store.Set([]byte(id), credentialBytes)

	return nil
}

// ValidateCredential validates a credential
func (k Keeper) ValidateCredential(credential types.Credential) error {
	if credential.Id == "" {
		return fmt.Errorf("credential ID cannot be empty")
	}

	if credential.IssuerDid == "" {
		return fmt.Errorf("issuer DID cannot be empty")
	}

	if credential.SubjectDid == "" {
		return fmt.Errorf("subject DID cannot be empty")
	}

	// Check if issuer exists
	if _, found := k.GetDIDDocument(sdk.Context{}, credential.IssuerDid); !found {
		return fmt.Errorf("issuer DID %s not found", credential.IssuerDid)
	}

	// Check if subject exists
	if _, found := k.GetDIDDocument(sdk.Context{}, credential.SubjectDid); !found {
		return fmt.Errorf("subject DID %s not found", credential.SubjectDid)
	}

	// Validate expiry
	if credential.ExpiresAt != nil && credential.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("credential has expired")
	}

	return nil
}

// DeleteCredential deletes a credential
func (k Keeper) DeleteCredential(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(CredentialPrefix))
	store.Delete([]byte(id))

	// Also clean up the subject index
	if credential, found := k.GetCredential(ctx, id); found {
		subjectStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("credential_by_subject/"))
		subjectKey := []byte(fmt.Sprintf("%s/%s", credential.SubjectDid, id))
		subjectStore.Delete(subjectKey)
	}
}

// GetAllCredentials returns all credentials
func (k Keeper) GetAllCredentials(ctx sdk.Context) []types.Credential {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(CredentialPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var credentials []types.Credential
	for ; iterator.Valid(); iterator.Next() {
		var cred types.Credential
		k.cdc.MustUnmarshal(iterator.Value(), &cred)
		credentials = append(credentials, cred)
	}
	return credentials
}
