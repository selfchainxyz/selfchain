package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

const (
	// RateLimitPrefix is the prefix for storing rate limit data
	RateLimitPrefix = "rate_limit/"

	// Default rate limits
	DefaultMaxRequestsPerMinute = 30
	DefaultMaxRequestsPerHour   = 100
	DefaultMaxRequestsPerDay    = 1000
)

// RateLimitData stores rate limiting information
type RateLimitData struct {
	LastRequestTime      int64
	RequestsInLastMinute int32
	RequestsInLastHour   int32
	RequestsInLastDay    int32
}

// CheckRateLimit checks if a request should be rate limited
func (k Keeper) CheckRateLimit(ctx sdk.Context, identifier string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(RateLimitPrefix))
	key := []byte(identifier)

	var data types.RateLimitData
	dataBytes := store.Get(key)
	if dataBytes != nil {
		k.cdc.MustUnmarshal(dataBytes, &data)
	}

	currentTime := ctx.BlockTime().Unix()

	// Reset counters if time windows have passed
	if currentTime-data.LastRequestTime >= 86400 { // 24 hours
		data.RequestsInLastDay = 0
		data.RequestsInLastHour = 0
		data.RequestsInLastMinute = 0
	} else if currentTime-data.LastRequestTime >= 3600 { // 1 hour
		data.RequestsInLastHour = 0
		data.RequestsInLastMinute = 0
	} else if currentTime-data.LastRequestTime >= 60 { // 1 minute
		data.RequestsInLastMinute = 0
	}

	// Check rate limits
	if data.RequestsInLastMinute >= DefaultMaxRequestsPerMinute {
		return fmt.Errorf("rate limit exceeded: too many requests per minute")
	}
	if data.RequestsInLastHour >= DefaultMaxRequestsPerHour {
		return fmt.Errorf("rate limit exceeded: too many requests per hour")
	}
	if data.RequestsInLastDay >= DefaultMaxRequestsPerDay {
		return fmt.Errorf("rate limit exceeded: too many requests per day")
	}

	// Update counters
	data.RequestsInLastMinute++
	data.RequestsInLastHour++
	data.RequestsInLastDay++
	data.LastRequestTime = currentTime

	// Store updated data
	store.Set(key, k.cdc.MustMarshal(&data))

	return nil
}

// ResetRateLimit resets rate limit data for an identifier
func (k Keeper) ResetRateLimit(ctx sdk.Context, identifier string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(RateLimitPrefix))
	store.Delete([]byte(identifier))
}

// GetRateLimitData gets rate limit data for an identifier
func (k Keeper) GetRateLimitData(ctx sdk.Context, identifier string) types.RateLimitData {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(RateLimitPrefix))
	dataBytes := store.Get([]byte(identifier))

	var data types.RateLimitData
	if dataBytes != nil {
		k.cdc.MustUnmarshal(dataBytes, &data)
	}
	return data
}
