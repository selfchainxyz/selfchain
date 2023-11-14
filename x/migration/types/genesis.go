package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		TokenMigrationList: []TokenMigration{},
		Acl:                nil,
		MigratorList:       []Migrator{},
		Config:             nil,
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in tokenMigration
	tokenMigrationIndexMap := make(map[string]struct{})

	for _, elem := range gs.TokenMigrationList {
		index := string(TokenMigrationKey(elem.MsgHash))
		if _, ok := tokenMigrationIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for tokenMigration")
		}
		tokenMigrationIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in migrator
	migratorIndexMap := make(map[string]struct{})

	for _, elem := range gs.MigratorList {
		index := string(MigratorKey(elem.Migrator))
		if _, ok := migratorIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for migrator")
		}
		migratorIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
