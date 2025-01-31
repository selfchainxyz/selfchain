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
	DIDDocumentPrefix        = "did_document/"
	DIDPrefix                = DIDDocumentPrefix
	CredentialPrefix         = "credential/"
	SocialIdentityPrefix     = "social_identity/"
	SocialIdentityByIDPrefix = "social_identity_by_id/"
	MFAConfigPrefix          = "mfa_config/"
	AuditLogPrefix           = "audit_log/"
	CredentialByDIDPrefix    = "credential_by_did/"
	MFAMethodPrefix          = "mfa_method/"
	MFAChallengePrefix       = "mfa_challenge:"
	OAuthProviderPrefix      = "oauth_provider/"
	OAuthSessionPrefix       = "oauth_session/"
	MFASessionPrefix         = "mfa_session/"
	RecoverySessionPrefix    = "recovery_session/"
	RecoveryPrefix           = "recovery/"
	KeySharePrefix           = "key_share/"
	RateLimitPrefix          = "rate_limit/"
	AuditEventPrefix         = "audit_event/"
	MFAChallengeExpiry       = 5 * time.Minute // MFA challenge expiry time
)

var (
	// DIDDocumentKey is the key for storing DID documents
	DIDDocumentKey = []byte("DIDDocument")
)

// KeyPrefix returns the prefix for a given key type
func KeyPrefix(prefix string) []byte {
	return []byte(prefix)
}

// GetDIDKey returns the store key to retrieve a DID document
func GetDIDKey(did string) []byte {
	return append(KeyPrefix(DIDDocumentPrefix), []byte(did)...)
}

// GetCredentialKey returns the store key to retrieve a Credential
func GetCredentialKey(id string) []byte {
	return append(KeyPrefix(CredentialPrefix), []byte(id)...)
}

// GetCredentialByDIDKey returns the store key to retrieve Credentials by DID
func GetCredentialByDIDKey(did, credentialID string) []byte {
	prefix := append(KeyPrefix(CredentialByDIDPrefix), []byte(did)...)
	return append(prefix, []byte(credentialID)...)
}

// GetCredentialByDIDPrefixKey returns the store prefix key to retrieve all Credentials for a DID
func GetCredentialByDIDPrefixKey(did string) []byte {
	return append(KeyPrefix(CredentialByDIDPrefix), []byte(did)...)
}

// GetSocialIdentityKey returns the store key to retrieve a social identity
func GetSocialIdentityKey(did string) []byte {
	return append(KeyPrefix(SocialIdentityPrefix), []byte(did)...)
}

// GetSocialIdentityBySocialIDKey returns the store key to retrieve a social identity by social ID
func GetSocialIdentityBySocialIDKey(provider string, socialID string) []byte {
	prefix := append(KeyPrefix(SocialIdentityByIDPrefix), []byte(provider)...)
	return append(prefix, []byte(socialID)...)
}

// GetMFAConfigKey returns the store key to retrieve MFA configuration
func GetMFAConfigKey(did string) []byte {
	return append(KeyPrefix(MFAConfigPrefix), []byte(did)...)
}

// GetAuditLogKey returns the store key for audit logs
func GetAuditLogKey(id string) []byte {
	return append(KeyPrefix(AuditLogPrefix), []byte(id)...)
}

// GetRecoverySessionKey returns the store key for recovery sessions
func GetRecoverySessionKey(id string) []byte {
	return append(KeyPrefix(RecoverySessionPrefix), []byte(id)...)
}

// GetRecoveryKey returns the store key for recovery data
func GetRecoveryKey(id string) []byte {
	return append(KeyPrefix(RecoveryPrefix), []byte(id)...)
}
