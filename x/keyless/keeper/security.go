package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

// SecurityManager handles security-related operations
type SecurityManager struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewSecurityManager creates a new security manager instance
func NewSecurityManager(cdc codec.BinaryCodec, storeKey storetypes.StoreKey) *SecurityManager {
	return &SecurityManager{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

var (
	// Store prefixes
	RateLimitPrefix = []byte("rate_limit/")
	AuditLogPrefix  = []byte("audit_log/")
	AccessPrefix    = []byte("access/")
)

// RateLimitConfig defines rate limiting parameters
type RateLimitConfig struct {
	MaxRequests uint32
	WindowSize  time.Duration
}

// Default rate limit configurations
var (
	DefaultCreateWalletLimit = RateLimitConfig{
		MaxRequests: 5,
		WindowSize:  24 * time.Hour,
	}
	DefaultSigningLimit = RateLimitConfig{
		MaxRequests: 100,
		WindowSize:  time.Hour,
	}
	DefaultRecoveryLimit = RateLimitConfig{
		MaxRequests: 3,
		WindowSize:  24 * time.Hour,
	}
)

// CheckRateLimit verifies if an operation is within rate limits
func (sm *SecurityManager) CheckRateLimit(ctx sdk.Context, operation string, user string) error {
	config := sm.getRateLimitConfig(operation)
	key := sm.getRateLimitKey(operation, user)
	store := ctx.KVStore(sm.storeKey)

	// Get current count and timestamp
	var data types.RateLimitData
	bz := store.Get(key)
	if bz != nil {
		if err := sm.cdc.Unmarshal(bz, &data); err != nil {
			return fmt.Errorf("failed to unmarshal rate limit data: %w", err)
		}
	}

	// Check if window has reset
	now := ctx.BlockTime()
	if data.LastReset != nil && now.Sub(*data.LastReset) > config.WindowSize {
		data.Count = 0
		data.LastReset = &now
	}

	// Check limit
	if data.Count >= config.MaxRequests {
		return fmt.Errorf("rate limit exceeded for operation %s", operation)
	}

	// Update counter
	data.Count++
	if data.LastReset == nil {
		data.LastReset = &now
	}

	bz, err := sm.cdc.Marshal(&data)
	if err != nil {
		return fmt.Errorf("failed to marshal rate limit data: %w", err)
	}
	store.Set(key, bz)

	return nil
}

// ValidateOwnership checks if a user has permission to access a wallet
func (sm *SecurityManager) ValidateOwnership(ctx sdk.Context, wallet string, user string) error {
	store := ctx.KVStore(sm.storeKey)
	key := append(AccessPrefix, []byte(fmt.Sprintf("%s/%s", wallet, user))...)

	if !store.Has(key) {
		return fmt.Errorf("user %s is not authorized for wallet %s", user, wallet)
	}

	return nil
}

// LogOperation records an audit log entry
func (sm *SecurityManager) LogOperation(ctx sdk.Context, operation string, wallet string, user string) error {
	store := ctx.KVStore(sm.storeKey)
	
	now := ctx.BlockTime()
	logEntry := types.AuditLog{
		Operation:    operation,
		WalletId:    wallet,
		User:        user,
		Timestamp:   &now,
		BlockHeight: ctx.BlockHeight(),
		Success:     true,
	}

	bz, err := sm.cdc.Marshal(&logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal audit log: %w", err)
	}

	// Use timestamp as part of the key for chronological ordering
	key := append(AuditLogPrefix, []byte(fmt.Sprintf("%d/%s/%s", ctx.BlockHeight(), wallet, operation))...)
	store.Set(key, bz)

	return nil
}

// GrantAccess grants a user access to a wallet
func (sm *SecurityManager) GrantAccess(ctx sdk.Context, wallet string, user string) error {
	store := ctx.KVStore(sm.storeKey)
	key := append(AccessPrefix, []byte(fmt.Sprintf("%s/%s", wallet, user))...)
	
	// Store access grant with timestamp
	now := ctx.BlockTime()
	data := types.AccessData{
		GrantedAt: &now,
	}
	bz, err := sm.cdc.Marshal(&data)
	if err != nil {
		return fmt.Errorf("failed to marshal access data: %w", err)
	}
	
	store.Set(key, bz)
	return nil
}

// RevokeAccess revokes a user's access to a wallet
func (sm *SecurityManager) RevokeAccess(ctx sdk.Context, wallet string, user string) error {
	store := ctx.KVStore(sm.storeKey)
	key := append(AccessPrefix, []byte(fmt.Sprintf("%s/%s", wallet, user))...)
	
	store.Delete(key)
	return nil
}

// Helper functions

func (sm *SecurityManager) getRateLimitConfig(operation string) RateLimitConfig {
	switch operation {
	case "create_wallet":
		return DefaultCreateWalletLimit
	case "sign_transaction":
		return DefaultSigningLimit
	case "recover_wallet":
		return DefaultRecoveryLimit
	default:
		return DefaultSigningLimit
	}
}

func (sm *SecurityManager) getRateLimitKey(operation string, user string) []byte {
	return append(RateLimitPrefix, []byte(fmt.Sprintf("%s/%s", operation, user))...)
}

// GetAuditLogs retrieves audit logs for a wallet
func (sm *SecurityManager) GetAuditLogs(ctx sdk.Context, wallet string, limit uint32) ([]*types.AuditLog, error) {
	store := ctx.KVStore(sm.storeKey)
	prefix := append(AuditLogPrefix, []byte(wallet)...)
	
	var logs []*types.AuditLog
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	count := uint32(0)
	for ; iterator.Valid() && count < limit; iterator.Next() {
		var log types.AuditLog
		if err := sm.cdc.Unmarshal(iterator.Value(), &log); err != nil {
			return nil, fmt.Errorf("failed to unmarshal audit log: %w", err)
		}
		logs = append(logs, &log)
		count++
	}

	return logs, nil
}
