package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"selfchain/testutil/keeper"
	"selfchain/x/keyless/types"
)

func TestKeyShare(t *testing.T) {
	k := keeper.NewKeylessKeeper(t)
	
	// Test key share management
	walletID := "test_wallet"
	partyID := "test_party"
	shareData := []byte("test_share_data")

	// Save key share
	err := k.SavePartyData(k.Ctx, walletID, &types.PartyData{
		PartyId:    partyID,
		PartyShare: shareData,
		ChainId:    "test-chain",
		Status:     "active",
	})
	require.NoError(t, err)

	// Get key share
	retrievedData, err := k.GetPartyData(k.Ctx, walletID)
	require.NoError(t, err)
	require.Equal(t, shareData, retrievedData.PartyShare)
	require.Equal(t, "active", retrievedData.Status)
}

func TestWalletKeyGeneration(t *testing.T) {
	k := keeper.NewKeylessKeeper(t)

	// Test cases
	testCases := []struct {
		name          string
		walletAddress string
		chainID       string
		creator       string
		wantErr      bool
	}{
		{
			name:          "valid key generation",
			walletAddress: "test_wallet_1",
			chainID:       "test-chain-1",
			creator:       "test_creator_1",
			wantErr:      false,
		},
		{
			name:          "empty wallet address",
			walletAddress: "",
			chainID:       "test-chain-2",
			creator:       "test_creator_2",
			wantErr:      true,
		},
		{
			name:          "empty chain ID",
			walletAddress: "test_wallet_3",
			chainID:       "",
			creator:       "test_creator_3",
			wantErr:      true,
		},
		{
			name:          "empty creator",
			walletAddress: "test_wallet_4",
			chainID:       "test-chain-4",
			creator:       "",
			wantErr:      true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create wallet
			wallet := &types.Wallet{
				WalletAddress: tc.walletAddress,
				ChainId:      tc.chainID,
				Creator:      tc.creator,
				Status:       types.WalletStatus_WALLET_STATUS_ACTIVE,
			}

			err := k.SaveWallet(k.Ctx, wallet)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Verify wallet was saved
			savedWallet, err := k.GetWallet(k.Ctx, tc.walletAddress)
			require.NoError(t, err)
			require.Equal(t, wallet.WalletAddress, savedWallet.WalletAddress)
			require.Equal(t, wallet.ChainId, savedWallet.ChainId)
			require.Equal(t, wallet.Creator, savedWallet.Creator)
			require.Equal(t, wallet.Status, savedWallet.Status)
		})
	}
}
