package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This file contains parameter-related keeper methods

// IsOAuthProviderAllowed checks if an OAuth provider is allowed
func (k Keeper) IsOAuthProviderAllowed(ctx sdk.Context, provider string) bool {
	params := k.GetParams(ctx)
	for _, p := range params.DidParams.AllowedMethods {
		if p == provider {
			return true
		}
	}
	return false
}

// GetMaxControllers returns the maximum number of controllers allowed per DID
func (k Keeper) GetMaxControllers(ctx sdk.Context) int64 {
	params := k.GetParams(ctx)
	return params.DidParams.MaxControllers
}

// GetMaxServices returns the maximum number of services allowed per DID
func (k Keeper) GetMaxServices(ctx sdk.Context) int64 {
	params := k.GetParams(ctx)
	return params.DidParams.MaxServices
}

// GetMaxVerificationMethods returns the maximum number of verification methods allowed per DID
func (k Keeper) GetMaxVerificationMethods(ctx sdk.Context) int64 {
	params := k.GetParams(ctx)
	return params.DidParams.MaxVerificationMethods
}

// GetMaxCredentialsPerDID returns the maximum number of credentials allowed per DID
func (k Keeper) GetMaxCredentialsPerDID(ctx sdk.Context) int64 {
	params := k.GetParams(ctx)
	return params.CredentialParams.MaxCredentialsPerDid
}

// GetMaxClaimSize returns the maximum size of a claim in bytes
func (k Keeper) GetMaxClaimSize(ctx sdk.Context) int64 {
	params := k.GetParams(ctx)
	return params.CredentialParams.MaxClaimSize
}

// GetMaxValidityDuration returns the maximum duration a credential can be valid for
func (k Keeper) GetMaxValidityDuration(ctx sdk.Context) int64 {
	params := k.GetParams(ctx)
	return params.CredentialParams.MaxValidityDuration
}

// GetMaxMFAMethods returns the maximum number of MFA methods allowed per DID
func (k Keeper) GetMaxMFAMethods(ctx sdk.Context) int64 {
	params := k.GetParams(ctx)
	return params.MfaParams.MaxMethods
}

// GetMFAChallengeExpiry returns the expiry duration for MFA challenges
func (k Keeper) GetMFAChallengeExpiry(ctx sdk.Context) int64 {
	params := k.GetParams(ctx)
	return params.MfaParams.ChallengeExpiry
}

// GetMaxMFAFailedAttempts returns the maximum number of failed MFA attempts allowed
func (k Keeper) GetMaxMFAFailedAttempts(ctx sdk.Context) int64 {
	params := k.GetParams(ctx)
	return params.MfaParams.MaxFailedAttempts
}

// IsMFAMethodAllowed checks if an MFA method is allowed
func (k Keeper) IsMFAMethodAllowed(ctx sdk.Context, method string) bool {
	params := k.GetParams(ctx)
	for _, m := range params.MfaParams.AllowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// IsCredentialTypeAllowed checks if a credential type is allowed
func (k Keeper) IsCredentialTypeAllowed(ctx sdk.Context, credType string) bool {
	params := k.GetParams(ctx)
	for _, t := range params.CredentialParams.AllowedTypes {
		if t == credType {
			return true
		}
	}
	return false
}
