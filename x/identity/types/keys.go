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

	// Key prefixes for different types of data
	DIDPrefix              = "did/"
	CredentialPrefix       = "credential/"
	SocialIdentityPrefix   = "social_identity/"
	SocialIdentityByIDPrefix = "social_identity_by_id/"
	MFAConfigPrefix        = "mfa_config/"
	AuditLogPrefix         = "audit_log/"
	CredentialByDIDPrefix  = "credential_by_did/"
	MFAChallengePrefix     = "mfa_challenge:"
	MFAChallengeExpiry     = 5 * time.Minute // MFA challenge expiry time
)

// DIDKey returns the store key to retrieve a DID document
func DIDKey(did string) []byte {
	return []byte(DIDPrefix + did)
}

// CredentialKey returns the store key to retrieve a Credential
func CredentialKey(id string) []byte {
	return []byte(CredentialPrefix + id)
}

// CredentialByDIDKey returns the store key to retrieve Credentials by DID
func CredentialByDIDKey(did, credentialID string) []byte {
	return []byte(CredentialByDIDPrefix + did + "/" + credentialID)
}

// CredentialByDIDPrefixKey returns the store prefix key to retrieve all Credentials for a DID
func CredentialByDIDPrefixKey(did string) []byte {
	return []byte(CredentialByDIDPrefix + did + "/")
}

// SocialIdentityKey returns the store key to retrieve a social identity
func SocialIdentityKey(did string) []byte {
	return []byte(SocialIdentityPrefix + did)
}

// MFAConfigKey returns the store key to retrieve MFA configuration
func MFAConfigKey(did string) []byte {
	return []byte(MFAConfigPrefix + did)
}

// GetAuditLogKey returns the store key for audit logs
func GetAuditLogKey(id string) []byte {
	return []byte(AuditLogPrefix + id)
}
