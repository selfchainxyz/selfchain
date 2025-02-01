package testutil_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"selfchain/x/keyless/testutil"
	"selfchain/x/keyless/types"
)

// setupBenchmarkSuite creates a new test suite instance for benchmarking
func setupBenchmarkSuite(t testing.TB) *testutil.IntegrationTestSuite {
	suite := new(testutil.IntegrationTestSuite)
	if tt, ok := t.(*testing.T); ok {
		suite.SetT(tt)
	}
	suite.SetupTest()
	return suite
}

// setupKeyShares creates mock key shares for testing
func setupKeyShares(suite *testutil.IntegrationTestSuite, creator, walletAddr string) error {
	// Create a simple mock key share data
	mockShare := map[string]interface{}{
		"key_share": "mock_share_data",
		"version":   1,
	}

	// Marshal share data
	shareData, err := json.Marshal(mockShare)
	if err != nil {
		return err
	}

	// Store key shares for both creator and wallet
	ctx := suite.Ctx()
	if err := suite.Keeper().StoreKeyShare(ctx, creator, shareData); err != nil {
		return err
	}
	if err := suite.Keeper().StoreKeyShare(ctx, walletAddr, shareData); err != nil {
		return err
	}

	return nil
}

// BenchmarkWalletCreation measures the performance of wallet creation
func BenchmarkWalletCreation(b *testing.B) {
	suite := setupBenchmarkSuite(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		creator := fmt.Sprintf("cosmos1creator%d", i)
		pubKey := fmt.Sprintf("test_pubkey%d", i)
		walletAddr := fmt.Sprintf("cosmos1wallet%d", i)
		chainID := "test-chain-1"

		msg := types.NewMsgCreateWallet(
			creator,
			pubKey,
			walletAddr,
			chainID,
		)

		_, err := suite.MsgServer().CreateWallet(sdk.WrapSDKContext(suite.Ctx()), msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSignTransaction measures the performance of transaction signing
func BenchmarkSignTransaction(b *testing.B) {
	suite := setupBenchmarkSuite(b)

	// Create a test wallet first
	creator := "cosmos1creator_sign"
	pubKey := "test_pubkey_sign"
	walletAddr := "cosmos1wallet_sign"
	chainID := "test-chain-1"

	createMsg := types.NewMsgCreateWallet(
		creator,
		pubKey,
		walletAddr,
		chainID,
	)

	_, err := suite.MsgServer().CreateWallet(sdk.WrapSDKContext(suite.Ctx()), createMsg)
	if err != nil {
		b.Fatal(err)
	}

	// Set up key shares
	if err := setupKeyShares(suite, creator, walletAddr); err != nil {
		b.Fatal(err)
	}

	// Prepare signing message
	msg := []byte("test message")
	unsignedTx := hex.EncodeToString(msg) // Convert to hex string for unsigned tx

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		signMsg := types.NewMsgSignTransaction(
			creator,
			walletAddr,
			unsignedTx,
		)

		_, err := suite.MsgServer().SignTransaction(sdk.WrapSDKContext(suite.Ctx()), signMsg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestStressParallelSigning tests parallel signing operations
func TestStressParallelSigning(t *testing.T) {
	suite := setupBenchmarkSuite(t)

	// Create test wallets
	numWallets := 10
	wallets := make([]string, numWallets)
	creators := make([]string, numWallets)

	for i := 0; i < numWallets; i++ {
		creator := fmt.Sprintf("cosmos1creator_stress%d", i)
		pubKey := fmt.Sprintf("test_pubkey_stress%d", i)
		walletAddr := fmt.Sprintf("cosmos1wallet_stress%d", i)
		chainID := "test-chain-1"

		createMsg := types.NewMsgCreateWallet(
			creator,
			pubKey,
			walletAddr,
			chainID,
		)

		_, err := suite.MsgServer().CreateWallet(sdk.WrapSDKContext(suite.Ctx()), createMsg)
		require.NoError(t, err)

		// Set up key shares
		err = setupKeyShares(suite, creator, walletAddr)
		require.NoError(t, err)

		wallets[i] = walletAddr
		creators[i] = creator
	}

	// Run parallel signing operations
	t.Run("ParallelSigning", func(t *testing.T) {
		for i := 0; i < numWallets; i++ {
			walletAddr := wallets[i]
			creator := creators[i]
			t.Run(fmt.Sprintf("Wallet%d", i), func(t *testing.T) {
				t.Parallel()
				msg := []byte("test message")
				unsignedTx := hex.EncodeToString(msg)

				signMsg := types.NewMsgSignTransaction(
					creator,
					walletAddr,
					unsignedTx,
				)

				_, err := suite.MsgServer().SignTransaction(sdk.WrapSDKContext(suite.Ctx()), signMsg)
				require.NoError(t, err)
			})
		}
	})
}

// TestGasConsumption measures gas consumption for various operations
func TestGasConsumption(t *testing.T) {
	suite := setupBenchmarkSuite(t)

	initialGas := suite.Ctx().GasMeter().GasConsumed()

	// Test wallet creation gas
	t.Run("WalletCreationGas", func(t *testing.T) {
		creator := "cosmos1creator_gas1"
		pubKey := "test_pubkey_gas1"
		walletAddr := "cosmos1wallet_gas1"
		chainID := "test-chain-1"

		createMsg := types.NewMsgCreateWallet(
			creator,
			pubKey,
			walletAddr,
			chainID,
		)

		_, err := suite.MsgServer().CreateWallet(sdk.WrapSDKContext(suite.Ctx()), createMsg)
		require.NoError(t, err)

		gasUsed := suite.Ctx().GasMeter().GasConsumed() - initialGas
		t.Logf("Gas used for wallet creation: %d", gasUsed)
		require.Less(t, gasUsed, sdk.Gas(1000000))
	})

	// Reset gas meter
	suite.SetCtx(suite.Ctx().WithGasMeter(sdk.NewGasMeter(1000000000)))
	initialGas = suite.Ctx().GasMeter().GasConsumed()

	// Test signing gas
	t.Run("SigningGas", func(t *testing.T) {
		creator := "cosmos1creator_gas2"
		pubKey := "test_pubkey_gas2"
		walletAddr := "cosmos1wallet_gas2"
		chainID := "test-chain-1"

		createMsg := types.NewMsgCreateWallet(
			creator,
			pubKey,
			walletAddr,
			chainID,
		)

		_, err := suite.MsgServer().CreateWallet(sdk.WrapSDKContext(suite.Ctx()), createMsg)
		require.NoError(t, err)

		// Set up key shares
		err = setupKeyShares(suite, creator, walletAddr)
		require.NoError(t, err)

		msg := []byte("test message")
		unsignedTx := hex.EncodeToString(msg)

		signMsg := types.NewMsgSignTransaction(
			creator,
			walletAddr,
			unsignedTx,
		)

		_, err = suite.MsgServer().SignTransaction(sdk.WrapSDKContext(suite.Ctx()), signMsg)
		require.NoError(t, err)

		gasUsed := suite.Ctx().GasMeter().GasConsumed() - initialGas
		t.Logf("Gas used for signing: %d", gasUsed)
		require.Less(t, gasUsed, sdk.Gas(500000))
	})
}
