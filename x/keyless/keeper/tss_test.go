package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"selfchain/x/keyless/testutil/keeper"
	"selfchain/x/keyless/types"
)

func TestKeyShare(t *testing.T) {
	k, ctx := keeper.KeylessKeeper(t)
	
	// Test key share management
	walletID := "test_wallet"
	partyID := "test_party"
	shareData := []byte("test_share_data")

	// Save key share
	err := k.SavePartyData(ctx, walletID, &types.PartyData{
		PartyId:    partyID,
		PartyShare: shareData,
		ChainId:    "test-chain",
		Status:     "active",
	})
	require.NoError(t, err)

	// Get key share
	retrievedData, err := k.GetPartyData(ctx, walletID)
	require.NoError(t, err)
	require.Equal(t, shareData, retrievedData.PartyShare)
	require.Equal(t, "active", retrievedData.Status)
}

func TestWalletKeyGeneration(t *testing.T) {
	k, ctx := keeper.KeylessKeeper(t)

	tests := []struct {
		name          string
		walletAddress string
		chainID       string
		creator       string
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "valid wallet",
			walletAddress: "cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq",
			chainID:      "test-1",
			creator:      "cosmos1creator",
			wantErr:      false,
		},
		{
			name:          "invalid wallet address",
			walletAddress: "invalid",
			chainID:      "test-1",
			creator:      "cosmos1creator",
			wantErr:      true,
			errMsg:       "invalid bech32 address",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create and save wallet
			wallet := &types.Wallet{
				WalletAddress: tc.walletAddress,
				ChainId:      tc.chainID,
				Creator:      tc.creator,
				Status:       types.WalletStatus_WALLET_STATUS_ACTIVE,
			}

			// SetWallet is a void function, so we don't need to capture its return value
			k.SetWallet(ctx, wallet)

			// Verify wallet was saved
			savedWallet, found := k.GetWallet(ctx, tc.walletAddress)
			if tc.wantErr {
				require.False(t, found)
				return
			}
			require.True(t, found)
			require.Equal(t, wallet.WalletAddress, savedWallet.WalletAddress)
			require.Equal(t, wallet.ChainId, savedWallet.ChainId)
			require.Equal(t, wallet.Creator, savedWallet.Creator)
			require.Equal(t, wallet.Status, savedWallet.Status)
		})
	}
}
