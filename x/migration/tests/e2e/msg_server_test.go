package test

import (
	test "selfchain/x/migration/tests"
	"selfchain/x/migration/types"
	selfvestingTypes "selfchain/x/selfvesting/types"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestShouldFailIfInvalidMigrator() {
	suite.setupSuiteWithBalances()
	ctx := sdk.WrapSDKContext(suite.ctx)

	// Alice who is not a migrator is trying to mint 1M tokens for free
	_, err := suite.msgServer.Migrate(ctx, &types.MsgMigrate{
		Creator:     test.Alice,
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount:      "1000000000000000000000000", // 1 Milion
		Token:       0,
		LogIndex:    0,
	})

	suite.Require().ErrorIs(err, types.ErrUnknownMigrator)
}

func (suite *IntegrationTestSuite) TestShouldMintCorrectRatioForFront() {
	suite.setupSuiteWithBalances()
	ctx := sdk.WrapSDKContext(suite.ctx)
	selfVestingAddr := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, selfvestingTypes.ModuleName).GetAddress()
	aliceAddr, _ := sdk.AccAddressFromBech32(test.Alice)

	balBefore := suite.app.BankKeeper.GetBalance(suite.ctx, selfVestingAddr, types.DENOM)
	balBeneficiaryBefore := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAddr, types.DENOM)

	_, err := suite.msgServer.Migrate(ctx, &types.MsgMigrate{
		Creator:     test.Migrator_1,
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount:      "1000000000000000000000000", // 1 Milion
		Token:       0,
		LogIndex:    0,
	})
	_ = err
	balAfter := suite.app.BankKeeper.GetBalance(suite.ctx, selfVestingAddr, types.DENOM)
	balBeneficiaryAfter := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAddr, types.DENOM)

	suite.Require().EqualValues(sdkmath.NewInt(999999000000), balAfter.Amount.Sub(balBefore.Amount))
	suite.Require().EqualValues(sdkmath.NewInt(1000000), balBeneficiaryAfter.Amount.Sub(balBeneficiaryBefore.Amount))
}

func (suite *IntegrationTestSuite) TestShouldMintCorrectRatioForHotcross() {
	suite.setupSuiteWithBalances()

	// update params
	suite.app.MigrationKeeper.SetParams(suite.ctx, types.Params{
		HotcrossRatio: 50,
	})

	ctx := sdk.WrapSDKContext(suite.ctx)

	selfVestingAddr := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, selfvestingTypes.ModuleName).GetAddress()
	bobAddr, _ := sdk.AccAddressFromBech32(test.Bob)

	balBeneficiaryBefore := suite.app.BankKeeper.GetBalance(suite.ctx, bobAddr, types.DENOM)

	balBefore := suite.app.BankKeeper.GetBalance(suite.ctx, selfVestingAddr, types.DENOM)
	_, err := suite.msgServer.Migrate(ctx, &types.MsgMigrate{
		Creator:     test.Migrator_2,
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Bob,
		Amount:      "1000000000000000000000000", // 1 Milion
		Token:       1,
		LogIndex:    0,
	})
	_ = err

	balAfter := suite.app.BankKeeper.GetBalance(suite.ctx, selfVestingAddr, types.DENOM)
	balBeneficiaryAfter := suite.app.BankKeeper.GetBalance(suite.ctx, bobAddr, types.DENOM)

	suite.Require().EqualValues(sdkmath.NewInt(499999000000), balAfter.Amount.Sub(balBefore.Amount))
	suite.Require().EqualValues(sdkmath.NewInt(1000000), balBeneficiaryAfter.Amount.Sub(balBeneficiaryBefore.Amount))
}

func (suite *IntegrationTestSuite) TestShouldFailWhenHotcrossRationIsZero() {
	suite.setupSuiteWithBalances()
	ctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.msgServer.Migrate(ctx, &types.MsgMigrate{
		Creator:     test.Migrator_2,
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Bob,
		Amount:      "1000000000000000000000000", // 1 Milion
		Token:       1,
		LogIndex:    0,
	})

	suite.Require().ErrorIs(err, types.ErrHotcrossRatioZero)
}

func (suite *IntegrationTestSuite) TestShouldFailIfMigrationProcessed() {
	suite.setupSuiteWithBalances()
	ctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.msgServer.Migrate(ctx, &types.MsgMigrate{
		Creator:     test.Migrator_1,
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount:      "1000000000000000000000000", // 1 Milion
		Token:       0,
		LogIndex:    0,
	})
	suite.Require().Nil(err)

	_, err2 := suite.msgServer.Migrate(ctx, &types.MsgMigrate{
		Creator:     test.Migrator_1,
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount:      "1000000000000000000000000", // 1 Milion
		Token:       0,
		LogIndex:    0,
	})

	suite.Require().ErrorIs(err2, types.ErrMigrationProcessed)
}

func (suite *IntegrationTestSuite) TestShouldInstanlyReleaseFullAmount() {
	suite.setupSuiteWithBalances()
	ctx := sdk.WrapSDKContext(suite.ctx)
	selfVestingAddr := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, selfvestingTypes.ModuleName).GetAddress()
	aliceAddr, _ := sdk.AccAddressFromBech32(test.Alice)

	balBefore := suite.app.BankKeeper.GetBalance(suite.ctx, selfVestingAddr, types.DENOM)
	balBeneficiaryBefore := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAddr, types.DENOM)

	_, err := suite.msgServer.Migrate(ctx, &types.MsgMigrate{
		Creator:     test.Migrator_1,
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount:      "1000000000000000000",
		Token:       0,
		LogIndex:    0,
	})
	_ = err

	balAfter := suite.app.BankKeeper.GetBalance(suite.ctx, selfVestingAddr, types.DENOM)
	balBeneficiaryAfter := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAddr, types.DENOM)

	suite.Require().EqualValues(sdkmath.NewInt(0), balAfter.Amount.Sub(balBefore.Amount))
	suite.Require().EqualValues(sdkmath.NewInt(1000000), balBeneficiaryAfter.Amount.Sub(balBeneficiaryBefore.Amount))
}
