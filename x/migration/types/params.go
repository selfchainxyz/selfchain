package types

import (
	fmt "fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

const (
	DefaultHotcrossRatio uint64 = 0
)

var (
	HotcrossRatio = []byte("HotcrossRatio")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(hotcrossRatio uint64) Params {
	return Params{
		HotcrossRatio: hotcrossRatio,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultHotcrossRatio,
	)
}

func validateHotcrossRatio(i interface {}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v > 100 {
		return fmt.Errorf("invalid ratio value")
	}

	return nil
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(HotcrossRatio, &p.HotcrossRatio, validateHotcrossRatio),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateHotcrossRatio(p.HotcrossRatio); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
