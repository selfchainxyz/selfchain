package types

const (
	// ModuleName defines the module name
	ModuleName = "identity"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_identity"
)

var (
	// DIDDocumentKeyPrefix is the prefix for storing DID documents
	DIDDocumentKeyPrefix = []byte("did_document/")

	// CredentialKeyPrefix is the prefix for storing verifiable credentials
	CredentialKeyPrefix = []byte("credential/")

	// CredentialBySubjectKeyPrefix is the prefix for indexing credentials by subject DID
	CredentialBySubjectKeyPrefix = []byte("credential_by_subject/")

	// VerificationKeyPrefix is the prefix for storing identity verification records
	VerificationKeyPrefix = []byte("verification/")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
