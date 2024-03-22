// This is used for integration tests
package test

import (
	"selfchain/app"
	"selfchain/x/selfvesting"
	"selfchain/x/selfvesting/keeper"
	test "selfchain/x/selfvesting/tests"
	"selfchain/x/selfvesting/types"
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
)

var (
	selfvestingModuleAddress string
)

func TestSelfvestingTestSuite(t *testing.T) {
	test.InitSDKConfig()
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
	app := app.Setup(suite.T(), false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	selfvestingModuleAddress = app.AccountKeeper.GetModuleAddress(types.ModuleName).String()

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.SelfvestingKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.msgServer = keeper.NewMsgServerImpl(app.SelfvestingKeeper)
	suite.ctx = ctx
	suite.queryClient = queryClient

	selfvesting.InitGenesis(ctx, app.SelfvestingKeeper, getModuleGenesis())
}

func getModuleGenesis() types.GenesisState {
	genesis := *types.DefaultGenesis()

	genesis.VestingPositionsList = []types.VestingPositions{}

	return genesis
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
		makeBalance(test.Alice, 1000000000, "uslf"),
		makeBalance(test.Bob, 1000000000, "uslf"),
		makeBalance(test.Carol, 1000000000, "uslf"),
	}
	supply := banktypes.Supply{
		Total: addAll(coins),
	}

	state := banktypes.NewGenesisState(
		banktypes.DefaultParams(),
		coins,
		supply.Total,
		[]banktypes.Metadata{},
   	[]banktypes.SendEnabled{})

	return state
}

func (suite *IntegrationTestSuite) setupSuiteWithBalances() {
	suite.app.BankKeeper.InitGenesis(suite.ctx, getBankGenesis())
}
