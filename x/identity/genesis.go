package identity

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/identity/keeper"
	"selfchain/x/identity/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *types.GenesisState {
	return &types.GenesisState{
		Params: types.Params{
			MfaParams: types.MFAParams{
				MaxMethods:        3,
				ChallengeExpiry:  300, // 5 minutes
				AllowedMethods:   []string{"totp", "email"},
				MaxFailedAttempts: 3,
			},
			CredentialParams: types.CredentialParams{
				MaxCredentialsPerDid: 100,
				MaxClaimSize:        1024 * 1024, // 1MB
				AllowedTypes:        []string{"VerifiableCredential"},
				MaxValidityDuration: 31536000, // 1 year
			},
			DidParams: types.DIDParams{
				AllowedMethods:         []string{"self"},
				MaxControllers:         5,
				MaxServices:           10,
				MaxVerificationMethods: 10,
			},
		},
		DidDocuments: []types.DIDDocument{},
		Credentials:  []types.Credential{},
		MfaConfigs:   []types.MFAConfig{},
		AuditLogs:    []types.AuditLogEntry{},
		SocialIdentities: []types.SocialIdentity{},
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data *types.GenesisState) error {
	// Validate MFA params
	if data.Params.MfaParams.MaxMethods <= 0 {
		return types.ErrInvalidMFAMethod.Wrap("max methods must be positive")
	}
	if data.Params.MfaParams.ChallengeExpiry <= 0 {
		return types.ErrInvalidMFAMethod.Wrap("challenge expiry must be positive")
	}
	if len(data.Params.MfaParams.AllowedMethods) == 0 {
		return types.ErrInvalidMFAMethod.Wrap("no MFA methods allowed")
	}

	// Validate credential params
	if data.Params.CredentialParams.MaxCredentialsPerDid <= 0 {
		return types.ErrInvalidCredentialID.Wrap("max credentials per DID must be positive")
	}
	if data.Params.CredentialParams.MaxClaimSize <= 0 {
		return types.ErrInvalidCredentialClaims.Wrap("max claim size must be positive")
	}
	if len(data.Params.CredentialParams.AllowedTypes) == 0 {
		return types.ErrInvalidCredentialType.Wrap("no credential types allowed")
	}

	// Validate DID params
	if len(data.Params.DidParams.AllowedMethods) == 0 {
		return types.ErrInvalidDID.Wrap("no DID methods allowed")
	}
	if data.Params.DidParams.MaxControllers <= 0 {
		return types.ErrInvalidController.Wrap("max controllers must be positive")
	}
	if data.Params.DidParams.MaxServices <= 0 {
		return types.ErrInvalidService.Wrap("max services must be positive")
	}
	if data.Params.DidParams.MaxVerificationMethods <= 0 {
		return types.ErrInvalidVerificationMethod.Wrap("max verification methods must be positive")
	}

	// Validate DID documents
	for _, doc := range data.DidDocuments {
		if err := doc.ValidateBasic(); err != nil {
			return types.ErrInvalidDID.Wrapf("invalid DID document: %v", err)
		}
	}

	// Validate credentials
	for _, cred := range data.Credentials {
		if err := cred.ValidateBasic(); err != nil {
			return types.ErrInvalidCredentialID.Wrapf("invalid credential: %v", err)
		}
	}

	// Validate MFA configs
	for _, config := range data.MfaConfigs {
		if config.Did == "" {
			return types.ErrInvalidMFAConfig.Wrap("DID cannot be empty")
		}
		for _, method := range config.Methods {
			if method.GetType() == "" {
				return types.ErrInvalidMFAConfig.Wrap("method type cannot be empty")
			}
		}
	}

	return nil
}

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Validate DID documents
	for _, doc := range genState.DidDocuments {
		if err := doc.ValidateBasic(); err != nil {
			panic(fmt.Sprintf("invalid DID document: %v", err))
		}
		k.SetDIDDocument(ctx, doc.Id, doc)
	}

	// Validate credentials
	for _, cred := range genState.Credentials {
		if err := cred.ValidateBasic(); err != nil {
			panic(fmt.Sprintf("invalid credential: %v", err))
		}
		k.SetCredential(ctx, &cred)
	}

	// Store social identities
	for _, identity := range genState.SocialIdentities {
		k.StoreSocialIdentity(ctx, identity)
	}

	// Store MFA configurations
	for _, config := range genState.MfaConfigs {
		k.StoreMFAConfig(ctx, config)
	}

	// Set module parameters
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Get all DID documents
	didDocuments := k.GetAllDIDDocuments(ctx)
	genesis.DidDocuments = didDocuments

	// Get all credentials
	credentials, err := k.GetAllCredentials(ctx)
	if err != nil {
		panic(err)
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
