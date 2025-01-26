package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "selfchain/testutil/keeper"
	"selfchain/x/keyless/types"
)

func TestGetParams(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	params := types.DefaultParams()

	k.SetParams(k.Ctx, params)

	require.EqualValues(t, params, k.GetParams(k.Ctx))
}

func TestParams(t *testing.T) {
	k := testkeeper.NewKeylessKeeper(t)
	
	// Test default values
	params := k.GetParams(k.Ctx)

	// Test wallet limits
	require.Equal(t, uint32(5), params.MaxParties)
	require.Equal(t, uint32(3), params.MaxThreshold)
	require.Equal(t, uint32(3), params.MaxSecurityLevel)
	require.Equal(t, uint32(100), params.MaxBatchSize)
	require.Equal(t, uint32(1024), params.MaxMetadataSize)

	// Test recovery settings
	require.Equal(t, uint32(5), params.MaxWalletsPerDid)
	require.Equal(t, uint32(3), params.MaxSharesPerWallet)
	require.Equal(t, uint32(2), params.MinRecoveryThreshold)
	require.Equal(t, uint32(3), params.MaxRecoveryThreshold)
	require.Equal(t, uint32(86400), params.RecoveryWindowSeconds) // 24 hours
	require.Equal(t, uint32(3), params.MaxSigningAttempts)

	// Test setting new values
	newParams := types.NewParams(
		10,  // MaxParties
		5,   // MaxThreshold
		4,   // MaxSecurityLevel
		200, // MaxBatchSize
		2048, // MaxMetadataSize
		10,   // MaxWalletsPerDid
		5,    // MaxSharesPerWallet
		3,    // MinRecoveryThreshold
		4,    // MaxRecoveryThreshold
		7200, // RecoveryWindowSeconds (2 hours)
		5,    // MaxSigningAttempts
	)
	k.SetParams(k.Ctx, newParams)

	// Verify new values
	params = k.GetParams(k.Ctx)

	// Verify wallet limits
	require.Equal(t, uint32(10), params.MaxParties)
	require.Equal(t, uint32(5), params.MaxThreshold)
	require.Equal(t, uint32(4), params.MaxSecurityLevel)
	require.Equal(t, uint32(200), params.MaxBatchSize)
	require.Equal(t, uint32(2048), params.MaxMetadataSize)

	// Verify recovery settings
	require.Equal(t, uint32(10), params.MaxWalletsPerDid)
	require.Equal(t, uint32(5), params.MaxSharesPerWallet)
	require.Equal(t, uint32(3), params.MinRecoveryThreshold)
	require.Equal(t, uint32(4), params.MaxRecoveryThreshold)
	require.Equal(t, uint32(7200), params.RecoveryWindowSeconds)
	require.Equal(t, uint32(5), params.MaxSigningAttempts)
}
