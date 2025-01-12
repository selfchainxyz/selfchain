package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

const (
	DIDDocumentPrefix = "did_document/"
)

// SetDIDDocument stores a DID document
func (k Keeper) SetDIDDocument(ctx sdk.Context, did string, document types.DIDDocument) error {
	if err := k.ValidateDIDDocument(document); err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(DIDDocumentPrefix))
	documentBytes := k.cdc.MustMarshal(&document)
	store.Set([]byte(did), documentBytes)
	return nil
}

// GetDIDDocument returns a DID document by ID
func (k Keeper) GetDIDDocument(ctx sdk.Context, did string) (types.DIDDocument, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(DIDDocumentPrefix))
	documentBytes := store.Get([]byte(did))
	if documentBytes == nil {
		return types.DIDDocument{}, false
	}

	var document types.DIDDocument
	k.cdc.MustUnmarshal(documentBytes, &document)
	return document, true
}

// GetAllDIDDocuments returns all DID documents
func (k Keeper) GetAllDIDDocuments(ctx sdk.Context) []types.DIDDocument {
	var documents []types.DIDDocument
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(DIDDocumentPrefix))

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var document types.DIDDocument
		k.cdc.MustUnmarshal(iterator.Value(), &document)
		documents = append(documents, document)
	}

	return documents
}

// DeleteDIDDocument deletes a DID document
func (k Keeper) DeleteDIDDocument(ctx sdk.Context, did string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(DIDDocumentPrefix))
	store.Delete([]byte(did))
}

// ValidateDIDDocument validates a DID document
func (k Keeper) ValidateDIDDocument(document types.DIDDocument) error {
	if document.Id == "" {
		return fmt.Errorf("DID document ID cannot be empty")
	}

	if len(document.VerificationMethod) == 0 {
		return fmt.Errorf("DID document must have at least one verification method")
	}

	for _, method := range document.VerificationMethod {
		if method.Id == "" {
			return fmt.Errorf("verification method ID cannot be empty")
		}
		if method.Type == "" {
			return fmt.Errorf("verification method type cannot be empty")
		}
		if method.Controller == "" {
			return fmt.Errorf("verification method controller cannot be empty")
		}
		if method.PublicKeyBase58 == "" {
			return fmt.Errorf("verification method public key cannot be empty")
		}
	}

	return nil
}

// HasDIDDocument checks if a DID document exists
func (k Keeper) HasDIDDocument(ctx sdk.Context, did string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(DIDDocumentPrefix))
	return store.Has([]byte(did))
}
