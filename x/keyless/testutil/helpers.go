package testutil

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"selfchain/x/keyless/keeper"
	"selfchain/x/keyless/testutil/mocks"
	"selfchain/x/keyless/types"
)

// CreateTestWallet creates a new wallet for testing
func CreateTestWallet(t *testing.T, k *keeper.Keeper, ctx sdk.Context) *types.Wallet {
	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey().Bytes()
	now := time.Now().UTC()
	addr := sdk.AccAddress(privKey.PubKey().Address()).String()

	wallet := &types.Wallet{
		Id:            addr,
		Creator:       addr,
		PublicKey:     string(pubKey),
		WalletAddress: addr,
		ChainId:       "test-chain",
		Status:        types.WalletStatus_WALLET_STATUS_ACTIVE,
		KeyVersion:    1,
		CreatedAt:     &now,
		UpdatedAt:     &now,
		LastUsed:      &now,
		UsageCount:    0,
	}

	err := k.SaveWallet(ctx, wallet)
	require.NoError(t, err)

	return wallet
}

// GenerateValidRecoveryProof creates a valid recovery proof for testing
func GenerateValidRecoveryProof(t *testing.T, wallet *types.Wallet) []byte {
	// Generate a mock recovery proof that will pass verification
	proof := []byte("valid_recovery_proof")
	return proof
}

// SetupMockIdentityKeeper configures the mock identity keeper for testing
func SetupMockIdentityKeeper(t *testing.T) *mocks.IdentityKeeper {
	mockIdentityKeeper := &mocks.IdentityKeeper{}

	// Setup default behavior for all required methods
	mockIdentityKeeper.On("VerifyRecoveryProof", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("ValidateIdentityStatus", mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("CheckRateLimit", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("LogAuditEvent", mock.Anything, mock.Anything).Return(nil)
	mockIdentityKeeper.On("GenerateRecoveryToken", mock.Anything, mock.Anything).Return("test_token", nil)
	mockIdentityKeeper.On("ValidateRecoveryToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	return mockIdentityKeeper
}

// SetupMockTSSProtocol configures the mock TSS protocol for testing
func SetupMockTSSProtocol(t *testing.T) *mocks.TSSProtocol {
	mockTSS := &mocks.TSSProtocol{}
	now := time.Now().UTC()
	// Setup default behavior for all required methods
	mockTSS.On("GenerateKey", mock.Anything, mock.Anything, mock.Anything).Return([]byte("test_pubkey"), nil)
	mockTSS.On("Sign", mock.Anything, mock.Anything, mock.Anything).Return([]byte("test_signature"), nil)
	mockTSS.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockTSS.On("GenerateKeyShares", mock.Anything, mock.Anything).Return(&types.KeyGenResponse{
		WalletAddress: "test_wallet",
		PublicKey:     []byte("test_pubkey"),
		Metadata: &types.KeyMetadata{
			CreatedAt:   now,
			LastRotated: now,
			LastUsed:    now,
			UsageCount:  0,
		},
	}, nil)

	return mockTSS
}
