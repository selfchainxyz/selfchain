package test

import (
	test "frontier/x/migration/tests"
	"frontier/x/migration/types"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestShouldFailIfInvalidMigrator() {
	suite.setupSuiteWithBalances()
	ctx := sdk.WrapSDKContext(suite.ctx)

	// Alice who is not a migrator is trying to mint 1M tokens for free
	_, err := suite.msgServer.Migrate(ctx, &types.MsgMigrate{
		Creator: test.Alice,
		TxHash:  "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount: "1000000000000000000000000", // 1 Milion
		Token: 0,
	})

	suite.Require().ErrorIs(err, types.ErrUnknownMigrator)
}

func (suite *IntegrationTestSuite) TestShouldMintAmount() {
	suite.setupSuiteWithBalances()
	ctx := sdk.WrapSDKContext(suite.ctx)
	aliceAddr, _ := sdk.AccAddressFromBech32(test.Alice)

	balBefore := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAddr, types.DENOM)
	_, err := suite.msgServer.Migrate(ctx, &types.MsgMigrate{
		Creator: test.Migrator_1,
		TxHash:  "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
		EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: test.Alice,
		Amount: "1000000000000000000000000", // 1 Milion
		Token: 0,
	})
	_ = err
	balAfter := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAddr, types.DENOM)

	suite.Require().EqualValues(sdkmath.NewInt(100000000000), balAfter.Amount.Sub(balBefore.Amount))
}
