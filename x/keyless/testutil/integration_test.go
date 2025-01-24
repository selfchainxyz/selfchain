package testutil

import (
	"context"
	"testing"
	"time"
	"strconv"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
	identitytypes "selfchain/x/identity/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx           sdk.Context
	keeper        *keeper.Keeper
	msgServer     types.MsgServer
	networkParams *types.NetworkParams
	mockIdentity  *keeper.MockIdentityKeeper
}

func (s *IntegrationTestSuite) SetupTest() {
	// Create test context
	s.ctx = sdk.NewContext(
		NewTestMultiStore(s.T(), sdk.NewKVStoreKey("test")),
		tmproto.Header{Height: 1, Time: time.Now().UTC()},
		false,
		nil,
	)

	// Create mock identity keeper
	s.mockIdentity = keeper.NewMockIdentityKeeper()

	// Initialize keeper
	s.keeper = keeper.NewKeeper(
		MakeTestEncodingConfig(),
		sdk.NewKVStoreKey(types.StoreKey),
		sdk.NewKVStoreKey("memKey"),
		paramtypes.Subspace{},
		s.mockIdentity,
	)

	// Initialize message server
	s.msgServer = keeper.NewMsgServerImpl(*s.keeper)

	// Set up test network params
	s.networkParams = &types.NetworkParams{
		NetworkType:      "testnet",
		ChainId:         "test-chain",
		SigningAlgorithm: "ECDSA",
		CurveType:       "secp256k1",
		AddressPrefix:   "self",
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// Test wallet creation flow
func (s *IntegrationTestSuite) TestCreateWalletFlow() {
	// Set up test data
	creator := "cosmos1creator"
	pubKey := "testPubKey"
	walletAddress := "cosmos1wallet"

	// Create wallet
	createMsg := &types.MsgCreateWallet{
		Creator:       creator,
		PubKey:        pubKey,
		WalletAddress: walletAddress,
		ChainId:      s.networkParams.ChainId,
	}

	// Execute create wallet
	resp, err := s.msgServer.CreateWallet(context.Background(), createMsg)
	s.Require().NoError(err)
	s.Require().Equal(walletAddress, resp.WalletAddress)

	// Verify wallet was created
	wallet, found := s.keeper.GetWalletFromStore(s.ctx, walletAddress)
	s.Require().True(found)
	s.Require().Equal(pubKey, wallet.PublicKey)
}

// Test transaction signing flow
func (s *IntegrationTestSuite) TestSignTransactionFlow() {
	// Set up test wallet
	creator := "cosmos1creator"
	pubKey := "testPubKey"
	walletAddress := "cosmos1wallet"
	unsignedTx := "testUnsignedTx"

	// Create wallet first
	createMsg := &types.MsgCreateWallet{
		Creator:       creator,
		PubKey:        pubKey,
		WalletAddress: walletAddress,
		ChainId:      s.networkParams.ChainId,
	}
	_, err := s.msgServer.CreateWallet(context.Background(), createMsg)
	s.Require().NoError(err)

	// Sign transaction
	signMsg := &types.MsgSignTransaction{
		Creator:       creator,
		WalletAddress: walletAddress,
		UnsignedTx:    unsignedTx,
		ChainId:      s.networkParams.ChainId,
	}

	resp, err := s.msgServer.SignTransaction(context.Background(), signMsg)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp.SignedTx)
}

// Test wallet recovery flow
func (s *IntegrationTestSuite) TestWalletRecoveryFlow() {
	// Set up test data
	creator := "cosmos1creator"
	pubKey := "testPubKey"
	walletAddress := "cosmos1wallet"
	recoveryProof := "testRecoveryProof"
	newPubKey := "newTestPubKey"

	// Create wallet first
	createMsg := &types.MsgCreateWallet{
		Creator:       creator,
		PubKey:        pubKey,
		WalletAddress: walletAddress,
		ChainId:      s.networkParams.ChainId,
	}
	_, err := s.msgServer.CreateWallet(context.Background(), createMsg)
	s.Require().NoError(err)

	// Create test DID document
	did := "did:self:123"
	now := time.Now()
	didDoc := identitytypes.DIDDocument{
		Id:        did,
		Created:   &now,
		Updated:   &now,
		Status:    identitytypes.Status_STATUS_ACTIVE,
		VerificationMethod: []*identitytypes.VerificationMethod{
			{
				Id:              did + "#key1",
				Type:            "Ed25519VerificationKey2018",
				Controller:      did,
				PublicKeyBase58: "test-pubkey",
			},
		},
	}

	// Set DID document in mock keeper
	s.mockIdentity.SetDIDDocument(s.ctx, did, &didDoc)

	// Recover wallet
	recoverMsg := &types.MsgRecoverWallet{
		Creator:       creator,
		WalletAddress: walletAddress,
		RecoveryProof: recoveryProof,
		NewPubKey:     newPubKey,
	}

	resp, err := s.msgServer.RecoverWallet(context.Background(), recoverMsg)
	s.Require().NoError(err)
	s.Require().Equal(walletAddress, resp.WalletAddress)

	// Verify wallet was recovered
	wallet, found := s.keeper.GetWalletFromStore(s.ctx, walletAddress)
	s.Require().True(found)
	s.Require().Equal(newPubKey, wallet.PublicKey)
}

// Test batch signing flow
func (s *IntegrationTestSuite) TestBatchSignFlow() {
	// Set up test wallet
	creator := "cosmos1creator"
	pubKey := "testPubKey"
	walletAddress := "cosmos1wallet"
	unsignedTxs := []string{"tx1", "tx2", "tx3"}

	// Create wallet first
	createMsg := &types.MsgCreateWallet{
		Creator:       creator,
		PubKey:        pubKey,
		WalletAddress: walletAddress,
		ChainId:      s.networkParams.ChainId,
	}
	_, err := s.msgServer.CreateWallet(context.Background(), createMsg)
	s.Require().NoError(err)

	// Batch sign transactions
	batchSignMsg := &types.MsgBatchSign{
		Creator:       creator,
		WalletAddress: walletAddress,
		UnsignedTxs:   unsignedTxs,
		ChainId:      s.networkParams.ChainId,
	}

	resp, err := s.msgServer.BatchSign(context.Background(), batchSignMsg)
	s.Require().NoError(err)
	s.Require().Equal(len(unsignedTxs), len(resp.SignedTxs))
}

// Test key rotation flow
func (s *IntegrationTestSuite) TestKeyRotationFlow() {
	// Set up test wallet
	creator := "cosmos1creator"
	pubKey := "testPubKey"
	walletAddress := "cosmos1wallet"
	newPubKey := "newTestPubKey"

	// Create wallet first
	createMsg := &types.MsgCreateWallet{
		Creator:       creator,
		PubKey:        pubKey,
		WalletAddress: walletAddress,
		ChainId:      s.networkParams.ChainId,
	}
	_, err := s.msgServer.CreateWallet(context.Background(), createMsg)
	s.Require().NoError(err)

	// Set up DID document in mock keeper
	did := "did:self:123"
	s.mockIdentity.SetDIDDocument(s.ctx, did, &identitytypes.DIDDocument{
		Id: did,
		// Add other required DID document fields
	})

	// Initiate key rotation
	rotateMsg := &types.MsgInitiateKeyRotation{
		Creator:       creator,
		WalletAddress: walletAddress,
		NewPubKey:     newPubKey,
		Signature:     "test-signature",
	}

	resp, err := s.msgServer.InitiateKeyRotation(context.Background(), rotateMsg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().NotZero(resp.NewVersion)

	// Complete key rotation
	completeMsg := &types.MsgCompleteKeyRotation{
		Creator:       creator,
		WalletAddress: walletAddress,
		Version:      strconv.FormatUint(uint64(resp.NewVersion), 10),
		Signature:    "test-signature",
	}

	completeResp, err := s.msgServer.CompleteKeyRotation(context.Background(), completeMsg)
	s.Require().NoError(err)
	s.Require().NotNil(completeResp)
	s.Require().Equal(walletAddress, completeResp.WalletAddress)

	// Verify key was rotated
	wallet, found := s.keeper.GetWalletFromStore(s.ctx, walletAddress)
	s.Require().True(found)
	s.Require().Equal(newPubKey, wallet.PublicKey)
}
