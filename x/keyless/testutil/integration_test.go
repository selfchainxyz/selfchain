package testutil

import (
	"context"
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
	"encoding/json"
	"math/big"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
)

type IntegrationTestSuite struct {
	suite.Suite
	keeper         *keeper.Keeper
	ctx            sdk.Context
	msgServer      types.MsgServer
	identityKeeper *mocks.IdentityKeeper
}

type MockTSSProtocol struct {
	mock.Mock
}

func (m *MockTSSProtocol) SignMessage(
	ctx context.Context,
	msg []byte,
	localPartySaveData []byte,
	remotePartySaveData []byte,
) ([]byte, error) {
	args := m.Called(ctx, msg, localPartySaveData, remotePartySaveData)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockTSSProtocol) InitiateSigning(
	ctx context.Context,
	msg []byte,
	walletID string,
) (*types.SigningResponse, error) {
	args := m.Called(ctx, msg, walletID)
	if response, ok := args.Get(0).(*types.SigningResponse); ok {
		return response, args.Error(1)
	}
	return &types.SigningResponse{
		Signature: []byte("mocked_signature"),
	}, nil
}

func (m *MockTSSProtocol) ProcessKeyGenRound(
	ctx context.Context,
	sessionID string,
	partyData *types.PartyData,
) error {
	args := m.Called(ctx, sessionID, partyData)
	return args.Error(0)
}

