package keeper

import (
	"fmt"
	"time"

	"selfchain/x/identity/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	AuditLogPrefix = "audit_log:"
)

// StoreAuditLog stores an audit log entry
func (k Keeper) StoreAuditLog(ctx sdk.Context, log types.AuditLogEntry) error {
	if err := log.ValidateBasic(); err != nil {
		return sdkerrors.Wrapf(err, "invalid audit log")
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AuditLogPrefix))
	key := []byte(log.Id)

	bz, err := k.cdc.Marshal(&log)
	if err != nil {
		return fmt.Errorf("failed to marshal audit log: %w", err)
	}

	store.Set(key, bz)
	return nil
}

// GetAuditLog retrieves an audit log entry by ID
func (k Keeper) GetAuditLog(ctx sdk.Context, id string) (*types.AuditLogEntry, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AuditLogPrefix))
	key := []byte(id)

	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("audit log not found: %s", id)
	}

	var log types.AuditLogEntry
	if err := k.cdc.Unmarshal(bz, &log); err != nil {
		return nil, fmt.Errorf("failed to unmarshal audit log: %w", err)
	}

	return &log, nil
}

// GetAllAuditLogs returns all audit logs
func (k Keeper) GetAllAuditLogs(ctx sdk.Context) []types.AuditLogEntry {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AuditLogPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var logs []types.AuditLogEntry
	for ; iterator.Valid(); iterator.Next() {
		var log types.AuditLogEntry
		k.cdc.MustUnmarshal(iterator.Value(), &log)
		logs = append(logs, log)
	}
	return logs
}

// GetAuditLogsByDID returns all audit logs for a DID
func (k Keeper) GetAuditLogsByDID(ctx sdk.Context, did string) []types.AuditLogEntry {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AuditLogPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var logs []types.AuditLogEntry
	for ; iterator.Valid(); iterator.Next() {
		var log types.AuditLogEntry
		k.cdc.MustUnmarshal(iterator.Value(), &log)
		if log.Did == did {
			logs = append(logs, log)
		}
	}
	return logs
}

// CreateAuditLog creates an audit log entry
func (k Keeper) CreateAuditLog(ctx sdk.Context, entry types.AuditLogEntry) error {
	if err := entry.ValidateBasic(); err != nil {
		return err
	}

	// Set timestamp if not set
	if entry.Timestamp == 0 {
		entry.Timestamp = time.Now().Unix()
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AuditLogPrefix))
	key := []byte(entry.Id)

	bz, err := k.cdc.Marshal(&entry)
	if err != nil {
		return fmt.Errorf("failed to marshal audit entry: %w", err)
	}

	store.Set(key, bz)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"audit_log",
			sdk.NewAttribute("id", entry.Id),
			sdk.NewAttribute("type", entry.Type.String()),
			sdk.NewAttribute("did", entry.Did),
			sdk.NewAttribute("action", entry.Action),
			sdk.NewAttribute("severity", entry.Severity.String()),
		),
	)

	return nil
}

// GetAuditLogs returns audit logs filtered by DID and type
func (k Keeper) GetAuditLogs(ctx sdk.Context, did string, logType types.AuditLogType) ([]types.AuditLogEntry, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AuditLogPrefix))
	var logs []types.AuditLogEntry

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var entry types.AuditLogEntry
		if err := k.cdc.Unmarshal(iterator.Value(), &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal audit entry: %w", err)
		}

		if (did == "" || entry.Did == did) && (logType == types.AuditLogType_AUDIT_LOG_TYPE_UNSPECIFIED || entry.Type == logType) {
			logs = append(logs, entry)
		}
	}

	return logs, nil
}

// GetAuditLogsByTimeRange returns audit logs within a time range
func (k Keeper) GetAuditLogsByTimeRange(ctx sdk.Context, did string, startTime, endTime time.Time) ([]types.AuditLogEntry, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AuditLogPrefix))
	var logs []types.AuditLogEntry

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var entry types.AuditLogEntry
		if err := k.cdc.Unmarshal(iterator.Value(), &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal audit entry: %w", err)
		}

		if (did == "" || entry.Did == did) &&
			(entry.Timestamp >= startTime.Unix() && entry.Timestamp <= endTime.Unix()) {
			logs = append(logs, entry)
		}
	}

	return logs, nil
}

// GetSecurityEvents returns security-related audit logs
func (k Keeper) GetSecurityEvents(ctx sdk.Context, did string, minSeverity types.SecuritySeverity) ([]types.AuditLogEntry, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AuditLogPrefix))
	var logs []types.AuditLogEntry

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var entry types.AuditLogEntry
		if err := k.cdc.Unmarshal(iterator.Value(), &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal audit entry: %w", err)
		}

		// Filter by severity level
		if (did == "" || entry.Did == did) &&
			(minSeverity == types.SecuritySeverity_SEVERITY_UNSPECIFIED || entry.Severity >= minSeverity) {
			logs = append(logs, entry)
		}
	}

	return logs, nil
}

// ListAuditLogs returns all audit logs for a DID
func (k Keeper) ListAuditLogs(ctx sdk.Context, did string) ([]*types.AuditLogEntry, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AuditLogPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte(did))
	defer iterator.Close()

	var logs []*types.AuditLogEntry
	for ; iterator.Valid(); iterator.Next() {
		var log types.AuditLogEntry
		if err := k.cdc.Unmarshal(iterator.Value(), &log); err != nil {
			continue
		}
		logs = append(logs, &log)
	}

	return logs, nil
}
