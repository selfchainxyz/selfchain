package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pquerna/otp/totp"
	"selfchain/x/identity/types"
)

const (
	// Key prefixes
	MFAConfigPrefix    = types.MFAConfigPrefix
	MFAChallengePrefix = types.MFAChallengePrefix
	MFAChallengeExpiry = types.MFAChallengeExpiry
	MaxMFAMethods      = 5 // Maximum number of MFA methods per DID
)

// StoreMFAConfig stores an MFA configuration
func (k Keeper) StoreMFAConfig(ctx sdk.Context, config types.MFAConfig) error {
	if err := config.ValidateBasic(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	key := append([]byte(MFAConfigPrefix), []byte(config.Did)...)
	bz, err := k.cdc.Marshal(&config)
	if err != nil {
		return err
	}

	store.Set(key, bz)
	return nil
}

// GetMFAConfig returns an MFA configuration
func (k Keeper) GetMFAConfig(ctx sdk.Context, did string) (*types.MFAConfig, bool) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(MFAConfigPrefix), []byte(did)...)
	bz := store.Get(key)
	if bz == nil {
		return nil, false
	}

	var config types.MFAConfig
	k.cdc.MustUnmarshal(bz, &config)
	return &config, true
}

// DeleteMFAConfig deletes an MFA configuration
func (k Keeper) DeleteMFAConfig(ctx sdk.Context, did string) error {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(MFAConfigPrefix), []byte(did)...)
	store.Delete(key)
	return nil
}

// AddMFAMethod adds an MFA method to a DID
func (k Keeper) AddMFAMethod(ctx sdk.Context, did string, method types.MFAMethod) error {
	config, found := k.GetMFAConfig(ctx, did)
	if !found {
		// Create new MFA config if it doesn't exist
		mfaConfig := types.MFAConfig{
			Did:     did,
			Methods: make([]*types.MFAMethod, 0),
		}
		config = &mfaConfig
	}

	// Check if method already exists
	for _, m := range config.Methods {
		if m.Type == method.Type {
			return sdkerrors.Wrap(types.ErrInvalidMFAMethod, "method already exists")
		}
	}

	// Check max methods
	if len(config.Methods) >= MaxMFAMethods {
		return sdkerrors.Wrap(types.ErrInvalidMFAMethod, "max methods exceeded")
	}

	// Create new MFA method
	blockTime := ctx.BlockTime()
	mfaMethod := types.MFAMethod{
		Type:      method.Type,
		Secret:    method.Secret,
		CreatedAt: &blockTime,
		Status:    types.MFAMethodStatus_MFA_METHOD_STATUS_ACTIVE,
	}

	// Add method to config
	config.Methods = append(config.Methods, &mfaMethod)
	return k.StoreMFAConfig(ctx, *config)
}

// RemoveMFAMethod removes an MFA method from a DID
func (k Keeper) RemoveMFAMethod(ctx sdk.Context, did string, methodType string) error {
	config, found := k.GetMFAConfig(ctx, did)
	if !found {
		return sdkerrors.Wrap(types.ErrMFAMethodNotFound, "MFA config not found")
	}

	// Find and remove method
	for i, method := range config.Methods {
		if method.Type == methodType {
			config.Methods = append(config.Methods[:i], config.Methods[i+1:]...)
			return k.StoreMFAConfig(ctx, *config)
		}
	}

	return sdkerrors.Wrap(types.ErrMFAMethodNotFound, "method not found")
}

// GetMFAMethod gets an MFA method for a DID
func (k Keeper) GetMFAMethod(ctx sdk.Context, did string, methodType string) (*types.MFAMethod, error) {
	config, found := k.GetMFAConfig(ctx, did)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrMFAMethodNotFound, "MFA config not found")
	}

	// Find method
	for _, method := range config.Methods {
		if method.Type == methodType {
			return method, nil
		}
	}

	return nil, sdkerrors.Wrap(types.ErrMFAMethodNotFound, "method not found")
}

// CreateMFAChallenge creates a new MFA challenge
func (k Keeper) CreateMFAChallenge(ctx sdk.Context, did string, methodType string) (*types.MFAChallenge, error) {
	method, err := k.GetMFAMethod(ctx, did, methodType)
	if err != nil {
		return nil, err
	}

	if !method.IsActive() {
		return nil, sdkerrors.Wrap(types.ErrMFAMethodInactive, "method is not active")
	}

	// Create challenge
	blockTime := ctx.BlockTime()
	expiry := blockTime.Add(MFAChallengeExpiry)
	challenge := types.MFAChallenge{
		Id:        fmt.Sprintf("%s:%s", did, methodType),
		Did:       did,
		Method:    methodType,
		CreatedAt: &blockTime,
		ExpiresAt: &expiry,
	}

	// Store challenge
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(MFAChallengePrefix), []byte(challenge.Id)...)
	bz, err := k.cdc.Marshal(&challenge)
	if err != nil {
		return nil, err
	}

	store.Set(key, bz)
	return &challenge, nil
}