func (m *MockTSSProtocol) GenerateKeyShares(
	ctx context.Context,
	request *types.KeyGenRequest,
) (*types.KeyGenResponse, error) {
	args := m.Called(ctx, request)
	if response, ok := args.Get(0).(*types.KeyGenResponse); ok {
		return response, args.Error(1)
	}
	return &types.KeyGenResponse{
		WalletAddress: request.WalletAddress,
		PublicKey:     []byte("mocked_public_key"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     time.Now(),
			LastRotated:   time.Now(),
			LastUsed:      time.Now(),
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: request.SecurityLevel,
		},
	}, nil
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
	s.identityKeeper.On("GetKeyShare", mock.Anything, mock.Anything).Return([]byte(`{
		"LocalPreParams": {
			"PaillierSK": {
				"L": "test_l",
				"U": "test_u"
			},
			"NTildei": "test_ntildei",
			"H1i": "test_h1i",
			"H2i": "test_h2i"
		},
		"LocalSecrets": {
			"Xi": "test_xi",
			"ShareID": "test_share_id"
		},
		"LocalData": {
			"Ks": ["test_k1", "test_k2"],
			"NTildej": ["test_ntildej1", "test_ntildej2"],
			"H1j": ["test_h1j1", "test_h1j2"],
			"H2j": ["test_h2j1", "test_h2j2"]
		}
	}`), true)
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

	// Store key shares for both creator and wallet
	preParams1, err := keygen.GeneratePreParams(1 * time.Minute)
	s.Require().NoError(err)
	preParams2, err := keygen.GeneratePreParams(1 * time.Minute)
	s.Require().NoError(err)

	// Create party IDs
	p1ID := tss.NewPartyID("party1", "P1", big.NewInt(1))
	p2ID := tss.NewPartyID("party2", "P2", big.NewInt(2))
	parties := tss.SortPartyIDs([]*tss.PartyID{p1ID, p2ID})

	// Create peer context
	peerCtx := tss.NewPeerContext(parties)
	params1 := tss.NewParameters(tss.S256(), peerCtx, p1ID, len(parties), 1)
	params2 := tss.NewParameters(tss.S256(), peerCtx, p2ID, len(parties), 1)

	// Create channels
	outCh1 := make(chan tss.Message, len(parties))
	outCh2 := make(chan tss.Message, len(parties))
	endCh := make(chan *keygen.LocalPartySaveData, len(parties))
	errCh := make(chan *tss.Error, len(parties))

	// Create keygen parties with unique pre-parameters
	p1 := keygen.NewLocalParty(params1, outCh1, endCh, *preParams1).(*keygen.LocalParty)
	p2 := keygen.NewLocalParty(params2, outCh2, endCh, *preParams2).(*keygen.LocalParty)

	// Start parties
	go func() {
		if err := p1.Start(); err != nil {
			errCh <- err
		}
	}()
	go func() {
		if err := p2.Start(); err != nil {
			errCh <- err
		}
	}()

	// Message routing for party 1
	go func() {
		for msg := range outCh1 {
			wireBytes, _, err := msg.WireBytes()
			if err != nil {
				errCh <- tss.NewError(err, "failed to get wire bytes", 1, msg.GetFrom(), msg.GetFrom())
				continue
			}

			dest := msg.GetTo()
			if dest == nil { // broadcast
				if _, err := p2.UpdateFromBytes(wireBytes, msg.GetFrom(), true); err != nil {
					errCh <- err
					continue
				}
			} else if dest[0].Index == p2.PartyID().Index {
				if _, err := p2.UpdateFromBytes(wireBytes, msg.GetFrom(), false); err != nil {
					errCh <- err
					continue
				}
			}
		}
	}()

	// Message routing for party 2
	go func() {
		for msg := range outCh2 {
			wireBytes, _, err := msg.WireBytes()
			if err != nil {
				errCh <- tss.NewError(err, "failed to get wire bytes", 2, msg.GetFrom(), msg.GetFrom())
				continue
			}

			dest := msg.GetTo()
			if dest == nil { // broadcast
				if _, err := p1.UpdateFromBytes(wireBytes, msg.GetFrom(), true); err != nil {
					errCh <- err
					continue
				}
			} else if dest[0].Index == p1.PartyID().Index {
				if _, err := p1.UpdateFromBytes(wireBytes, msg.GetFrom(), false); err != nil {
					errCh <- err
					continue
				}
			}
		}
	}()

	// Wait for both parties to finish
	var party1Data, party2Data *keygen.LocalPartySaveData
	for i := 0; i < 2; i++ {
		select {
		case err := <-errCh:
			s.T().Fatalf("Error in key generation: %v", err)
			return
		case data := <-endCh:
			if data.ShareID.Cmp(big.NewInt(1)) == 0 {
				party1Data = data
			} else {
				party2Data = data
			}
		}
	}

	// Store key shares
	saveDataBytes1, err := json.Marshal(party1Data)
	s.Require().NoError(err)
	saveDataBytes2, err := json.Marshal(party2Data)
	s.Require().NoError(err)

	err = s.keeper.StoreKeyShare(s.ctx, createMsg.Creator, saveDataBytes1)
	s.Require().NoError(err)
	err = s.keeper.StoreKeyShare(s.ctx, createMsg.WalletAddress, saveDataBytes2)
	s.Require().NoError(err)

	// Create and configure mock TSS protocol for signing
	mockTSS := new(MockTSSProtocol)
	mockTSS.On("SignMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]byte("mocked_signature"), nil)
	mockTSS.On("GenerateKeyShares", mock.Anything, mock.Anything).
		Return(&types.KeyGenResponse{
			WalletAddress: createMsg.WalletAddress,
			PublicKey:     []byte("mocked_public_key"),
			Metadata: &types.KeyMetadata{
				CreatedAt:     time.Now(),
				LastRotated:   time.Now(),
				LastUsed:      time.Now(),
				UsageCount:    0,
				BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
				SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
			},
		}, nil)
	mockTSS.On("InitiateSigning", mock.Anything, mock.Anything, mock.Anything).
		Return(&types.SigningResponse{
			Signature: []byte("mocked_signature"),
		}, nil)
	mockTSS.On("ProcessKeyGenRound", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	s.keeper.SetTSSProtocol(mockTSS)

	// Test signing
	signMsg := &types.MsgSignTransaction{
		Creator:       createMsg.Creator,
		WalletAddress: createMsg.WalletAddress,
		UnsignedTx:   "test_tx",
	}

	resp, err := s.msgServer.SignTransaction(sdk.WrapSDKContext(s.ctx), signMsg)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal("mocked_signature", string(resp.SignedTx))
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
