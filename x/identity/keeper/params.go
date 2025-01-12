package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This file contains parameter-related keeper methods

// IsOAuthProviderAllowed checks if an OAuth provider is allowed
func (k Keeper) IsOAuthProviderAllowed(ctx sdk.Context, provider string) bool {
	params := k.GetParams(ctx)
	for _, p := range params.CredentialParams.AllowedTypes {
		if p == provider {
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

// IsDIDMethodAllowed checks if a DID method is allowed
func (k Keeper) IsDIDMethodAllowed(ctx sdk.Context, method string) bool {
	params := k.GetParams(ctx)
	for _, m := range params.DidParams.AllowedMethods {
		if m == method {
			return true
		}
	}
	return false
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