// GetMFAChallenge returns an MFA challenge
func (k Keeper) GetMFAChallenge(ctx sdk.Context, did string, methodType string) (*types.MFAChallenge, error) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(MFAChallengePrefix), []byte(fmt.Sprintf("%s:%s", did, methodType))...)
	bz := store.Get(key)
	if bz == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMFAMethod, "challenge not found")
	}

	var challenge types.MFAChallenge
	k.cdc.MustUnmarshal(bz, &challenge)

	// Check if challenge has expired
	if challenge.ExpiresAt.Before(ctx.BlockTime()) {
		return nil, sdkerrors.Wrap(types.ErrInvalidMFAMethod, "challenge expired")
	}

	return &challenge, nil
}

// DeleteMFAChallenge deletes an MFA challenge
func (k Keeper) DeleteMFAChallenge(ctx sdk.Context, did string, methodType string) error {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(MFAChallengePrefix), []byte(fmt.Sprintf("%s:%s", did, methodType))...)
	store.Delete(key)
	return nil
}

// VerifyMFACode verifies an MFA code for a given method
func (k Keeper) VerifyMFACode(ctx sdk.Context, did string, methodID string, code string) error {
	method, err := k.GetMFAMethod(ctx, did, methodID)
	if err != nil {
		return err
	}

	if !method.IsActive() {
		return sdkerrors.Wrap(types.ErrMFAMethodDisabled, "method is disabled")
	}

	var verifyErr error
	switch method.Type {
	case "OTP":
		verifyErr = k.verifyOTP(ctx, method, code)
	case "TOTP":
		verifyErr = k.verifyTOTP(ctx, method, code)
	default:
		return sdkerrors.Wrap(types.ErrInvalidMFAMethod, "unsupported MFA method type")
	}

	if verifyErr != nil {
		// Emit failure event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"mfa_verification",
				sdk.NewAttribute("method_id", methodID),
				sdk.NewAttribute("error", verifyErr.Error()),
			),
		)
		return verifyErr
	}

	// Emit success event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"mfa_verification",
			sdk.NewAttribute("method_id", methodID),
		),
	)

	return nil
}

// verifyOTP verifies a one-time password
func (k Keeper) verifyOTP(ctx sdk.Context, method *types.MFAMethod, code string) error {
	if method.Type != "OTP" {
		return sdkerrors.Wrap(types.ErrInvalidMFAMethod, "invalid OTP method type")
	}

	// Implement OTP verification logic here
	// For now, return unimplemented
	return sdkerrors.Wrap(types.ErrInvalidMFAMethod, "OTP verification not implemented")
}

// verifyTOTP verifies a time-based one-time password
func (k Keeper) verifyTOTP(ctx sdk.Context, method *types.MFAMethod, code string) error {
	if method.Type != "TOTP" {
		return sdkerrors.Wrap(types.ErrInvalidMFAMethod, "invalid TOTP method type")
	}

	// Verify TOTP code
	valid := totp.Validate(code, method.Secret)
	if !valid {
		return sdkerrors.Wrap(types.ErrInvalidMFAMethod, "invalid TOTP code")
	}

	return nil
}

// VerifyMFAChallenge verifies an MFA challenge
func (k Keeper) VerifyMFAChallenge(ctx sdk.Context, did string, methodType string, code string) error {
	// Get challenge
	challenge, err := k.GetMFAChallenge(ctx, did, methodType)
	if err != nil {
		return err
	}

	// Get method
	method, err := k.GetMFAMethod(ctx, did, methodType)
	if err != nil {
		return err
	}

	// Verify challenge hasn't expired
	if challenge.IsExpired() {
		return sdkerrors.Wrap(types.ErrInvalidMFAMethod, "challenge expired")
	}

	// Verify the code using the appropriate method
	switch methodType {
	case "totp":
		err := k.verifyTOTP(ctx, method, code)
		if err != nil {
			return sdkerrors.Wrap(types.ErrMFAVerificationFailed, "invalid TOTP code")
		}
	default:
		return sdkerrors.Wrap(types.ErrInvalidMFAMethod, "unsupported MFA method type")
	}

	// Delete the challenge after successful verification
	k.DeleteMFAChallenge(ctx, did, methodType)

	return nil
}

// GetAllMFAConfigs returns all MFA configurations
func (k Keeper) GetAllMFAConfigs(ctx sdk.Context) []types.MFAConfig {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(MFAConfigPrefix))
	defer iterator.Close()

	var configs []types.MFAConfig
	for ; iterator.Valid(); iterator.Next() {
		var config types.MFAConfig
		k.cdc.MustUnmarshal(iterator.Value(), &config)
		configs = append(configs, config)
	}
	return configs
}

// SetMFAMethod stores an MFA method for a DID
func (k Keeper) SetMFAMethod(ctx sdk.Context, did string, method string, secret string) error {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.MFAMethodPrefix), []byte(did+":"+method)...)

	blockTime := ctx.BlockTime()
	mfaMethod := types.MFAMethod{
		Type:      method,
		Secret:    secret,
		CreatedAt: &blockTime,
		Status:    types.MFAMethodStatus_MFA_METHOD_STATUS_ACTIVE,
	}

	bz, err := json.Marshal(mfaMethod)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to marshal MFA method")
	}

	store.Set(key, bz)
	return nil
}
