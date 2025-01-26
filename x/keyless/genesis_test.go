package keyless_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"selfchain/x/keyless"
	"selfchain/x/keyless/testutil"
	"selfchain/x/keyless/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		// Add test wallets, party data, etc. if needed
		Wallets: []types.Wallet{
			{
				WalletAddress: "test_wallet_1",
				ChainId:      "test_chain_1",
				Owner:        "test_owner_1",
				Status:       types.WalletStatus_ACTIVE,
			},
		},
		PartyData: []types.PartyData{
			{
				PartyId:          "test_party_1",
				PublicKey:        []byte("test_public_key_1"),
				PartyShare:       []byte("test_party_share_1"),
				VerificationData: []byte("test_verification_data_1"),
				ChainId:         "test_chain_1",
				Status:          types.PartyStatus_ACTIVE,
			},
		},
	}

	k, ctx := testutil.NewTestKeeper(t)
	keyless.InitGenesis(ctx, k, genesisState)
	got := keyless.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	// Verify params
	require.Equal(t, genesisState.Params, got.Params)

	// Verify wallets
	require.Equal(t, len(genesisState.Wallets), len(got.Wallets))
	for i := range genesisState.Wallets {
		require.Equal(t, genesisState.Wallets[i].WalletAddress, got.Wallets[i].WalletAddress)
		require.Equal(t, genesisState.Wallets[i].ChainId, got.Wallets[i].ChainId)
		require.Equal(t, genesisState.Wallets[i].Owner, got.Wallets[i].Owner)
		require.Equal(t, genesisState.Wallets[i].Status, got.Wallets[i].Status)
	}

	// Verify party data
	require.Equal(t, len(genesisState.PartyData), len(got.PartyData))
	for i := range genesisState.PartyData {
		require.Equal(t, genesisState.PartyData[i].PartyId, got.PartyData[i].PartyId)
		require.Equal(t, genesisState.PartyData[i].PublicKey, got.PartyData[i].PublicKey)
		require.Equal(t, genesisState.PartyData[i].PartyShare, got.PartyData[i].PartyShare)
		require.Equal(t, genesisState.PartyData[i].VerificationData, got.PartyData[i].VerificationData)
		require.Equal(t, genesisState.PartyData[i].ChainId, got.PartyData[i].ChainId)
		require.Equal(t, genesisState.PartyData[i].Status, got.PartyData[i].Status)
	}
}
