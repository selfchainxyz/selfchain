package identity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/keeper"
	"selfchain/x/identity/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *types.GenesisState {
	return &types.GenesisState{
		Params: types.Params{
			AllowedOauthProviders: []string{"google", "github", "apple"},
			VerificationTimeoutHours: 24,
			MaxCredentialsPerDid: 100,
			AllowedCredentialTypes: []string{
				"IdentityCredential",
				"KYCCredential",
				"MembershipCredential",
			},
		},
		Dids: []types.DIDDocument{},
		Credentials: []types.Credential{},
		SocialIdentities: []types.SocialIdentity{},
		CredentialSchemas: []types.CredentialSchema{},
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data *types.GenesisState) error {
	if data.Params.VerificationTimeoutHours == 0 {
		return types.ErrInvalidVerificationTimeout
	}

	if data.Params.MaxCredentialsPerDid == 0 {
		return types.ErrInvalidMaxCredentials
	}

	if len(data.Params.AllowedOauthProviders) == 0 {
		return types.ErrNoOAuthProvidersAllowed
	}

	if len(data.Params.AllowedCredentialTypes) == 0 {
		return types.ErrNoCredentialTypesAllowed
	}

	// Validate DID documents
	for _, doc := range data.Dids {
		if doc.Id == "" {
			return types.ErrInvalidDID
		}
	}

	// Validate credentials
	for _, cred := range data.Credentials {
		if cred.Id == "" {
			return types.ErrInvalidCredential
		}
		if cred.Subject == "" || cred.Issuer == "" {
			return types.ErrInvalidDID
		}
		if cred.SchemaId == "" {
			return types.ErrInvalidCredential
		}
		if cred.Created.IsZero() {
			return types.ErrInvalidCredential
		}
	}

	// Validate credential schemas
	for _, schema := range data.CredentialSchemas {
		if schema.Id == "" {
			return types.ErrInvalidCredentialSchema
		}
		if schema.Name == "" {
			return types.ErrInvalidCredentialSchema
		}
		if len(schema.Properties) == 0 {
			return types.ErrInvalidCredentialSchema
		}
	}

	return nil
}

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState *types.GenesisState) {
	// Set all the DID documents
	for _, doc := range genState.Dids {
		k.SetDIDDocument(ctx, doc.Id, doc)
	}

	// Set all the credentials
	for _, cred := range genState.Credentials {
		k.SetCredential(ctx, cred)
	}

	// Set all the credential schemas
	for _, schema := range genState.CredentialSchemas {
		k.SetCredentialSchema(ctx, schema)
	}

	// Set all the social identities
	for _, identity := range genState.SocialIdentities {
		k.StoreSocialIdentity(ctx, identity)
	}

	// Set module parameters
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.GenesisState{
		Params:             k.GetParams(ctx),
		Dids:              k.GetAllDIDDocuments(ctx),
		Credentials:        k.GetAllCredentials(ctx),
		SocialIdentities:   k.GetAllSocialIdentities(ctx),
		CredentialSchemas: k.GetAllCredentialSchemas(ctx),
	}

	return &genesis
}
