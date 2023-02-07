package test

import (
	test "frontier/x/migration/tests"
	"frontier/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestShouldIfSignerIsNotAdmin() {
	suite.setupSuiteWithBalances()
	ctx := sdk.WrapSDKContext(suite.ctx)

	// Alice tries to add herself to the list of migrators
	_, err := suite.msgServer.AddMigrator(ctx, &types.MsgAddMigrator{
		Creator:  test.Alice,
		Migrator: test.Alice,
	})

	suite.Require().ErrorIs(err, types.ErrOnlyAdmin)
}
