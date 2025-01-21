package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyMFAParams        = []byte("MFAParams")
	KeyCredentialParams = []byte("CredentialParams")
	KeyDIDParams        = []byte("DIDParams")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns default identity module parameters
func DefaultParams() Params {
	return Params{
		MfaParams:        DefaultMFAParams(),
		CredentialParams: DefaultCredentialParams(),
		DidParams:        DefaultDIDParams(),
	}
}

// DefaultMFAParams returns default MFA parameters
func DefaultMFAParams() MFAParams {
	return MFAParams{
		MaxMethods:        3,
		ChallengeExpiry:   int64(5 * time.Minute),
		AllowedMethods:    []string{"totp", "recovery"},
		MaxFailedAttempts: 5,
	}
}

// DefaultCredentialParams returns default credential parameters
func DefaultCredentialParams() CredentialParams {
	return CredentialParams{
		MaxCredentialsPerDid: 100,
		MaxClaimSize:         1024 * 1024, // 1MB
		AllowedTypes:         []string{"VerifiableCredential"},
		MaxValidityDuration:  int64(365 * 24 * time.Hour), // 1 year
	}
}

// DefaultDIDParams returns default DID parameters
func DefaultDIDParams() DIDParams {
	return DIDParams{
		AllowedMethods:         []string{"did:self"},
		MaxControllers:         5,
		MaxServices:           10,
		MaxVerificationMethods: 5,
	}
}

// Validate validates Params
func (p Params) Validate() error {
	if err := validateMFAParams(p.MfaParams); err != nil {
		return err
	}
	if err := validateCredentialParams(p.CredentialParams); err != nil {
		return err
	}
	if err := validateDIDParams(p.DidParams); err != nil {
		return err
	}
	return nil
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMFAParams, p.MfaParams, validateMFAParams),
		paramtypes.NewParamSetPair(KeyCredentialParams, p.CredentialParams, validateCredentialParams),
		paramtypes.NewParamSetPair(KeyDIDParams, p.DidParams, validateDIDParams),
	}
}

// validateMFAParams validates MFAParams
func validateMFAParams(i interface{}) error {
	v, ok := i.(MFAParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.MaxMethods <= 0 {
		return fmt.Errorf("max methods must be positive")
	}

	if v.ChallengeExpiry <= 0 {
		return fmt.Errorf("challenge expiry must be positive")
	}

	if len(v.AllowedMethods) == 0 {
		return fmt.Errorf("allowed methods cannot be empty")
	}

	if v.MaxFailedAttempts <= 0 {
		return fmt.Errorf("max failed attempts must be positive")
	}

	return nil
}

// validateCredentialParams validates CredentialParams
func validateCredentialParams(i interface{}) error {
	v, ok := i.(CredentialParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.MaxCredentialsPerDid <= 0 {
		return fmt.Errorf("max credentials per did must be positive")
	}

	if v.MaxClaimSize <= 0 {
		return fmt.Errorf("max claim size must be positive")
	}

	if len(v.AllowedTypes) == 0 {
		return fmt.Errorf("allowed types cannot be empty")
	}

	if v.MaxValidityDuration <= 0 {
		return fmt.Errorf("max validity duration must be positive")
	}

	return nil
}

// validateDIDParams validates DIDParams
func validateDIDParams(i interface{}) error {
	v, ok := i.(DIDParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.MaxServices <= 0 {
		return fmt.Errorf("max services must be positive")
	}

	if v.MaxVerificationMethods <= 0 {
		return fmt.Errorf("max verification methods must be positive")
	}

	if v.MaxControllers <= 0 {
		return fmt.Errorf("max controllers must be positive")
	}

	return nil
}

var _ paramtypes.ParamSet = (*Params)(nil)

// validateParams validates the params
func validateParams(i interface{}) error {
	v, ok := i.(*Params)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return v.Validate()
}
