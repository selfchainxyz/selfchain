package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"selfchain/testutil/keeper"
	"selfchain/x/keyless/types"
)

func TestRecoverWallet(t *testing.T) {
	k := keeper.NewKeylessKeeper(t)
	
	// Create a wallet first
	walletID := "test_wallet"
	creator := "test_creator"
	wallet := &types.Wallet{
		WalletAddress: walletID,
		ChainId:      "test-chain",
		Creator:      creator,
		Status:       types.WalletStatus_WALLET_STATUS_INACTIVE,
	}
	err := k.SaveWallet(k.Ctx, wallet)
	require.NoError(t, err)

	// Test recovery process
	newOwner := "new_owner"
	newPubKey := "new_pubkey"
	recoveryMsg := &types.MsgRecoverWallet{
		Creator:       newOwner,
		WalletAddress: walletID,
		NewPubKey:    newPubKey,
		RecoveryProof: "test_proof",
	}

	err = k.RecoverWallet(k.Ctx, recoveryMsg)
	require.NoError(t, err)

	// Verify wallet is recovered with new owner and key
	recoveredWallet, err := k.GetWallet(k.Ctx, walletID)
	require.NoError(t, err)
	require.Equal(t, types.WalletStatus_WALLET_STATUS_ACTIVE, recoveredWallet.Status)
	require.Equal(t, newOwner, recoveredWallet.Creator)
	require.Equal(t, newPubKey, recoveredWallet.PublicKey)
}

func TestCreateRecoverySession(t *testing.T) {
	k := keeper.NewKeylessKeeper(t)
	
	// Create a wallet first
	walletID := "test_wallet"
	creator := "test_creator"
	wallet := &types.Wallet{
		WalletAddress: walletID,
		ChainId:      "test-chain",
		Creator:      creator,
		Status:       types.WalletStatus_WALLET_STATUS_INACTIVE,
	}
	err := k.SaveWallet(k.Ctx, wallet)
	require.NoError(t, err)

	// Test recovery session creation
	err = k.CreateRecoverySession(k.Ctx, creator, walletID)
	require.NoError(t, err)
}
