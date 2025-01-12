package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

const (
	AuditLogPrefix = "audit_log:"
)

// StoreAuditLog stores an audit log entry
func (k Keeper) StoreAuditLog(ctx sdk.Context, log types.AuditLogEntry) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(AuditLogPrefix))
	key := []byte(log.Id)
	bz := k.cdc.MustMarshal(&log)
	store.Set(key, bz)
	return nil
}

// GetAuditLog returns an audit log entry by ID
func (k Keeper) GetAuditLog(ctx sdk.Context, id string) (*types.AuditLogEntry, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(AuditLogPrefix))
	bz := store.Get([]byte(id))
	if bz == nil {
		return nil, false
	}

	var log types.AuditLogEntry
	k.cdc.MustUnmarshal(bz, &log)
	return &log, true
}

// GetAllAuditLogs returns all audit logs
func (k Keeper) GetAllAuditLogs(ctx sdk.Context) []types.AuditLogEntry {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(AuditLogPrefix))
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(AuditLogPrefix))
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

// GetAuditLogsByCredential returns all audit logs for a credential
func (k Keeper) GetAuditLogsByCredential(ctx sdk.Context, credentialID string) []types.AuditLogEntry {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(AuditLogPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var logs []types.AuditLogEntry
	for ; iterator.Valid(); iterator.Next() {
		var log types.AuditLogEntry
		k.cdc.MustUnmarshal(iterator.Value(), &log)
		if log.Metadata["credential_id"] == credentialID {
			logs = append(logs, log)
		}
	}
	return logs
}
