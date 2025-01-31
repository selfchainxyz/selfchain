package types

// Event types for the identity module
const (
	EventTypeAuditLog = "audit_log"
	EventTypeMFAVerification = "mfa_verification"
	EventTypeRateLimit = "rate_limit"
	EventTypeOAuthVerification = "oauth_verification"
	EventTypeOAuthSuccess      = "oauth_success"
	EventTypeOAuthFailure      = "oauth_failure"
)

// Event attribute keys
const (
	AttributeKeyDID       = "did"
	AttributeKeyEventType = "event_type"
	AttributeKeySuccess   = "success"
	AttributeKeyDetails   = "details"
	AttributeKeyTimestamp = "timestamp"
	AttributeKeyMethodID = "method_id"
	AttributeKeyError    = "error"
	AttributeKeyOperation = "operation"
	AttributeKeyProvider  = "provider"
	AttributeKeySocialID  = "social_id"
	AttributeKeyMethod    = "method"
	AttributeKeyCode      = "code"
	AttributeKeyType      = "type"
	AttributeKeyStatus    = "status"
)

// Rate limit key prefixes
const (
	RateLimitCountKey = "rate_limit_count"
)
