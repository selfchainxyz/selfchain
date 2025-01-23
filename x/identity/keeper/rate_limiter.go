package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

const (
	// RateLimitPrefix is the prefix for storing rate limit data
	RateLimitPrefix = "rate_limit/"

	// Default rate limits
	DefaultMaxRequests = 100  // Default max requests per time window
	DefaultTimeWindow  = 3600 // Default time window in seconds (1 hour)
	DefaultBurstLimit = 10   // Default burst limit
)

// CheckRateLimit checks if an operation is rate limited
func (k Keeper) CheckRateLimit(ctx sdk.Context, did string, operation string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(RateLimitPrefix))
	key := []byte(fmt.Sprintf("%s/%s", did, operation))

	var rateLimit types.RateLimit
	limitBytes := store.Get(key)
	if limitBytes != nil {
		k.cdc.MustUnmarshal(limitBytes, &rateLimit)
	} else {
		// Initialize new rate limit with defaults
		rateLimit = types.RateLimit{
			Did:          did,
			Operation:    operation,
			MaxRequests:  DefaultMaxRequests,
			TimeWindow:   DefaultTimeWindow,
			BurstLimit:   DefaultBurstLimit,
			CurrentCount: 0,
			LastReset:    ctx.BlockTime(),
		}
	}

	currentTime := ctx.BlockTime()

	// Reset if time window has passed
	if currentTime.Sub(rateLimit.LastReset) >= time.Duration(rateLimit.TimeWindow)*time.Second {
		rateLimit.LastReset = currentTime
		rateLimit.CurrentCount = 0
	}

	// Check if we're within burst limit
	if rateLimit.CurrentCount <= rateLimit.BurstLimit {
		rateLimit.CurrentCount++
		store.Set(key, k.cdc.MustMarshal(&rateLimit))
		return nil
	}

	// Check if we're within max requests for the time window
	if rateLimit.CurrentCount >= rateLimit.MaxRequests {
		return types.ErrRateLimitExceeded
	}

	// Increment counter and update
	rateLimit.CurrentCount++
	store.Set(key, k.cdc.MustMarshal(&rateLimit))
	return nil
}

// ResetRateLimit resets rate limit data for a DID and operation
func (k Keeper) ResetRateLimit(ctx sdk.Context, did string, operation string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(RateLimitPrefix))
	key := []byte(fmt.Sprintf("%s/%s", did, operation))
	store.Delete(key)
}

// GetRateLimit gets rate limit data for a DID and operation
func (k Keeper) GetRateLimit(ctx sdk.Context, did string, operation string) *types.RateLimit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(RateLimitPrefix))
	key := []byte(fmt.Sprintf("%s/%s", did, operation))

	limitBytes := store.Get(key)
	if limitBytes == nil {
		return nil
	}

	var rateLimit types.RateLimit
	k.cdc.MustUnmarshal(limitBytes, &rateLimit)
	return &rateLimit
}

// SetRateLimit sets a custom rate limit for a DID and operation
func (k Keeper) SetRateLimit(ctx sdk.Context, rateLimit types.RateLimit) error {
	if err := rateLimit.ValidateBasic(); err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(RateLimitPrefix))
	key := []byte(fmt.Sprintf("%s/%s", rateLimit.Did, rateLimit.Operation))
	store.Set(key, k.cdc.MustMarshal(&rateLimit))
	return nil
}
