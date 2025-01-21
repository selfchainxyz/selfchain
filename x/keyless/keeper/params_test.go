package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "selfchain/testutil/keeper"
	"selfchain/x/keyless/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.KeylessKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)
	require.EqualValues(t, params, k.GetParams(ctx))
}

func TestParams(t *testing.T) {
	k, ctx := testkeeper.KeylessKeeper(t)
	
	// Test default values
	params := k.GetParams(ctx)
	require.Equal(t, uint32(5), params.MaxWalletsPerDid)
	require.Equal(t, uint32(3), params.MaxSharesPerWallet)
	require.Equal(t, uint32(2), params.MinRecoveryThreshold)
	require.Equal(t, uint32(3), params.MaxRecoveryThreshold)
	require.Equal(t, uint32(86400), params.RecoveryWindowSeconds) // 24 hours
	require.Equal(t, uint32(3), params.MaxSigningAttempts)

	// Test setting new values
	newParams := types.NewParams(
		10,  // MaxWalletsPerDid
		5,   // MaxSharesPerWallet
		3,   // MinRecoveryThreshold
		4,   // MaxRecoveryThreshold
		7200, // RecoveryWindowSeconds (2 hours)
		5,    // MaxSigningAttempts
	)
	k.SetParams(ctx, newParams)

	// Verify new values
	params = k.GetParams(ctx)
	require.Equal(t, uint32(10), params.MaxWalletsPerDid)
	require.Equal(t, uint32(5), params.MaxSharesPerWallet)
	require.Equal(t, uint32(3), params.MinRecoveryThreshold)
	require.Equal(t, uint32(4), params.MaxRecoveryThreshold)
	require.Equal(t, uint32(7200), params.RecoveryWindowSeconds)
	require.Equal(t, uint32(5), params.MaxSigningAttempts)
}
