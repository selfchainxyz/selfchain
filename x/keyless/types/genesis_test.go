package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"selfchain/x/keyless/types"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params: types.Params{
					MaxParties:           5,
					MaxThreshold:         3,
					MaxSecurityLevel:     3,
					MaxBatchSize:         100,
					MaxMetadataSize:      1024,
					MaxWalletsPerDid:     5,
					MaxSharesPerWallet:   3,
					MinRecoveryThreshold: 2,
					MaxRecoveryThreshold: 3,
					RecoveryWindowSeconds: 86400,
					MaxSigningAttempts:    3,
				},
			},
			valid: true,
		},
		{
			desc: "invalid params - zero MaxParties",
			genState: &types.GenesisState{
				Params: types.Params{
					MaxParties:           0, // Invalid
					MaxThreshold:         3,
					MaxSecurityLevel:     3,
					MaxBatchSize:         100,
					MaxMetadataSize:      1024,
					MaxWalletsPerDid:     5,
					MaxSharesPerWallet:   3,
					MinRecoveryThreshold: 2,
					MaxRecoveryThreshold: 3,
					RecoveryWindowSeconds: 86400,
					MaxSigningAttempts:    3,
				},
			},
			valid: false,
		},
		{
			desc: "invalid params - MaxThreshold > MaxParties",
			genState: &types.GenesisState{
				Params: types.Params{
					MaxParties:           3,
					MaxThreshold:         5, // Invalid
					MaxSecurityLevel:     3,
					MaxBatchSize:         100,
					MaxMetadataSize:      1024,
					MaxWalletsPerDid:     5,
					MaxSharesPerWallet:   3,
					MinRecoveryThreshold: 2,
					MaxRecoveryThreshold: 3,
					RecoveryWindowSeconds: 86400,
					MaxSigningAttempts:    3,
				},
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
