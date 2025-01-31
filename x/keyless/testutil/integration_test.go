package testutil

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	db "github.com/cometbft/cometbft-db"

	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/types"
	identitytypes "selfchain/x/identity/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	keeper         *keeper.Keeper
	identityKeeper *MockIdentityKeeper
	tssProtocol    *MockTSSProtocol
}

// MockIdentityKeeper is a mock type for the IdentityKeeper interface
type MockIdentityKeeper struct {
	mock.Mock
}

func (m *MockIdentityKeeper) GetDIDDocument(ctx sdk.Context, did string) (identitytypes.DIDDocument, bool) {
	args := m.Called(ctx, did)
	return args.Get(0).(identitytypes.DIDDocument), args.Bool(1)
}

func (m *MockIdentityKeeper) SaveDIDDocument(ctx sdk.Context, doc identitytypes.DIDDocument) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockIdentityKeeper) ReconstructWalletFromDID(ctx sdk.Context, doc identitytypes.DIDDocument) ([]byte, error) {
	args := m.Called(ctx, doc)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockIdentityKeeper) VerifyOAuth2Token(ctx sdk.Context, token string, scope string) error {
	args := m.Called(ctx, token, scope)
	return args.Error(0)
}

func (m *MockIdentityKeeper) VerifyMFA(ctx sdk.Context, code string) error {
	args := m.Called(ctx, code)
	return args.Error(0)
}

func (m *MockIdentityKeeper) CheckRateLimit(ctx sdk.Context, did string, action string) error {
	args := m.Called(ctx, did, action)
	return args.Error(0)
}

func (m *MockIdentityKeeper) LogAuditEvent(ctx sdk.Context, event *identitytypes.AuditEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// MockTSSProtocol implements the TSSProtocol interface for testing
type MockTSSProtocol struct {
	mock.Mock
}

func (m *MockTSSProtocol) GenerateKeyShares(ctx sdk.Context, did string, threshold uint32, securityLevel types.SecurityLevel) (*types.KeyGenResponse, error) {
	args := m.Called(ctx, did, threshold, securityLevel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.KeyGenResponse), args.Error(1)
}

func (m *MockTSSProtocol) SignMessage(ctx sdk.Context, msg []byte, shares [][]byte) ([]byte, error) {
	args := m.Called(ctx, msg, shares)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockTSSProtocol) VerifyShare(ctx sdk.Context, share []byte, publicKey []byte) error {
	args := m.Called(ctx, share, publicKey)
	return args.Error(0)
}

func (m *MockTSSProtocol) ReconstructKey(ctx sdk.Context, shares [][]byte) ([]byte, error) {
	args := m.Called(ctx, shares)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockTSSProtocol) GetPartyData(ctx sdk.Context, partyID string) (*types.PartyData, error) {
	args := m.Called(ctx, partyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.PartyData), args.Error(1)
}

func (m *MockTSSProtocol) SetPartyData(ctx sdk.Context, data *types.PartyData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockTSSProtocol) VerifySignature(ctx sdk.Context, message []byte, signature []byte, publicKey []byte) error {
	args := m.Called(ctx, message, signature, publicKey)
	return args.Error(0)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupTest() {
	db := db.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(s.T(), stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := paramtypes.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"KeylessParams",
	)

	s.identityKeeper = new(MockIdentityKeeper)
	s.tssProtocol = new(MockTSSProtocol)

	// Set up mock expectations
	now := time.Now().UTC()
	s.identityKeeper.On("GetDIDDocument", mock.Anything, mock.Anything).Return(identitytypes.DIDDocument{
		Id:         "test-did",
		Controller: []string{"test-controller"},
		Created:    &now,
		Updated:    &now,
		Status:     identitytypes.Status_STATUS_ACTIVE,
	}, true)
	s.identityKeeper.On("SaveDIDDocument", mock.Anything, mock.Anything).Return(nil)
	s.identityKeeper.On("ReconstructWalletFromDID", mock.Anything, mock.Anything).Return([]byte("test-wallet"), nil)
	s.identityKeeper.On("VerifyOAuth2Token", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.identityKeeper.On("VerifyMFA", mock.Anything, mock.Anything).Return(nil)
	s.identityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.identityKeeper.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil)

	s.tssProtocol.On("GenerateKeyShares", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&types.KeyGenResponse{
		WalletAddress: "test-wallet",
		PublicKey:     []byte("test-pubkey"),
		Metadata: &types.KeyMetadata{
			CreatedAt:     now,
			LastRotated:   now,
			LastUsed:      now,
			UsageCount:    0,
			BackupStatus:  types.BackupStatus_BACKUP_STATUS_COMPLETED,
			SecurityLevel: types.SecurityLevel_SECURITY_LEVEL_HIGH,
		},
	}, nil)
	s.tssProtocol.On("SignMessage", mock.Anything, mock.Anything, mock.Anything).Return([]byte("test-signature"), nil)
	s.tssProtocol.On("VerifyShare", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	s.tssProtocol.On("ReconstructKey", mock.Anything, mock.Anything).Return([]byte("test-key"), nil)
	s.tssProtocol.On("GetPartyData", mock.Anything, mock.Anything).Return(&types.PartyData{
		PartyId:          "test-party",
		PublicKey:        []byte("test-pubkey"),
		PartyShare:       []byte("test-share"),
		VerificationData: []byte("test-verification"),
		ChainId:          "test-chain",
		Status:           "active",
	}, nil)
	s.tssProtocol.On("SetPartyData", mock.Anything, mock.Anything).Return(nil)
	s.tssProtocol.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	s.keeper = keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		s.identityKeeper,
		s.tssProtocol,
	)

	s.ctx = sdk.NewContext(stateStore, tmproto.Header{}, false, nil)
}

func (s *IntegrationTestSuite) TestWalletCreation() {
	msg := &types.MsgCreateWallet{
		Creator:       "test-creator",
		WalletAddress: "test-wallet",
		PubKey:       "test-pubkey",
		ChainId:      "test-chain",
	}

	// Test wallet creation
	msgServer := keeper.NewMsgServerImpl(s.keeper)
	response, err := msgServer.CreateWallet(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().NotNil(response)

	// Verify wallet was created
	wallet, found := s.keeper.GetWallet(s.ctx, msg.WalletAddress)
	s.Require().True(found)
	s.Require().NotNil(wallet)
}

func (s *IntegrationTestSuite) TestRecoverySession() {
	createMsg := &types.MsgCreateWallet{
		Creator:       "test-creator",
		WalletAddress: "test-wallet",
		PubKey:       "test-pubkey",
		ChainId:      "test-chain",
	}

	// Create wallet first
	msgServer := keeper.NewMsgServerImpl(s.keeper)
	_, err := msgServer.CreateWallet(s.ctx, createMsg)
	s.Require().NoError(err)

	// Test recovery session creation
	session, err := s.keeper.CreateRecoverySession(s.ctx, createMsg.WalletAddress)
	s.Require().NoError(err)
	s.Require().NotNil(session)

	// Verify wallet exists
	wallet, found := s.keeper.GetWallet(s.ctx, createMsg.WalletAddress)
	s.Require().True(found)
	s.Require().NotNil(wallet)
}
