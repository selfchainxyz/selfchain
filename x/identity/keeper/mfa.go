package keeper

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pquerna/otp/totp"
	"selfchain/x/identity/types"
)

const (
	MFAConfigPrefix    = "mfa_config/"
	MFAChallengePrefix = "mfa_challenge/"
	MFASecretLength    = 20
	MFAChallengeExpiry = 5 * time.Minute
)

// SetMFAConfig stores an MFA configuration
func (k Keeper) SetMFAConfig(ctx sdk.Context, config types.MFAConfig) error {
	if !k.HasDIDDocument(ctx, config.Did) {
		return fmt.Errorf("DID not found: %s", config.Did)
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(MFAConfigPrefix))
	configBytes := k.cdc.MustMarshal(&config)
	store.Set([]byte(config.Did), configBytes)
	return nil
}

// GetMFAConfig retrieves an MFA configuration
func (k Keeper) GetMFAConfig(ctx sdk.Context, did string) (types.MFAConfig, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(MFAConfigPrefix))
	configBytes := store.Get([]byte(did))
	if configBytes == nil {
		return types.MFAConfig{}, false
	}

	var config types.MFAConfig
	k.cdc.MustUnmarshal(configBytes, &config)
	return config, true
}

// EnableMFA enables MFA for a DID
func (k Keeper) EnableMFA(ctx sdk.Context, did string, methodType string, identifier string) (*types.MFAMethod, error) {
	config, found := k.GetMFAConfig(ctx, did)
	if !found {
		config = types.MFAConfig{
			Did:     did,
			Enabled: false,
			Methods: []*types.MFAMethod{},
		}
	}

	// Generate secret based on method type
	secret := make([]byte, MFASecretLength)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("failed to generate secret: %v", err)
	}

	method := &types.MFAMethod{
		Type:       methodType,
		Identifier: identifier,
		Secret:     secret,
		Verified:   false,
		CreatedAt:  ctx.BlockTime(),
	}

	config.Methods = append(config.Methods, method)
	config.UpdatedAt = ctx.BlockTime()

	if err := k.SetMFAConfig(ctx, config); err != nil {
		return nil, err
	}

	return method, nil
}

// CreateMFAChallenge creates a new MFA challenge
func (k Keeper) CreateMFAChallenge(ctx sdk.Context, did string, methodType string) (*types.MFAChallenge, error) {
	config, found := k.GetMFAConfig(ctx, did)
	if !found || !config.Enabled {
		return nil, fmt.Errorf("MFA not enabled for DID: %s", did)
	}

	// Find the specified method
	var method *types.MFAMethod
	for _, m := range config.Methods {
		if m.Type == methodType && m.Verified {
			method = m
			break
		}
	}
	if method == nil {
		return nil, fmt.Errorf("no verified MFA method found of type: %s", methodType)
	}

	// Generate challenge ID
	challengeID := make([]byte, 16)
	if _, err := rand.Read(challengeID); err != nil {
		return nil, fmt.Errorf("failed to generate challenge ID: %v", err)
	}

	challenge := &types.MFAChallenge{
		Id:            base32.StdEncoding.EncodeToString(challengeID),
		Did:           did,
		MethodType:    methodType,
		ChallengeData: method.Secret, // For TOTP, this is the shared secret
		ExpiresAt:     ctx.BlockTime().Add(MFAChallengeExpiry),
		Completed:     false,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(MFAChallengePrefix))
	challengeBytes := k.cdc.MustMarshal(challenge)
	store.Set([]byte(challenge.Id), challengeBytes)

	return challenge, nil
}

// GetMFAChallenge retrieves an MFA challenge by ID
func (k Keeper) GetMFAChallenge(ctx sdk.Context, id string) (types.MFAChallenge, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(MFAChallengePrefix))
	challengeBytes := store.Get([]byte(id))
	if challengeBytes == nil {
		return types.MFAChallenge{}, false
	}

	var challenge types.MFAChallenge
	k.cdc.MustUnmarshal(challengeBytes, &challenge)
	return challenge, true
}

// VerifyMFAChallenge verifies an MFA challenge response
func (k Keeper) VerifyMFAChallenge(ctx sdk.Context, challengeID string, response string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(MFAChallengePrefix))
	challengeBytes := store.Get([]byte(challengeID))
	if challengeBytes == nil {
		return fmt.Errorf("challenge not found: %s", challengeID)
	}

	var challenge types.MFAChallenge
	k.cdc.MustUnmarshal(challengeBytes, &challenge)

	if ctx.BlockTime().After(challenge.ExpiresAt) {
		return fmt.Errorf("challenge expired")
	}

	if challenge.Completed {
		return fmt.Errorf("challenge already completed")
	}

	// Verify based on method type
	switch challenge.MethodType {
	case "totp":
		secret := base32.StdEncoding.EncodeToString(challenge.ChallengeData)
		if !totp.Validate(response, secret) {
			return fmt.Errorf("invalid TOTP code")
		}
	case "recovery_codes":
		// Implement recovery code verification
		return fmt.Errorf("recovery codes not implemented yet")
	default:
		return fmt.Errorf("unsupported MFA method type: %s", challenge.MethodType)
	}

	challenge.Completed = true
	challengeBytes = k.cdc.MustMarshal(&challenge)
	store.Set([]byte(challengeID), challengeBytes)

	return nil
}
