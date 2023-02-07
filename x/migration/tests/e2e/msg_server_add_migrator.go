package test

import (
	test "frontier/x/migration/tests"
	"frontier/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestShouldIfSignerIsNotAdmin() {
	ctx := sdk.WrapSDKContext(suite.ctx)

	// Alice tries to add herself to the list of migrators
	_, err := suite.msgServer.AddMigrator(ctx, &types.MsgAddMigrator{
		Creator:  test.Alice,
		Migrator: test.Alice,
	})

	suite.Require().ErrorIs(err, types.ErrOnlyAdmin)
}

func (suite *IntegrationTestSuite) TestShouldSetTheNewMigrator() {
	ctx := sdk.WrapSDKContext(suite.ctx)

	// Alice tries to add herself to the list of migrators
	_, err := suite.msgServer.AddMigrator(ctx, &types.MsgAddMigrator{
		Creator:  test.AclAdmin,
		Migrator: test.Alice,
	})

	suite.Require().Nil(err)

	_, exists := suite.app.MigrationKeeper.GetMigrator(suite.ctx, test.Alice);
	suite.Require().True(exists);
}
