package types

import (
	"time"
)

const (
	// ModuleName defines the module name
	ModuleName = "identity"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_identity"

	// ParamsKey is the store key for module parameters
	ParamsKey = "params"

	// Store prefixes
	DIDDocumentPrefix     = "did_document/"
	DIDPrefix             = DIDDocumentPrefix
	CredentialPrefix      = "credential/"
	SocialIdentityPrefix  = "social_identity/"
	SocialIdentityByIDPrefix = "social_identity_by_id/"
	MFAConfigPrefix       = "mfa_config/"
	AuditLogPrefix        = "audit_log/"
	CredentialByDIDPrefix = "credential_by_did/"
	MFAMethodPrefix       = "mfa_method/"
	MFAChallengePrefix    = "mfa_challenge:"
	OAuthProviderPrefix   = "oauth_provider/"
	OAuthSessionPrefix    = "oauth_session/"
	MFASessionPrefix      = "mfa_session/"
	RecoverySessionPrefix = "recovery_session/"
	KeySharePrefix        = "key_share/"
	RateLimitPrefix       = "rate_limit/"
	AuditEventPrefix      = "audit_event/"
	MFAChallengeExpiry    = 5 * time.Minute // MFA challenge expiry time
)

// KeyPrefix returns the prefix for a given key type
func KeyPrefix(prefix string) []byte {
	return []byte(prefix)
}

// GetDIDKey returns the store key to retrieve a DID document
func GetDIDKey(did string) []byte {
	return []byte(DIDPrefix + did)
}

// GetCredentialKey returns the store key to retrieve a Credential
func GetCredentialKey(id string) []byte {
	return []byte(CredentialPrefix + id)
}

// GetCredentialByDIDKey returns the store key to retrieve Credentials by DID
func GetCredentialByDIDKey(did, credentialID string) []byte {
	return []byte(CredentialByDIDPrefix + did + "/" + credentialID)
}

// GetCredentialByDIDPrefixKey returns the store prefix key to retrieve all Credentials for a DID
func GetCredentialByDIDPrefixKey(did string) []byte {
	return []byte(CredentialByDIDPrefix + did + "/")
}

// GetSocialIdentityKey returns the store key to retrieve a social identity
func GetSocialIdentityKey(did string) []byte {
	return []byte(SocialIdentityPrefix + did)
}

// GetMFAConfigKey returns the store key to retrieve MFA configuration
func GetMFAConfigKey(did string) []byte {
	return []byte(MFAConfigPrefix + did)
}

// GetAuditLogKey returns the store key for audit logs
func GetAuditLogKey(id string) []byte {
	return []byte(AuditLogPrefix + id)
}
