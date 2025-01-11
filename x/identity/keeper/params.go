package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

// GetVerificationStatus returns the verification status for a DID
func (k Keeper) GetVerificationStatus(ctx sdk.Context, did string) (types.VerificationStatus, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(VerificationPrefix))
	value := store.Get([]byte(did))
	if value == nil {
		return types.VerificationStatus{}, false
	}

	var status types.VerificationStatus
	k.cdc.MustUnmarshal(value, &status)
	return status, true
}

// SetVerificationStatus sets the verification status for a DID
func (k Keeper) SetVerificationStatus(ctx sdk.Context, did string, status types.VerificationStatus) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(VerificationPrefix))
	value := k.cdc.MustMarshal(&status)
	store.Set([]byte(did), value)
}

// GetAllVerificationStatuses returns all verification statuses
func (k Keeper) GetAllVerificationStatuses(ctx sdk.Context) []types.VerificationStatus {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(VerificationPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var statuses []types.VerificationStatus
	for ; iterator.Valid(); iterator.Next() {
		var status types.VerificationStatus
		k.cdc.MustUnmarshal(iterator.Value(), &status)
		statuses = append(statuses, status)
	}

	return statuses
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
