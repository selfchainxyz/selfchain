package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

// SaveAuditEvent saves an audit event to the store
func (k Keeper) SaveAuditEvent(ctx sdk.Context, event *types.AuditEvent) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AuditEventKey))
	
	// Set timestamp if not already set
	if event.Timestamp == nil {
		now := time.Now().UTC()
		event.Timestamp = &now
	}

	bz := k.cdc.MustMarshal(event)
	key := types.GetAuditEventKey(event.WalletId, event.Timestamp.Unix())
	store.Set(key, bz)

	return nil
}

// GetAuditEventsForWallet returns all audit events for a wallet
func (k Keeper) GetAuditEventsForWallet(ctx sdk.Context, walletId string) []*types.AuditEvent {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AuditEventKey))
	iterator := store.Iterator([]byte(walletId), nil)
	defer iterator.Close()

	var events []*types.AuditEvent
	for ; iterator.Valid(); iterator.Next() {
		var event types.AuditEvent
		k.cdc.MustUnmarshal(iterator.Value(), &event)
		if event.WalletId == walletId {
			events = append(events, &event)
		}
	}

	return events
}

// CreateAuditEvent creates and saves a new audit event
func (k Keeper) CreateAuditEvent(ctx sdk.Context, walletId string, eventType string, success bool, details string, creator string, chainId string) error {
	event := &types.AuditEvent{
		WalletId:  walletId,
		EventType: eventType,
		Success:   success,
		Details:   details,
		Creator:   creator,
		ChainId:   chainId,
		Timestamp: nil, // Will be set in SaveAuditEvent
	}

	return k.SaveAuditEvent(ctx, event)
}
