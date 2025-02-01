package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
			walletAddress: "cosmos1w3jhxap3ta047h6lta047h6lta047h6lx84s66",
			chainID:      "test-1",
			creator:      "cosmos1w3jhxap3ta047h6lta047h6lx84s66",
			wantErr:      false,
		},
		{
			name:          "invalid wallet address",
			walletAddress: "invalid",
			chainID:      "test-1",
			creator:      "cosmos1w3jhxap3ta047h6lta047h6lx84s66",
			wantErr:      true,
			errMsg:       "decoding bech32 failed: invalid bech32 string length 7",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create and save wallet
			now := time.Now().UTC()
			wallet := &types.Wallet{
				Id:            tc.walletAddress,
				Creator:       tc.creator,
				PublicKey:     "test_pubkey",
				WalletAddress: tc.walletAddress,
				ChainId:      tc.chainID,
				Status:       types.WalletStatus_WALLET_STATUS_ACTIVE,
				KeyVersion:   1,
				CreatedAt:    &now,
				UpdatedAt:    &now,
				LastUsed:     &now,
				UsageCount:   0,
			}

			// Validate wallet address before saving
			_, err := sdk.AccAddressFromBech32(tc.walletAddress)
			if tc.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMsg)
				return
			}
			require.NoError(t, err)

			// Save wallet
			k.SaveWallet(ctx, wallet)

			// Verify wallet was saved
			savedWallet, found := k.GetWallet(ctx, tc.walletAddress)
			require.True(t, found)
			require.Equal(t, wallet.Id, savedWallet.Id)
			require.Equal(t, wallet.Creator, savedWallet.Creator)
			require.Equal(t, wallet.PublicKey, savedWallet.PublicKey)
			require.Equal(t, wallet.WalletAddress, savedWallet.WalletAddress)
			require.Equal(t, wallet.ChainId, savedWallet.ChainId)
			require.Equal(t, wallet.Status, savedWallet.Status)
			require.Equal(t, wallet.KeyVersion, savedWallet.KeyVersion)
		})
	}
}
