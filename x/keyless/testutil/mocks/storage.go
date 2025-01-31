package mocks

import (
	"context"
	"fmt"
	"sync"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"

	"selfchain/x/keyless/types"
)

// MockStorage implements the Storage interface for testing
type MockStorage struct {
	shares    map[string]*types.EncryptedShare
	partyData map[string]*keygen.LocalPartySaveData
	metadata  map[string]*types.WalletMetadata
	mu        sync.RWMutex
}

// SetupMockStorage creates a new mock storage instance
func SetupMockStorage() *MockStorage {
	return &MockStorage{
		shares:    make(map[string]*types.EncryptedShare),
		partyData: make(map[string]*keygen.LocalPartySaveData),
		metadata:  make(map[string]*types.WalletMetadata),
	}
}

// SavePartyShare stores an encrypted key share
func (m *MockStorage) SavePartyShare(ctx context.Context, key string, share *types.EncryptedShare) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shares[key] = share
	return nil
}

// GetPartyShare retrieves an encrypted key share
func (m *MockStorage) GetPartyShare(ctx context.Context, key string) (*types.EncryptedShare, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	share, ok := m.shares[key]
	if !ok {
		return nil, fmt.Errorf("share not found: %s", key)
	}
	return share, nil
}

// DeletePartyShare removes a key share
func (m *MockStorage) DeletePartyShare(ctx context.Context, key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.shares, key)
}

// SavePartyData stores TSS party data for a wallet
func (m *MockStorage) SavePartyData(ctx context.Context, walletAddress string, party1Data, party2Data *keygen.LocalPartySaveData) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.partyData[walletAddress+"_party1"] = party1Data
	m.partyData[walletAddress+"_party2"] = party2Data
	return nil
}

// GetPartyData retrieves TSS party data for a wallet
func (m *MockStorage) GetPartyData(ctx context.Context, walletAddress string) (*keygen.LocalPartySaveData, *keygen.LocalPartySaveData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	party1Data, ok := m.partyData[walletAddress+"_party1"]
	if !ok {
		return nil, nil, fmt.Errorf("party1 data not found for wallet: %s", walletAddress)
	}
	party2Data, ok := m.partyData[walletAddress+"_party2"]
	if !ok {
		return nil, nil, fmt.Errorf("party2 data not found for wallet: %s", walletAddress)
	}
	return party1Data, party2Data, nil
}

// DeletePartyData removes TSS party data for a wallet
func (m *MockStorage) DeletePartyData(ctx context.Context, walletAddress string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.partyData, walletAddress+"_party1")
	delete(m.partyData, walletAddress+"_party2")
	return nil
}

// GetWalletMetadata retrieves metadata for a wallet
func (m *MockStorage) GetWalletMetadata(ctx context.Context, walletAddress string) (*types.WalletMetadata, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	metadata, ok := m.metadata[walletAddress]
	if !ok {
		return nil, fmt.Errorf("metadata not found for wallet: %s", walletAddress)
	}
	return metadata, nil
}

// SaveWalletMetadata stores metadata for a wallet
func (m *MockStorage) SaveWalletMetadata(ctx context.Context, walletAddress string, metadata *types.WalletMetadata) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metadata[walletAddress] = metadata
	return nil
}

// Helper method for tests to pre-populate shares
func (m *MockStorage) AddTestShare(key string, share *types.EncryptedShare) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shares[key] = share
}

// Helper method for tests to pre-populate party data
func (m *MockStorage) AddTestPartyData(walletAddress string, party1Data, party2Data *keygen.LocalPartySaveData) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.partyData[walletAddress+"_party1"] = party1Data
	m.partyData[walletAddress+"_party2"] = party2Data
}

// Helper method for tests to pre-populate metadata
func (m *MockStorage) AddTestMetadata(walletAddress string, metadata *types.WalletMetadata) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metadata[walletAddress] = metadata
}

// MockStorageWrapper wraps MockStorage to add additional functionality for testing
type MockStorageWrapper struct {
	*MockStorage
	shares     map[string]*types.EncryptedShare
	partyData  map[string]*keygen.LocalPartySaveData
	metadata   map[string]*types.WalletMetadata
}

// GetWalletMetadata retrieves metadata for a wallet
func (m *MockStorageWrapper) GetWalletMetadata(ctx context.Context, walletAddress string) (*types.WalletMetadata, error) {
	metadata, ok := m.metadata[walletAddress]
	if !ok {
		return nil, types.ErrWalletMetadataNotFound
	}
	return metadata, nil
}

// SaveWalletMetadata stores metadata for a wallet
func (m *MockStorageWrapper) SaveWalletMetadata(ctx context.Context, walletAddress string, metadata *types.WalletMetadata) error {
	if m.metadata == nil {
		m.metadata = make(map[string]*types.WalletMetadata)
	}
	m.metadata[walletAddress] = metadata
	return nil
}
