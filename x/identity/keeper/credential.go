package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

const (
	CredentialPrefix = "credential/"
)

// SetCredential stores a credential
func (k Keeper) SetCredential(ctx sdk.Context, credential types.Credential) error {
	if err := k.ValidateCredential(credential); err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(CredentialPrefix))
	credentialBytes := k.cdc.MustMarshal(&credential)
	store.Set([]byte(credential.Id), credentialBytes)
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(CredentialPrefix))

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var credential types.Credential
		k.cdc.MustUnmarshal(iterator.Value(), &credential)
		if credential.SubjectDid == subjectDid {
			credentials = append(credentials, credential)
		}
	}

	return credentials
}

// GetAllCredentials returns all credentials
func (k Keeper) GetAllCredentials(ctx sdk.Context) []types.Credential {
	var credentials []types.Credential
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(CredentialPrefix))

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var credential types.Credential
		k.cdc.MustUnmarshal(iterator.Value(), &credential)
		credentials = append(credentials, credential)
	}

	return credentials
}

// ValidateCredential validates a credential
func (k Keeper) ValidateCredential(credential types.Credential) error {
	if credential.Id == "" {
		return types.ErrInvalidCredential
	}

	if credential.SubjectDid == "" {
		return types.ErrInvalidDID
	}

	if credential.IssuerDid == "" {
		return types.ErrInvalidDID
	}

	if credential.Type == "" {
		return types.ErrInvalidCredential
	}

	if credential.IssuedAt.IsZero() {
		return types.ErrInvalidCredential
	}

	return nil
}

// DeleteCredential deletes a credential by ID
func (k Keeper) DeleteCredential(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(CredentialPrefix))
	store.Delete([]byte(id))
}
