package storage

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Key prefixes for the store
	partyDataPrefix = "party_data"
)

// Storage handles secure storage of TSS key shares
type Storage struct {
	mu    sync.RWMutex
	store prefix.Store
}

// NewStorage creates a new storage instance
func NewStorage(store prefix.Store) *Storage {
	return &Storage{
		store: store,
	}
}

// SavePartyData stores TSS party data for a wallet
func (s *Storage) SavePartyData(ctx sdk.Context, walletAddress string, party1Data, party2Data *keygen.LocalPartySaveData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Marshal party data to JSON
	party1Bytes, err := json.Marshal(party1Data)
	if err != nil {
		return fmt.Errorf("failed to marshal party1 data: %w", err)
	}

	party2Bytes, err := json.Marshal(party2Data)
	if err != nil {
		return fmt.Errorf("failed to marshal party2 data: %w", err)
	}

	// Store party data with wallet address as key
	key1 := []byte(fmt.Sprintf("%s/%s/1", partyDataPrefix, walletAddress))
	key2 := []byte(fmt.Sprintf("%s/%s/2", partyDataPrefix, walletAddress))

	s.store.Set(key1, party1Bytes)
	s.store.Set(key2, party2Bytes)

	return nil
}

// GetPartyData retrieves TSS party data for a wallet
func (s *Storage) GetPartyData(ctx sdk.Context, walletAddress string) (*keygen.LocalPartySaveData, *keygen.LocalPartySaveData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Get party data from store
	key1 := []byte(fmt.Sprintf("%s/%s/1", partyDataPrefix, walletAddress))
	key2 := []byte(fmt.Sprintf("%s/%s/2", partyDataPrefix, walletAddress))

	party1Bytes := s.store.Get(key1)
	if party1Bytes == nil {
		return nil, nil, fmt.Errorf("party1 data not found for wallet: %s", walletAddress)
	}

	party2Bytes := s.store.Get(key2)
	if party2Bytes == nil {
		return nil, nil, fmt.Errorf("party2 data not found for wallet: %s", walletAddress)
	}

	// Unmarshal party data
	var party1Data, party2Data keygen.LocalPartySaveData
	if err := json.Unmarshal(party1Bytes, &party1Data); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal party1 data: %w", err)
	}

	if err := json.Unmarshal(party2Bytes, &party2Data); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal party2 data: %w", err)
	}

	return &party1Data, &party2Data, nil
}

// DeletePartyData removes TSS party data for a wallet
func (s *Storage) DeletePartyData(ctx sdk.Context, walletAddress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Delete party data from store
	key1 := []byte(fmt.Sprintf("%s/%s/1", partyDataPrefix, walletAddress))
	key2 := []byte(fmt.Sprintf("%s/%s/2", partyDataPrefix, walletAddress))

	s.store.Delete(key1)
	s.store.Delete(key2)

	return nil
}
