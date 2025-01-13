package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	MFAParamsKey        = []byte("MFAParams")
	CredentialParamsKey = []byte("CredentialParams")
	DIDParamsKey        = []byte("DIDParams")
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
		AllowedMethods:         []string{"selfchain"},
		MaxControllers:         5,
		MaxServices:            10,
		MaxVerificationMethods: 10,
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

// validateMFAParams validates MFAParams
func validateMFAParams(params MFAParams) error {
	if params.MaxMethods <= 0 {
		return fmt.Errorf("max methods must be positive")
	}
	if params.ChallengeExpiry <= 0 {
		return fmt.Errorf("challenge expiry must be positive")
	}
	if len(params.AllowedMethods) == 0 {
		return fmt.Errorf("allowed methods cannot be empty")
	}
	if params.MaxFailedAttempts <= 0 {
		return fmt.Errorf("max failed attempts must be positive")
	}
	return nil
}

// validateCredentialParams validates CredentialParams
func validateCredentialParams(params CredentialParams) error {
	if params.MaxCredentialsPerDid <= 0 {
		return fmt.Errorf("max credentials per DID must be positive")
	}
	if params.MaxClaimSize <= 0 {
		return fmt.Errorf("max claim size must be positive")
	}
	if len(params.AllowedTypes) == 0 {
		return fmt.Errorf("allowed types cannot be empty")
	}
	if params.MaxValidityDuration <= 0 {
		return fmt.Errorf("max validity duration must be positive")
	}
	return nil
}

// validateDIDParams validates DIDParams
func validateDIDParams(params DIDParams) error {
	if len(params.AllowedMethods) == 0 {
		return fmt.Errorf("allowed methods cannot be empty")
	}
	if params.MaxControllers <= 0 {
		return fmt.Errorf("max controllers must be positive")
	}
	if params.MaxServices <= 0 {
		return fmt.Errorf("max services must be positive")
	}
	if params.MaxVerificationMethods <= 0 {
		return fmt.Errorf("max verification methods must be positive")
	}
	return nil
}

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}
