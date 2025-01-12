package types

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Security-related constants
const (
	SecurityEventPrefix = "security_event/"
	RateLimitPrefix     = "rate_limit/"
)

var (
	ErrInvalidAuditLogType = sdkerrors.Register(ModuleName, 1401, "invalid audit log type")
	ErrInvalidAuditAction  = sdkerrors.Register(ModuleName, 1402, "invalid audit action")
	ErrInvalidAuditActor   = sdkerrors.Register(ModuleName, 1403, "invalid audit actor")
	ErrInvalidAuditDID     = sdkerrors.Register(ModuleName, 1404, "invalid audit DID")
)

// ValidateBasic performs basic validation of the audit log entry
func (a *AuditLogEntry) ValidateBasic() error {
	if a.Id == "" {
		return sdkerrors.Register(ModuleName, 1401, "audit log ID cannot be empty")
	}

	if a.Did == "" {
		return sdkerrors.Register(ModuleName, 1402, "audit log DID cannot be empty")
	}

	if a.Action == "" {
		return sdkerrors.Register(ModuleName, 1403, "audit log action cannot be empty")
	}

	if a.Actor == "" {
		return sdkerrors.Register(ModuleName, 1404, "audit log actor cannot be empty")
	}

	if a.Timestamp == 0 {
		return sdkerrors.Register(ModuleName, 1405, "audit log timestamp cannot be zero")
	}

	return nil
}

// NewAuditLogEntry creates a new audit log entry
func NewAuditLogEntry(
	id string,
	did string,
	action string,
	actor string,
	logType AuditLogType,
	severity SecuritySeverity,
	metadata map[string]string,
) *AuditLogEntry {
	return &AuditLogEntry{
		Id:        id,
		Type:      logType,
		Did:       did,
		Action:    action,
		Actor:     actor,
		Timestamp: time.Now().Unix(),
		Severity:  severity,
		Metadata:  metadata,
	}
}
