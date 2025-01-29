package testutil

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
)

type IntegrationTestSuite struct {
	suite.Suite
	keeper    *keeper.Keeper
	ctx       sdk.Context
	msgServer types.MsgServer
}

func (s *IntegrationTestSuite) SetupSuite() {
	// Any one-time setup for the entire suite
}

func (s *IntegrationTestSuite) SetupTest() {
	// Create a fresh keeper and context for each test
	k, ctx := NewTestKeeper(s.T())
	s.keeper = k
	s.ctx = ctx
	s.msgServer = keeper.NewMsgServerImpl(k)
}

func (s *IntegrationTestSuite) TestA_CreateWallet() {
	msg := &types.MsgCreateWallet{
		Creator:       "cosmos1creator",
		PubKey:       "testPubKey",
		WalletAddress: "cosmos1wallet1",
		ChainId:      "test-chain",
	}

	resp, err := s.msgServer.CreateWallet(sdk.WrapSDKContext(s.ctx), msg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(msg.WalletAddress, resp.WalletAddress)

	// Verify wallet was created
	wallet, err := s.keeper.GetWallet(s.ctx, msg.WalletAddress)
	s.Require().NoError(err)
	s.Require().NotNil(wallet)
	s.Require().Equal(msg.Creator, wallet.Creator)
	s.Require().Equal(msg.PubKey, wallet.PublicKey)
	s.Require().Equal(msg.ChainId, wallet.ChainId)
	s.Require().Equal(types.WalletStatus_WALLET_STATUS_ACTIVE, wallet.Status)
}

func (s *IntegrationTestSuite) TestB_SignTransaction() {
	// Create a wallet first
	createMsg := &types.MsgCreateWallet{
		Creator:       "cosmos1creator",
		PubKey:       "testPubKey",
		WalletAddress: "cosmos1wallet2",
		ChainId:      "test-chain",
	}

	_, err := s.msgServer.CreateWallet(sdk.WrapSDKContext(s.ctx), createMsg)
	s.Require().NoError(err)

	// Test signing
	signMsg := &types.MsgSignTransaction{
		Creator:       createMsg.Creator,
		WalletAddress: createMsg.WalletAddress,
		UnsignedTx:   "test_tx",
	}

	resp, err := s.msgServer.SignTransaction(sdk.WrapSDKContext(s.ctx), signMsg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().NotEmpty(resp.SignedTx)
}

func (s *IntegrationTestSuite) TestC_RecoverWallet() {
	// Create a wallet first
	createMsg := &types.MsgCreateWallet{
		Creator:       "cosmos1creator",
		PubKey:       "testPubKey",
		WalletAddress: "cosmos1wallet3",
		ChainId:      "test-chain",
	}

	_, err := s.msgServer.CreateWallet(sdk.WrapSDKContext(s.ctx), createMsg)
	s.Require().NoError(err)

	// Test recovery
	recoverMsg := &types.MsgRecoverWallet{
		Creator:       "cosmos1newowner",
		WalletAddress: createMsg.WalletAddress,
		NewPubKey:    "newTestPubKey",
		RecoveryProof: "test_proof",
	}

	resp, err := s.msgServer.RecoverWallet(sdk.WrapSDKContext(s.ctx), recoverMsg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	// Verify ownership changed
	wallet, err := s.keeper.GetWallet(s.ctx, createMsg.WalletAddress)
	s.Require().NoError(err)
	s.Require().NotNil(wallet)
	s.Require().Equal(recoverMsg.Creator, wallet.Creator)
	s.Require().Equal(recoverMsg.NewPubKey, wallet.PublicKey)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
