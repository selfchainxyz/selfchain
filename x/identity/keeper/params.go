package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// IsOAuthProviderAllowed checks if an OAuth provider is allowed
func (k Keeper) IsOAuthProviderAllowed(ctx sdk.Context, provider string) bool {
	params := k.GetParams(ctx)
	for _, p := range params.AllowedOauthProviders {
		if p == provider {
			return true
		}
	}
	return false
}

// IsCredentialTypeAllowed checks if a credential type is allowed
func (k Keeper) IsCredentialTypeAllowed(ctx sdk.Context, credentialType string) bool {
	params := k.GetParams(ctx)
	for _, t := range params.AllowedCredentialTypes {
		if t == credentialType {
			return true
		}
	}
	return false
}

// GetVerificationTimeout returns the verification timeout duration
func (k Keeper) GetVerificationTimeout(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).VerificationTimeoutHours
}

// GetMaxCredentialsPerDID returns the maximum number of credentials allowed per DID
func (k Keeper) GetMaxCredentialsPerDID(ctx sdk.Context) uint32 {
	return k.GetParams(ctx).MaxCredentialsPerDid
}

// ValidateParams validates the module parameters
func ValidateParams(params types.Params) error {
	if params.VerificationTimeoutHours == 0 {
		return types.ErrInvalidVerificationTimeout
	}

	if params.MaxCredentialsPerDid == 0 {
		return types.ErrInvalidMaxCredentials
	}

	if len(params.AllowedOauthProviders) == 0 {
		return types.ErrNoOAuthProvidersAllowed
	}

	if len(params.AllowedCredentialTypes) == 0 {
		return types.ErrNoCredentialTypesAllowed
	}

	return nil
}
