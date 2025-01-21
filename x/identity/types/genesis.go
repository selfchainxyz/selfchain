package types

import (
	"fmt"
	// this line is used by starport scaffolding # genesis/types/import
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:       DefaultParams(),
		DidDocuments: []DIDDocument{},
		Credentials:  []Credential{},
		MfaConfigs:   []MFAConfig{},
		AuditLogs:    []AuditLogEntry{},
		// this line is used by starport scaffolding # genesis/types/default
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs GenesisState) Validate() error {
	// Check for duplicate DID documents
	didDocumentIdMap := make(map[string]bool)
	for _, doc := range gs.DidDocuments {
		if _, ok := didDocumentIdMap[doc.Id]; ok {
			return fmt.Errorf("duplicate DID document ID: %s", doc.Id)
		}
		didDocumentIdMap[doc.Id] = true
	}

	// Check for duplicate credentials
	credentialIdMap := make(map[string]bool)
	for _, cred := range gs.Credentials {
		if _, ok := credentialIdMap[cred.Id]; ok {
			return fmt.Errorf("duplicate credential ID: %s", cred.Id)
		}
		credentialIdMap[cred.Id] = true
	}

	// Check for duplicate MFA configs
	mfaConfigMap := make(map[string]bool)
	for _, config := range gs.MfaConfigs {
		if _, ok := mfaConfigMap[config.Did]; ok {
			return fmt.Errorf("duplicate MFA config for DID: %s", config.Did)
		}
		mfaConfigMap[config.Did] = true
	}

	// Check for duplicate audit log entries
	auditLogMap := make(map[string]bool)
	for _, log := range gs.AuditLogs {
		if _, ok := auditLogMap[log.Id]; ok {
			return fmt.Errorf("duplicate audit log ID: %s", log.Id)
		}
		auditLogMap[log.Id] = true
	}

	// Validate module parameters
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// this line is used by starport scaffolding # genesis/types/validate

	return nil
}
