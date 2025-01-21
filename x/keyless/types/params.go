package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	// Default values
	DefaultMaxParties       = uint32(5)
	DefaultMaxThreshold     = uint32(3)
	DefaultMaxSecurityLevel = uint32(3)
	DefaultMaxBatchSize     = uint32(100)
	DefaultMaxMetadataSize  = uint32(1024) // 1KB
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		MaxParties:       DefaultMaxParties,
		MaxThreshold:     DefaultMaxThreshold,
		MaxSecurityLevel: DefaultMaxSecurityLevel,
		MaxBatchSize:     DefaultMaxBatchSize,
		MaxMetadataSize:  DefaultMaxMetadataSize,
	}
}

// NewParams creates a new Params instance
func NewParams(
	maxParties uint32,
	maxThreshold uint32,
	maxSecurityLevel uint32,
	maxBatchSize uint32,
	maxMetadataSize uint32,
) Params {
	return Params{
		MaxParties:       maxParties,
		MaxThreshold:     maxThreshold,
		MaxSecurityLevel: maxSecurityLevel,
		MaxBatchSize:     maxBatchSize,
		MaxMetadataSize:  maxMetadataSize,
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
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
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
