package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	// Default values
	DefaultMaxWalletsPerDID     = uint32(5)
	DefaultMaxSharesPerWallet   = uint32(3)
	DefaultMinRecoveryThreshold = uint32(2)
	DefaultMaxRecoveryThreshold = uint32(3)
	DefaultRecoveryWindowSecs   = uint32(86400) // 24 hours
	DefaultMaxSigningAttempts   = uint32(3)
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		MaxWalletsPerDid:     DefaultMaxWalletsPerDID,
		MaxSharesPerWallet:   DefaultMaxSharesPerWallet,
		MinRecoveryThreshold: DefaultMinRecoveryThreshold,
		MaxRecoveryThreshold: DefaultMaxRecoveryThreshold,
		RecoveryWindowSeconds: DefaultRecoveryWindowSecs,
		MaxSigningAttempts:   DefaultMaxSigningAttempts,
	}
}

// NewParams creates a new Params instance
func NewParams(
	maxWallets uint32,
	maxShares uint32,
	minThreshold uint32,
	maxThreshold uint32,
	windowSecs uint32,
	maxAttempts uint32,
) Params {
	return Params{
		MaxWalletsPerDid:     maxWallets,
		MaxSharesPerWallet:   maxShares,
		MinRecoveryThreshold: minThreshold,
		MaxRecoveryThreshold: maxThreshold,
		RecoveryWindowSeconds: windowSecs,
		MaxSigningAttempts:   maxAttempts,
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.MaxWalletsPerDid == 0 {
		return ErrInvalidMaxWallets
	}
	if p.MaxSharesPerWallet == 0 {
		return ErrInvalidMaxShares
	}
	if p.MinRecoveryThreshold == 0 || p.MinRecoveryThreshold > p.MaxRecoveryThreshold {
		return ErrInvalidRecoveryThreshold
	}
	if p.MaxRecoveryThreshold > p.MaxSharesPerWallet {
		return ErrInvalidRecoveryThreshold
	}
	if p.RecoveryWindowSeconds == 0 {
		return ErrInvalidRecoveryWindow
	}
	if p.MaxSigningAttempts == 0 {
		return ErrInvalidMaxAttempts
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
