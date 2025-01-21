// Package types contains the types used in the identity module
package types

// Constants for identity module
const (
	// Version defines the current version the identity module
	Version = 1
)

// Store prefixes
const (
	// SocialIdentityKeyPrefix is the key prefix for social identities
	SocialIdentityKeyPrefix = "social_identity"

	// AuditLogKey is the key prefix for audit logs
	AuditLogKey = "audit_log"
)

// Rate limiting constants
const (
	// DefaultRateLimit defines the default rate limit for API calls
	DefaultRateLimit = 100

	// DefaultBurst defines the default burst size for rate limiting
	DefaultBurst = 10
)
