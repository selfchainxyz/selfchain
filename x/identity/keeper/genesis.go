package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	// Set all the DID documents
	for _, elem := range genState.DidDocuments {
		k.SetDIDDocument(ctx, elem.Id, elem)
	}

	// Set all the credentials
	for _, elem := range genState.Credentials {
		if err := k.CreateCredential(ctx, &elem); err != nil {
			panic(fmt.Sprintf("failed to create credential: %v", err))
		}
	}

	// Set all the social identities
	for _, elem := range genState.SocialIdentities {
		k.StoreSocialIdentity(ctx, elem)
	}

	// Set all the MFA configurations
	for _, elem := range genState.MfaConfigs {
		k.StoreMFAConfig(ctx, elem)
	}

	// Set module parameters
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Get all DID documents
	didDocuments := k.GetAllDIDDocuments(ctx)
	genesis.DidDocuments = didDocuments

	// Get all credentials
	credentials, err := k.GetAllCredentials(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to get all credentials: %v", err))
	}
	genesis.Credentials = credentials

	// Get all social identities
	socialIdentities := k.GetAllSocialIdentities(ctx)
	genesis.SocialIdentities = socialIdentities

	// Get all MFA configurations
	mfaConfigs := k.GetAllMFAConfigs(ctx)
	genesis.MfaConfigs = mfaConfigs

	// Get module parameters
	genesis.Params = k.GetParams(ctx)

	return genesis
}
