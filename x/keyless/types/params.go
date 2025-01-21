package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	// Default values for wallet limits
	DefaultMaxParties       = uint32(5)
	DefaultMaxThreshold     = uint32(3)
	DefaultMaxSecurityLevel = uint32(3)
	DefaultMaxBatchSize     = uint32(100)
	DefaultMaxMetadataSize  = uint32(1024) // 1KB

	// Default values for recovery settings
	DefaultMaxWalletsPerDid     = uint32(5)
	DefaultMaxSharesPerWallet   = uint32(3)
	DefaultMinRecoveryThreshold = uint32(2)
	DefaultMaxRecoveryThreshold = uint32(3)
	DefaultRecoveryWindowSeconds = uint32(86400) // 24 hours
	DefaultMaxSigningAttempts    = uint32(3)
)

// Parameter store keys
var (
	// Wallet limit keys
	KeyMaxParties       = []byte("MaxParties")
	KeyMaxThreshold     = []byte("MaxThreshold")
	KeyMaxSecurityLevel = []byte("MaxSecurityLevel")
	KeyMaxBatchSize     = []byte("MaxBatchSize")
	KeyMaxMetadataSize  = []byte("MaxMetadataSize")

	// Recovery setting keys
	KeyMaxWalletsPerDid      = []byte("MaxWalletsPerDid")
	KeyMaxSharesPerWallet    = []byte("MaxSharesPerWallet")
	KeyMinRecoveryThreshold  = []byte("MinRecoveryThreshold")
	KeyMaxRecoveryThreshold  = []byte("MaxRecoveryThreshold")
	KeyRecoveryWindowSeconds = []byte("RecoveryWindowSeconds")
	KeyMaxSigningAttempts    = []byte("MaxSigningAttempts")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		// Wallet limits
		MaxParties:       DefaultMaxParties,
		MaxThreshold:     DefaultMaxThreshold,
		MaxSecurityLevel: DefaultMaxSecurityLevel,
		MaxBatchSize:     DefaultMaxBatchSize,
		MaxMetadataSize:  DefaultMaxMetadataSize,

		// Recovery settings
		MaxWalletsPerDid:      DefaultMaxWalletsPerDid,
		MaxSharesPerWallet:    DefaultMaxSharesPerWallet,
		MinRecoveryThreshold:  DefaultMinRecoveryThreshold,
		MaxRecoveryThreshold:  DefaultMaxRecoveryThreshold,
		RecoveryWindowSeconds: DefaultRecoveryWindowSeconds,
		MaxSigningAttempts:    DefaultMaxSigningAttempts,
	}
}

// NewParams creates a new Params instance
func NewParams(
	maxParties uint32,           // Wallet limit params
	maxThreshold uint32,
	maxSecurityLevel uint32,
	maxBatchSize uint32,
	maxMetadataSize uint32,
	maxWalletsPerDid uint32,     // Recovery params
	maxSharesPerWallet uint32,
	minRecoveryThreshold uint32,
	maxRecoveryThreshold uint32,
	recoveryWindowSeconds uint32,
	maxSigningAttempts uint32,
) Params {
	return Params{
		// Wallet limits
		MaxParties:       maxParties,
		MaxThreshold:     maxThreshold,
		MaxSecurityLevel: maxSecurityLevel,
		MaxBatchSize:     maxBatchSize,
		MaxMetadataSize:  maxMetadataSize,

		// Recovery settings
		MaxWalletsPerDid:      maxWalletsPerDid,
		MaxSharesPerWallet:    maxSharesPerWallet,
		MinRecoveryThreshold:  minRecoveryThreshold,
		MaxRecoveryThreshold:  maxRecoveryThreshold,
		RecoveryWindowSeconds: recoveryWindowSeconds,
		MaxSigningAttempts:    maxSigningAttempts,
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxParties, &p.MaxParties, validateUint32Param),
		paramtypes.NewParamSetPair(KeyMaxThreshold, &p.MaxThreshold, validateUint32Param),
		paramtypes.NewParamSetPair(KeyMaxSecurityLevel, &p.MaxSecurityLevel, validateUint32Param),
		paramtypes.NewParamSetPair(KeyMaxBatchSize, &p.MaxBatchSize, validateUint32Param),
		paramtypes.NewParamSetPair(KeyMaxMetadataSize, &p.MaxMetadataSize, validateUint32Param),
		paramtypes.NewParamSetPair(KeyMaxWalletsPerDid, &p.MaxWalletsPerDid, validateUint32Param),
		paramtypes.NewParamSetPair(KeyMaxSharesPerWallet, &p.MaxSharesPerWallet, validateUint32Param),
		paramtypes.NewParamSetPair(KeyMinRecoveryThreshold, &p.MinRecoveryThreshold, validateUint32Param),
		paramtypes.NewParamSetPair(KeyMaxRecoveryThreshold, &p.MaxRecoveryThreshold, validateUint32Param),
		paramtypes.NewParamSetPair(KeyRecoveryWindowSeconds, &p.RecoveryWindowSeconds, validateUint32Param),
		paramtypes.NewParamSetPair(KeyMaxSigningAttempts, &p.MaxSigningAttempts, validateUint32Param),
	}
}

func validateUint32Param(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return ErrInvalidParam.Wrap("invalid parameter type: expected uint32")
	}

	if v == 0 {
		return ErrInvalidParam.Wrap("parameter must be greater than 0")
	}

	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	// Validate wallet limits
	if p.MaxParties == 0 {
		return ErrInvalidParam.Wrap("MaxParties must be greater than 0")
	}
	if p.MaxThreshold == 0 || p.MaxThreshold > p.MaxParties {
		return ErrInvalidParam.Wrap("MaxThreshold must be greater than 0 and less than or equal to MaxParties")
	}
	if p.MaxSecurityLevel == 0 {
		return ErrInvalidParam.Wrap("MaxSecurityLevel must be greater than 0")
	}
	if p.MaxBatchSize == 0 {
		return ErrInvalidParam.Wrap("MaxBatchSize must be greater than 0")
	}
	if p.MaxMetadataSize == 0 {
		return ErrInvalidParam.Wrap("MaxMetadataSize must be greater than 0")
	}

	// Validate recovery settings
	if p.MaxWalletsPerDid == 0 {
		return ErrInvalidParam.Wrap("MaxWalletsPerDid must be greater than 0")
	}
	if p.MaxSharesPerWallet == 0 {
		return ErrInvalidParam.Wrap("MaxSharesPerWallet must be greater than 0")
	}
	if p.MinRecoveryThreshold == 0 || p.MinRecoveryThreshold > p.MaxSharesPerWallet {
		return ErrInvalidParam.Wrap("MinRecoveryThreshold must be greater than 0 and less than or equal to MaxSharesPerWallet")
	}
	if p.MaxRecoveryThreshold < p.MinRecoveryThreshold || p.MaxRecoveryThreshold > p.MaxSharesPerWallet {
		return ErrInvalidParam.Wrap("MaxRecoveryThreshold must be greater than or equal to MinRecoveryThreshold and less than or equal to MaxSharesPerWallet")
	}
	if p.RecoveryWindowSeconds == 0 {
		return ErrInvalidParam.Wrap("RecoveryWindowSeconds must be greater than 0")
	}
	if p.MaxSigningAttempts == 0 {
		return ErrInvalidParam.Wrap("MaxSigningAttempts must be greater than 0")
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
