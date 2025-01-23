package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateBasic performs basic validation of rate limit configuration
func (r *RateLimit) ValidateBasic() error {
	if r.Did == "" {
		return errors.Wrap(ErrInvalidRateLimit, "DID cannot be empty")
	}

	if r.Operation == "" {
		return errors.Wrap(ErrInvalidRateLimit, "operation cannot be empty")
	}

	if r.MaxRequests <= 0 {
		return errors.Wrap(ErrInvalidRateLimit, "max requests must be positive")
	}

	if r.TimeWindow <= 0 {
		return errors.Wrap(ErrInvalidRateLimit, "time window must be positive")
	}

	if r.BurstLimit <= 0 {
		return errors.Wrap(ErrInvalidRateLimit, "burst limit must be positive")
	}

	if r.BurstLimit > r.MaxRequests {
		return errors.Wrap(ErrInvalidRateLimit, "burst limit cannot be greater than max requests")
	}

	return nil
}

// IsExceeded checks if the rate limit has been exceeded
func (r *RateLimit) IsExceeded(currentCount uint64, currentTime time.Time) bool {
	// Check if we're still within the time window
	if currentTime.Sub(r.LastReset) >= time.Duration(r.TimeWindow)*time.Second {
		return false // Not exceeded if we're in a new time window
	}

	// Check if we're within burst limit
	if currentCount <= r.BurstLimit {
		return false
	}

	// Check if we're within max requests for the time window
	return currentCount > r.MaxRequests
}
