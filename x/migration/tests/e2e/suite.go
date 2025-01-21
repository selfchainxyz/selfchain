package test

import (
	"context"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"selfchain/app"
	"selfchain/testutil"
	"selfchain/x/migration/keeper"
	"selfchain/x/migration/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	app         *app.App
	ctx         sdk.Context
	msgServer   types.MsgServer
	queryClient types.QueryClient
	addrs       []sdk.AccAddress
}

func (suite *IntegrationTestSuite) SetupTest() {
	app := testutil.Setup(suite.T(), false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	suite.app = app
	suite.ctx = ctx
	suite.msgServer = keeper.NewMsgServerImpl(app.MigrationKeeper)

	// Create a connection that bypasses gRPC
	queryHelper := &QueryHelper{keeper: app.MigrationKeeper}
	suite.queryClient = types.NewQueryClient(queryHelper)
	suite.setupSuiteWithBalances()
}

func (suite *IntegrationTestSuite) setupSuiteWithBalances() {
	// Generate test addresses
	for i := 0; i < 3; i++ {
		priv := secp256k1.GenPrivKey()
		addr := sdk.AccAddress(priv.PubKey().Address())
		acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr)
		suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
		suite.addrs = append(suite.addrs, addr)

		// Fund the account
		coins := sdk.NewCoins(
			sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000)),
		)
		err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, coins)
		suite.Require().NoError(err)
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, addr, coins)
		suite.Require().NoError(err)
	}
}

// QueryHelper implements the grpc.ClientConnInterface
type QueryHelper struct {
	keeper keeper.Keeper
}

func (q *QueryHelper) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	// Here you would implement the specific query handling based on the method
	// For example:
	switch method {
	// Add cases for each query method you need to support
	default:
		return nil
	}
}

func (q *QueryHelper) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, nil
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
