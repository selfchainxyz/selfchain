package testutil

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"selfchain/x/keyless/keeper"
	testkeeper "selfchain/x/keyless/testutil/keeper"
	"selfchain/x/keyless/testutil/mocks"
	"selfchain/x/keyless/types"
	identitytypes "selfchain/x/identity/types"
)

type IntegrationTestSuite struct {
	suite.Suite
	keeper         *keeper.Keeper
	ctx            sdk.Context
	msgServer      types.MsgServer
	identityKeeper *mocks.IdentityKeeper
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("Setting up integration test suite")
}

func (s *IntegrationTestSuite) SetupTest() {
	// Create a fresh keeper and context for each test
	k, ctx := testkeeper.NewTestKeeper(s.T())
	s.identityKeeper = mocks.NewIdentityKeeper()

	// Mock identity keeper methods
	s.identityKeeper.On("GenerateRecoveryToken", mock.Anything, mock.Anything).Return("test_recovery_token", nil)
	s.identityKeeper.On("VerifyRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.identityKeeper.On("ValidateRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.identityKeeper.On("VerifyRecoveryProof", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.identityKeeper.On("ValidateIdentityStatus", mock.Anything, mock.Anything).Return(nil)
	s.identityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.identityKeeper.On("GetDIDDocument", mock.Anything, mock.Anything).Return(identitytypes.DIDDocument{}, true)
	s.identityKeeper.On("VerifyDIDOwnership", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.identityKeeper.On("GetKeyShare", mock.Anything, mock.Anything).Return([]byte("test_key_share"), true)
	s.identityKeeper.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil)
	s.identityKeeper.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Initialize keeper with the identity keeper
	k2 := keeper.NewKeeper(
		k.GetCodec(),
		k.GetStoreKey(),
		k.GetMemKey(),
		k.GetParamStore(),
		s.identityKeeper,
		nil, // TSS protocol can be nil for tests
	)
	s.keeper = k2

	s.ctx = ctx
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
}

func (s *IntegrationTestSuite) TestA_CreateWallet() {
	msg := &types.MsgCreateWallet{
		Creator:       "cosmos1vx8knpllrj7n963p9ttd80w47kpacrhuts497x",
		PubKey:       "testPubKey",
		WalletAddress: "cosmos1m9l358xunhhwds0568za49mzhvuxx9uxre5tud",
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
		Creator:       "cosmos1vx8knpllrj7n963p9ttd80w47kpacrhuts497x",
		PubKey:       "testPubKey",
		WalletAddress: "cosmos1m9l358xunhhwds0568za49mzhvuxx9uxre5tud",
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
		Creator:       "cosmos1vx8knpllrj7n963p9ttd80w47kpacrhuts497x",
		PubKey:       "testPubKey",
		WalletAddress: "cosmos1m9l358xunhhwds0568za49mzhvuxx9uxre5tud",
		ChainId:      "test-chain",
	}

	_, err := s.msgServer.CreateWallet(sdk.WrapSDKContext(s.ctx), createMsg)
	s.Require().NoError(err)

	// Set wallet status to active
	err = s.keeper.SetWalletStatus(s.ctx, createMsg.WalletAddress, types.WalletStatus_WALLET_STATUS_ACTIVE)
	s.Require().NoError(err)

	// Grant recovery permission to new owner
	newOwner := "cosmos1qzskhrcjnkdz2ln4yeafzsdwht8ch08j4wed69"
	expiresAt := time.Now().Add(24 * time.Hour)
	permission := &types.Permission{
		WalletAddress: createMsg.WalletAddress,
		Grantee:      newOwner,
		Permissions:  []string{types.WalletPermission_WALLET_PERMISSION_RECOVER.String()},
		ExpiresAt:    &expiresAt,
	}

	err = s.keeper.GrantPermission(s.ctx, permission)
	s.Require().NoError(err)

	// Create recovery session
	err = s.keeper.CreateRecoverySession(s.ctx, newOwner, createMsg.WalletAddress)
	s.Require().NoError(err)

	// Test recovery before timelock - should fail
	recoverMsg := &types.MsgRecoverWallet{
		Creator:       newOwner,
		WalletAddress: createMsg.WalletAddress,
		NewPubKey:    "newTestPubKey",
		RecoveryProof: "test_proof",
	}

	_, err = s.msgServer.RecoverWallet(sdk.WrapSDKContext(s.ctx), recoverMsg)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "recovery not allowed")

	// Advance block time by 24 hours
	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(25 * time.Hour))

	// Test recovery after timelock - should succeed
	resp, err := s.msgServer.RecoverWallet(sdk.WrapSDKContext(s.ctx), recoverMsg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	// Verify wallet was recovered
	wallet, err := s.keeper.GetWallet(s.ctx, createMsg.WalletAddress)
	s.Require().NoError(err)
	s.Require().NotNil(wallet)
	s.Require().Equal(newOwner, wallet.Creator)
	s.Require().Equal(recoverMsg.NewPubKey, wallet.PublicKey)
	s.Require().Equal(createMsg.ChainId, wallet.ChainId)
	s.Require().Equal(types.WalletStatus_WALLET_STATUS_ACTIVE, wallet.Status)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
