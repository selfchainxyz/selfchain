package types

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidAuditLogType = sdkerrors.Register(ModuleName, 1401, "invalid audit log type")
	ErrInvalidAuditAction  = sdkerrors.Register(ModuleName, 1402, "invalid audit action")
	ErrInvalidAuditActor   = sdkerrors.Register(ModuleName, 1403, "invalid audit actor")
	ErrInvalidAuditDID     = sdkerrors.Register(ModuleName, 1404, "invalid audit DID")
)

// ValidateBasic performs basic validation of audit log entry
func (a *AuditLogEntry) ValidateBasic() error {
	if a.Id == "" {
		return sdkerrors.Wrap(ErrInvalidAuditLogID, "ID cannot be empty")
	}
	if a.Did == "" {
		return sdkerrors.Wrap(ErrInvalidAuditDID, "DID cannot be empty")
	}
	if a.Action == "" {
		return sdkerrors.Wrap(ErrInvalidAuditAction, "action cannot be empty")
	}
	if a.Actor == "" {
		return sdkerrors.Wrap(ErrInvalidAuditActor, "actor cannot be empty")
	}
	if a.Type == AuditLogType_AUDIT_LOG_TYPE_UNSPECIFIED {
		return sdkerrors.Wrap(ErrInvalidAuditLogType, "audit log type must be specified")
	}
	if a.Timestamp == nil {
		return sdkerrors.Wrap(ErrInvalidAuditLogType, "timestamp must be set")
	}
	return nil
}

// IsCredentialEvent checks if the audit log is related to a credential
func (a *AuditLogEntry) IsCredentialEvent() bool {
	switch a.Type {
	case AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_CREATED,
		AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_UPDATED,
		AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_DELETED,
		AuditLogType_AUDIT_LOG_TYPE_CREDENTIAL_REVOKED:
		return true
	default:
		return false
	}
}

// IsDIDEvent checks if the audit log is related to a DID
func (a *AuditLogEntry) IsDIDEvent() bool {
	switch a.Type {
	case AuditLogType_AUDIT_LOG_TYPE_DID_CREATED,
		AuditLogType_AUDIT_LOG_TYPE_DID_UPDATED,
		AuditLogType_AUDIT_LOG_TYPE_DID_DELETED:
		return true
	default:
		return false
	}
}

// IsMFAEvent checks if the audit log is related to MFA
func (a *AuditLogEntry) IsMFAEvent() bool {
	switch a.Type {
	case AuditLogType_AUDIT_LOG_TYPE_MFA_CONFIGURED,
		AuditLogType_AUDIT_LOG_TYPE_MFA_UPDATED,
		AuditLogType_AUDIT_LOG_TYPE_MFA_DISABLED:
		return true
	default:
		return false
	}
}

// IsSocialIdentityEvent checks if the audit log is related to social identity
func (a *AuditLogEntry) IsSocialIdentityEvent() bool {
	switch a.Type {
	case AuditLogType_AUDIT_LOG_TYPE_SOCIAL_IDENTITY_LINKED,
		AuditLogType_AUDIT_LOG_TYPE_SOCIAL_IDENTITY_UNLINKED:
		return true
	default:
		return false
	}
}

// NewAuditLogEntry creates a new audit log entry
func NewAuditLogEntry(id string, logType AuditLogType, did string, action string, actor string, metadata map[string]string) *AuditLogEntry {
	return &AuditLogEntry{
		Id:        id,
		Type:      logType,
		Did:       did,
		Action:    action,
		Actor:     actor,
		Metadata:  metadata,
		Timestamp: &time.Time{},
	}
}
