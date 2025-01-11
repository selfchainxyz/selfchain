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

	// Version defines the current version the IBC module supports
	Version = "identity-1"

	// DIDDocumentPrefix is the prefix for storing DID documents
	DIDDocumentPrefix = "did_document"

	// CredentialPrefix is the prefix for storing credentials
	CredentialPrefix = "credential"

	// VerificationPrefix is the prefix for storing verification statuses
	VerificationPrefix = "verification"

	// SocialIdentityPrefix is the prefix for storing social identities
	SocialIdentityPrefix = "social_identity"

	// DIDSocialPrefix is the prefix for storing DID to social identity mappings
	DIDSocialPrefix = "did_social"

	// PortID is the default port id that module binds to
	PortID = "identity"

	// Key prefixes
	DIDDocumentKey     = "did_document/"
	CredentialKey      = "credential/"
	CredentialSchemaKey = "credential_schema/"
	SocialIdentityKey  = "social_identity/"
)

var (
	// DIDDocumentKeyPrefix is the prefix for storing DID documents
	DIDDocumentKeyPrefix = []byte(DIDDocumentPrefix + "/")

	// CredentialKeyPrefix is the prefix for storing verifiable credentials
	CredentialKeyPrefix = []byte(CredentialPrefix + "/")

	// CredentialBySubjectKeyPrefix is the prefix for indexing credentials by subject DID
	CredentialBySubjectKeyPrefix = []byte("credential_by_subject/")

	// VerificationKeyPrefix is the prefix for storing identity verification records
	VerificationKeyPrefix = []byte(VerificationPrefix + "/")

	// SocialIdentityKeyPrefix is the prefix for storing social identities
	SocialIdentityKeyPrefix = []byte(SocialIdentityPrefix + "/")

	// DIDSocialKeyPrefix is the prefix for storing DID to social identity mappings
	DIDSocialKeyPrefix = []byte(DIDSocialPrefix + "/")

	// PortKey defines the key to store the port ID in store
	PortKey = []byte(ModuleName + "_" + "port")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
