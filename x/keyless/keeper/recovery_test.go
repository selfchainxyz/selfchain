package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "selfchain/testutil/keeper"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
)

func TestRecoverWallet(t *testing.T) {
	k := keepertest.NewKeylessKeeper(t)
	srv := keeper.NewMsgServerImpl(k.Keeper)
	wctx := sdk.WrapSDKContext(k.Ctx)

	// Clear store before test
	k.ClearStore()

	// Create test wallet first
	walletAddr := "cosmos1w3jhxap3ta047h6lta047h6lta047h6lx84s66"
	creator := "cosmos1w3jhxap3ta047h6lta047h6lta047h6lx84s66"
	grantee := "cosmos1w3jhxapjta047h6lta047h6lta047h6lwuy8a3"

	msg := &types.MsgCreateWallet{
		Creator:       creator,
		PubKey:        "pubkey1",
		WalletAddress: walletAddr,
		ChainId:       "test-1",
	}
	_, err := srv.CreateWallet(wctx, msg)
	require.NoError(t, err)

	// Set wallet status to active
	err = k.SetWalletStatus(k.Ctx, walletAddr, types.WalletStatus_WALLET_STATUS_ACTIVE)
	require.NoError(t, err)

	// Grant recovery permission
	expiresAt := time.Now().Add(24 * time.Hour)
	permission := &types.Permission{
		WalletAddress: walletAddr,
		Grantee:      grantee,
		Permissions:  []string{types.WalletPermission_WALLET_PERMISSION_RECOVER.String()},
		ExpiresAt:    &expiresAt,
	}

	err = k.GrantPermission(k.Ctx, permission)
	require.NoError(t, err)

	// Create recovery session
	recoveryInfo, err := k.CreateRecoverySession(k.Ctx, walletAddr)
	require.NoError(t, err)
	require.NotNil(t, recoveryInfo)
	require.Equal(t, walletAddr, recoveryInfo.Did)
	require.Equal(t, types.RecoveryStatus_RECOVERY_STATUS_PENDING, recoveryInfo.Status)

	// Test recovery with invalid token - should fail
	recoveryMsg := &types.MsgRecoverWallet{
		Creator:       grantee,
		WalletAddress: walletAddr,
		NewPubKey:    "newpubkey1",
		RecoveryProof: "invalid_token",
	}

	_, err = srv.RecoverWallet(wctx, recoveryMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid recovery token")

	// Test recovery with valid token - should succeed
	recoveryMsg.RecoveryProof = recoveryInfo.RecoveryToken
	_, err = srv.RecoverWallet(wctx, recoveryMsg)
	require.NoError(t, err)

	// Verify wallet status is updated
	wallet, found := k.GetWallet(k.Ctx, walletAddr)
	require.True(t, found)
	require.Equal(t, types.WalletStatus_WALLET_STATUS_ACTIVE, wallet.Status)
}

func TestCreateRecoverySession(t *testing.T) {
	k := keepertest.NewKeylessKeeper(t)
	srv := keeper.NewMsgServerImpl(k.Keeper)
	wctx := sdk.WrapSDKContext(k.Ctx)

	// Clear store before test
	k.ClearStore()

	// Create test wallet first
	walletAddr := "cosmos1w3jhxap3ta047h6lta047h6lta047h6lx84s66"
	creator := "cosmos1w3jhxap3ta047h6lta047h6lta047h6lx84s66"

	msg := &types.MsgCreateWallet{
		Creator:       creator,
		PubKey:        "pubkey1",
		WalletAddress: walletAddr,
		ChainId:       "test-1",
	}
	_, err := srv.CreateWallet(wctx, msg)
	require.NoError(t, err)

	// Set wallet status to active
	err = k.SetWalletStatus(k.Ctx, walletAddr, types.WalletStatus_WALLET_STATUS_ACTIVE)
	require.NoError(t, err)

	// Test creating recovery session
	recoveryInfo, err := k.CreateRecoverySession(k.Ctx, walletAddr)
	require.NoError(t, err)
	require.NotNil(t, recoveryInfo)
	require.Equal(t, walletAddr, recoveryInfo.Did)
	require.Equal(t, types.RecoveryStatus_RECOVERY_STATUS_PENDING, recoveryInfo.Status)

	// Test creating duplicate recovery session
	_, err = k.CreateRecoverySession(k.Ctx, walletAddr)
	require.Error(t, err)
	require.Contains(t, err.Error(), "recovery session already exists")

	// Test creating recovery session with invalid wallet address
	_, err = k.CreateRecoverySession(k.Ctx, "invalid")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid bech32 address")

	// Test creating recovery session with non-existent wallet
	_, err = k.CreateRecoverySession(k.Ctx, "cosmos1w3jhxapnta047h6lta047h6lta047h6l34t280")
	require.Error(t, err)
	require.Contains(t, err.Error(), "wallet not found")
}
