package testutil

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cometbft/cometbft/libs/log"
	db "github.com/cometbft/cometbft-db"

	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
	identitytypes "selfchain/x/identity/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx           sdk.Context
	keeper        *keeper.Keeper
	msgServer     types.MsgServer
	storeKey      *storetypes.KVStoreKey
	memKey        *storetypes.MemoryStoreKey
	paramsKey     *storetypes.KVStoreKey
	tParamsKey    *storetypes.TransientStoreKey
	networkParams types.NetworkParams
	mockIdentity  *keeper.MockIdentityKeeper
}

func (s *IntegrationTestSuite) SetupTest() {
	// Initialize store keys
	s.storeKey = storetypes.NewKVStoreKey(types.StoreKey)
	s.memKey = storetypes.NewMemoryStoreKey(types.MemStoreKey)
	s.paramsKey = storetypes.NewKVStoreKey("params")
	s.tParamsKey = storetypes.NewTransientStoreKey("transient_params")

	// Create db and multistore
	db := db.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(s.storeKey, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(s.memKey, storetypes.StoreTypeMemory, nil)
	ms.MountStoreWithDB(s.paramsKey, storetypes.StoreTypeIAVL, db)
	ms.MountStoreWithDB(s.tParamsKey, storetypes.StoreTypeTransient, nil)
	s.Require().NoError(ms.LoadLatestVersion())

	// Create context
	s.ctx = sdk.NewContext(ms, tmproto.Header{Height: 1, Time: time.Now().UTC()}, false, log.NewNopLogger())

	// Initialize encoding config
	encConfig := MakeTestEncodingConfig()

	// Initialize params keeper
	paramsKeeper := paramtypes.NewSubspace(
		encConfig.Codec,
		encConfig.Amino,
		s.paramsKey,
		s.tParamsKey,
		"KeylessParams",
	)
	paramsKeeper = paramsKeeper.WithKeyTable(types.ParamKeyTable())

	// Initialize mock identity keeper
	s.mockIdentity = keeper.NewMockIdentityKeeper()

	// Create keeper
	s.keeper = keeper.NewKeeper(
		encConfig.Codec,
		s.storeKey,
		s.memKey,
		paramsKeeper,
		s.mockIdentity,
	)

	// Create message server
	s.msgServer = keeper.NewMsgServerImpl(*s.keeper)

	// Set up test network params
	s.networkParams = types.NetworkParams{
		NetworkType:      "testnet",
		ChainId:         "test-chain",
		SigningAlgorithm: "ECDSA",
		CurveType:       "secp256k1",
		AddressPrefix:   "self",
	}

	// Set up initial params
	s.keeper.SetParams(s.ctx, types.DefaultParams())
}

// Helper function to create and authorize a wallet
func (s *IntegrationTestSuite) createAndAuthorizeWallet(creator, pubKey, walletAddress string) error {
	// Create wallet
	createMsg := &types.MsgCreateWallet{
		Creator:       creator,
		PublicKey:     pubKey,
		WalletAddress: walletAddress,
		ChainId:       s.networkParams.ChainId,
	}

	_, err := s.msgServer.CreateWallet(sdk.WrapSDKContext(s.ctx), createMsg)
	if err != nil {
		return err
	}

	// Set up DID document
	did := "did:self:" + creator
	now := time.Now()
	didDoc := &identitytypes.DIDDocument{
		Id:      did,
		Created: &now,
		Updated: &now,
		Status:  identitytypes.Status_STATUS_ACTIVE,
		VerificationMethod: []*identitytypes.VerificationMethod{
			{
				Id:              did + "#key1",
				Type:            "Ed25519VerificationKey2018",
				Controller:      did,
				PublicKeyBase58: pubKey,
			},
		},
	}
	s.mockIdentity.SetDIDDocument(s.ctx, did, didDoc)

	// Authorize wallet for the creator
	err = s.keeper.ValidateWalletAccess(s.ctx, creator, walletAddress)
	if err != nil {
		return err
	}

	return nil
}

func (s *IntegrationTestSuite) TestCreateWalletFlow() {
	// Set up test data
	creator := "cosmos1creator"
	pubKey := "testPubKey"
	walletAddress := "cosmos1wallet"

	err := s.createAndAuthorizeWallet(creator, pubKey, walletAddress)
	s.Require().NoError(err)

	// Verify wallet was created
	wallet, err := s.keeper.GetWallet(s.ctx, walletAddress)
	s.Require().NoError(err)
	s.Require().NotNil(wallet)
	s.Require().Equal(pubKey, wallet.PublicKey)
}

func (s *IntegrationTestSuite) TestSignTransactionFlow() {
	// First create and authorize wallet
	creator := "cosmos1creator"
	pubKey := "testPubKey"
	walletAddress := "cosmos1wallet"

	err := s.createAndAuthorizeWallet(creator, pubKey, walletAddress)
	s.Require().NoError(err)

	// Now test signing
	unsignedTx := "testUnsignedTx"
	signMsg := &types.MsgSignTransaction{
		Creator:       creator,
		WalletAddress: walletAddress,
		UnsignedTx:    unsignedTx,
	}

	resp, err := s.msgServer.SignTransaction(sdk.WrapSDKContext(s.ctx), signMsg)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp.SignedTx)
}

