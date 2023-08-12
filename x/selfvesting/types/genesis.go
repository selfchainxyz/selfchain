package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		VestingPositionsList: []VestingPositions{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in vestingPositions
	vestingPositionsIndexMap := make(map[string]struct{})

	for _, elem := range gs.VestingPositionsList {
		index := string(VestingPositionsKey(elem.Beneficiary))
		if _, ok := vestingPositionsIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for vestingPositions")
		}
		vestingPositionsIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
