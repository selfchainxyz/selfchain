package test

import (
	test "selfchain/x/migration/tests"
	"selfchain/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestShouldFailIfUpdateConfigWhenSignerIsNotAdmin() {
	ctx := sdk.WrapSDKContext(suite.ctx)

	// Alice tries to add herself to the list of migrators
	_, err := suite.msgServer.UpdateConfig(ctx, &types.MsgUpdateConfig{
		Creator:            test.Alice,
		VestingDuration:    100,
		VestingCliff:       10,
		MinMigrationAmount: 10000,
	})

	suite.Require().ErrorIs(err, types.ErrOnlyAdmin)
}

func (suite *IntegrationTestSuite) TestShouldSetTheNewConfig() {
	ctx := sdk.WrapSDKContext(suite.ctx)

	var config types.Config
	config, _ = suite.app.MigrationKeeper.GetConfig(suite.ctx)
	suite.Require().EqualValues(config.VestingDuration, 2592000)
	suite.Require().EqualValues(config.VestingCliff, 604800)
	suite.Require().EqualValues(config.MinMigrationAmount, 1000000000000000000)

	_, err := suite.msgServer.UpdateConfig(ctx, &types.MsgUpdateConfig{
		Creator:            test.AclAdmin,
		VestingDuration:    100,
		VestingCliff:       10,
		MinMigrationAmount: 10000,
	})

	suite.Require().Nil(err)

	_, exists := suite.app.MigrationKeeper.GetConfig(suite.ctx)
	suite.Require().True(exists)

	config, _ = suite.app.MigrationKeeper.GetConfig(suite.ctx)
	suite.Require().EqualValues(config.VestingDuration, 100)
	suite.Require().EqualValues(config.VestingCliff, 10)
	suite.Require().EqualValues(config.MinMigrationAmount, 10000)
}
