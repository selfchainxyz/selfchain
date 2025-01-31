package mocks

import (
	"context"
	"crypto/elliptic"
	"fmt"
	"time"
	"math/big"

	"github.com/bnb-chain/tss-lib/v2/crypto"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"selfchain/x/keyless/types"
)

// PrePopulateShares adds test shares to the mock storage
func PrePopulateShares(m *MockStorage, walletAddress string, chainId string, securityLevel types.SecurityLevel) error {
	// Create test shares
	share1 := &types.EncryptedShare{
		EncryptedData: "test_share_1",
		KeyId:        walletAddress + "_1",
		Version:      1,
		CreatedAt:    time.Now(),
	}
	share2 := &types.EncryptedShare{
		EncryptedData: "test_share_2",
		KeyId:        walletAddress + "_2",
		Version:      1,
		CreatedAt:    time.Now(),
	}

	// Save shares
	key1 := fmt.Sprintf("%s_share_1", walletAddress)
	key2 := fmt.Sprintf("%s_share_2", walletAddress)
	if err := m.SavePartyShare(context.Background(), key1, share1); err != nil {
		return fmt.Errorf("failed to save share 1: %w", err)
	}
	if err := m.SavePartyShare(context.Background(), key2, share2); err != nil {
		return fmt.Errorf("failed to save share 2: %w", err)
	}

	return nil
}

// PrePopulatePartyData adds test party data to the mock storage
func PrePopulatePartyData(m *MockStorage, walletAddress string) error {
	// Create test secrets
	secret1 := big.NewInt(123)
	secret2 := big.NewInt(456)

	// Create test ECPoints
	curve, err := crypto.NewECPoint(elliptic.P256(), big.NewInt(1), big.NewInt(2))
	if err != nil {
		return fmt.Errorf("failed to create EC point: %w", err)
	}

	// Create party data
	party1Data := keygen.NewLocalPartySaveData(2)
	party1Data.Ks = []*big.Int{secret1}
	party1Data.BigXj = []*crypto.ECPoint{curve}
	party1Data.ECDSAPub = curve

	party2Data := keygen.NewLocalPartySaveData(2)
	party2Data.Ks = []*big.Int{secret2}
	party2Data.BigXj = []*crypto.ECPoint{curve}
	party2Data.ECDSAPub = curve

	// Save party data
	if err := m.SavePartyData(context.Background(), walletAddress, &party1Data, &party2Data); err != nil {
		return fmt.Errorf("failed to save party data: %w", err)
	}

	return nil
}

// PrePopulateMetadata adds test metadata to the mock storage
func PrePopulateMetadata(m *MockStorage, walletAddress string, chainId string, securityLevel types.SecurityLevel) error {
	metadata := &types.WalletMetadata{
		ChainId:       chainId,
		SecurityLevel: securityLevel,
	}

	if err := m.SaveWalletMetadata(context.Background(), walletAddress, metadata); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}
