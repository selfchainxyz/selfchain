package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	// Set all the DIDs
	for _, did := range genState.Dids {
		k.SetDIDDocument(ctx, did.Id, did)
	}

	// Set all the credentials
	for _, cred := range genState.Credentials {
		k.SetCredential(ctx, cred)
	}

	// Set all the social identities
	for _, identity := range genState.SocialIdentities {
		if err := k.StoreSocialIdentity(ctx, identity); err != nil {
			panic(err)
		}
	}

	// Set all the credential schemas
	for _, schema := range genState.CredentialSchemas {
		k.SetCredentialSchema(ctx, schema)
	}

	// Set module parameters
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Get all DIDs
	genesis.Dids = k.GetAllDIDDocuments(ctx)

	// Get all credentials
	genesis.Credentials = k.GetAllCredentials(ctx)

	// Get all social identities
	genesis.SocialIdentities = k.GetAllSocialIdentities(ctx)

	// Get all credential schemas
	genesis.CredentialSchemas = k.GetAllCredentialSchemas(ctx)

	// Get module parameters
	genesis.Params = k.GetParams(ctx)

	return genesis
}

// GetParams returns the current module parameters
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("params"))
	if bz == nil {
		return params
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the module parameters
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set([]byte("params"), bz)
}

// GetAllSocialIdentities returns all social identities
func (k Keeper) GetAllSocialIdentities(ctx sdk.Context) []types.SocialIdentity {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte("social:"))
	defer iterator.Close()

	var identities []types.SocialIdentity
	for ; iterator.Valid(); iterator.Next() {
		var identity types.SocialIdentity
		k.cdc.MustUnmarshal(iterator.Value(), &identity)
		identities = append(identities, identity)
	}
	return identities
}
