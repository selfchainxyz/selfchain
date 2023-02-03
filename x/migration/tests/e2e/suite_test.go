// This is used for integration tests
package test

import (
	"frontier/app"
	"frontier/x/migration/keeper"
	test "frontier/x/migration/tests"
	"frontier/x/migration/types"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

var (
	migrationModuleAddress string
)

func TestMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

type IntegrationTestSuite struct {
	suite.Suite

	app         *app.App
	msgServer   types.MsgServer
	ctx         sdk.Context
	queryClient types.QueryClient
}

func (suite *IntegrationTestSuite) SetupTest() {
	test.InitSDKConfig()
	app := app.Setup(suite.T(), false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	migrationModuleAddress = app.AccountKeeper.GetModuleAddress(types.ModuleName).String()

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.MigrationKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.msgServer = keeper.NewMsgServerImpl(app.MigrationKeeper)
	suite.ctx = ctx
	suite.queryClient = queryClient
}

func makeBalance(address string, balance int64, denom string) banktypes.Balance {
	return banktypes.Balance{
		Address: address,
		Coins: sdk.Coins{
			sdk.Coin{
				Denom:  denom,
				Amount: sdk.NewInt(balance),
			},
		},
	}
}

func addAll(balances []banktypes.Balance) sdk.Coins {
	total := sdk.NewCoins()
	for _, balance := range balances {
		total = total.Add(balance.Coins...)
	}
	return total
}

func getBankGenesis() *banktypes.GenesisState {
	coins := []banktypes.Balance{
		makeBalance(test.Alice, 1000000000, "ufront"),
		makeBalance(test.Bob, 1000000000, "ufront"),
		makeBalance(test.Carol, 1000000000, "ufront"),
		makeBalance(test.Migrator_1, 1000000000, "ufront"),
		makeBalance(test.Migrator_2, 1000000000, "ufront"),
	}
	supply := banktypes.Supply{
		Total: addAll(coins),
	}

	state := banktypes.NewGenesisState(
		banktypes.DefaultParams(),
		coins,
		supply.Total,
		[]banktypes.Metadata{})

	return state
}

func (suite *IntegrationTestSuite) setupSuiteWithBalances() {
	suite.app.BankKeeper.InitGenesis(suite.ctx, getBankGenesis())
}