func (s *IntegrationTestSuite) TestWalletRecoveryFlow() {
	// First create and authorize wallet
	creator := "cosmos1creator"
	pubKey := "testPubKey"
	walletAddress := "cosmos1wallet"

	err := s.createAndAuthorizeWallet(creator, pubKey, walletAddress)
	s.Require().NoError(err)

	// Test recovery
	newPubKey := "newTestPubKey"
	recoveryProof := "recovery_proof"
	recoveryMsg := &types.MsgRecoverWallet{
		Creator:       creator,
		WalletAddress: walletAddress,
		NewPubKey:     newPubKey,
		RecoveryProof: recoveryProof,
	}

	resp, err := s.msgServer.RecoverWallet(sdk.WrapSDKContext(s.ctx), recoveryMsg)
	s.Require().NoError(err)
	s.Require().Equal(walletAddress, resp.WalletAddress)

	// Verify wallet was recovered
	wallet, err := s.keeper.GetWallet(s.ctx, walletAddress)
	s.Require().NoError(err)
	s.Require().NotNil(wallet)
	s.Require().Equal(newPubKey, wallet.PublicKey)
}

func (s *IntegrationTestSuite) TestBatchSignFlow() {
	// First create and authorize wallet
	creator := "cosmos1creator"
	pubKey := "testPubKey"
	walletAddress := "cosmos1wallet"

	err := s.createAndAuthorizeWallet(creator, pubKey, walletAddress)
	s.Require().NoError(err)

	// Now test batch signing
	messages := [][]byte{
		[]byte("tx1"),
		[]byte("tx2"),
	}
	batchSignMsg := &types.MsgBatchSignRequest{
		Creator:       creator,
		WalletAddress: walletAddress,
		Messages:      messages,
	}

	resp, err := s.msgServer.BatchSign(sdk.WrapSDKContext(s.ctx), batchSignMsg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	// Verify batch sign status
	status, err := s.keeper.GetBatchSignStatus(s.ctx, walletAddress)
	s.Require().NoError(err)
	s.Require().Equal(types.BatchSignStatus_BATCH_SIGN_STATUS_COMPLETED, status.Status)
}

func (s *IntegrationTestSuite) TestKeyRotationFlow() {
	// First create and authorize wallet
	creator := "cosmos1creator"
	pubKey := "testPubKey"
	walletAddress := "cosmos1wallet"

	err := s.createAndAuthorizeWallet(creator, pubKey, walletAddress)
	s.Require().NoError(err)

	// Test key rotation
	newPubKey := "newTestPubKey"
	initiateMsg := &types.MsgInitiateKeyRotation{
		Creator:       creator,
		WalletAddress: walletAddress,
		NewPubKey:     newPubKey,
		Signature:     "test-signature",
	}

	initiateResp, err := s.msgServer.InitiateKeyRotation(sdk.WrapSDKContext(s.ctx), initiateMsg)
	s.Require().NoError(err)
	s.Require().NotNil(initiateResp)
	s.Require().NotZero(initiateResp.NewVersion)

	// Complete key rotation
	completeMsg := &types.MsgCompleteKeyRotation{
		Creator:       creator,
		WalletAddress: walletAddress,
		Version:      strconv.FormatUint(uint64(initiateResp.NewVersion), 10),
		Signature:    "test-signature",
	}

	completeResp, err := s.msgServer.CompleteKeyRotation(sdk.WrapSDKContext(s.ctx), completeMsg)
	s.Require().NoError(err)
	s.Require().NotNil(completeResp)
	s.Require().Equal(walletAddress, completeResp.WalletAddress)

	// Verify key was rotated
	wallet, err := s.keeper.GetWallet(s.ctx, walletAddress)
	s.Require().NoError(err)
	s.Require().NotNil(wallet)
	s.Require().Equal(newPubKey, wallet.PublicKey)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
