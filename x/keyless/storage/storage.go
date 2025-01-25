package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"

	"selfchain/x/keyless/types"
)

// Storage manages data persistence
type Storage struct {
	store sdk.KVStore
}

// NewStorage creates a new storage instance
func NewStorage(store sdk.KVStore) *Storage {
	return &Storage{
		store: store,
	}
}

// SavePartyShare stores an encrypted key share
func (s *Storage) SavePartyShare(ctx context.Context, key string, share *types.EncryptedShare) error {
	store := prefix.NewStore(s.store, []byte(types.KeyShareKey))
	bz, err := share.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal share: %w", err)
	}
	store.Set([]byte(key), bz)
	return nil
}

// GetPartyShare retrieves an encrypted key share
func (s *Storage) GetPartyShare(ctx context.Context, key string) (*types.EncryptedShare, error) {
	store := prefix.NewStore(s.store, []byte(types.KeyShareKey))
	bz := store.Get([]byte(key))
	if bz == nil {
		return nil, fmt.Errorf("share not found: %s", key)
	}

	var share types.EncryptedShare
	if err := share.Unmarshal(bz); err != nil {
		return nil, fmt.Errorf("failed to unmarshal share: %w", err)
	}
	return &share, nil
}

// DeletePartyShare removes a key share
func (s *Storage) DeletePartyShare(ctx context.Context, key string) {
	store := prefix.NewStore(s.store, []byte(types.KeyShareKey))
	store.Delete([]byte(key))
}

// SavePartyData stores TSS party data for a wallet
func (s *Storage) SavePartyData(ctx context.Context, walletAddress string, party1Data, party2Data *keygen.LocalPartySaveData) error {
	store := prefix.NewStore(s.store, []byte(types.PartyDataKey))
	party1Key := fmt.Sprintf("%s_party1", walletAddress)
	party2Key := fmt.Sprintf("%s_party2", walletAddress)

	party1Bytes, err := json.Marshal(party1Data)
	if err != nil {
		return fmt.Errorf("failed to marshal party1 data: %w", err)
	}

	party2Bytes, err := json.Marshal(party2Data)
	if err != nil {
		return fmt.Errorf("failed to marshal party2 data: %w", err)
	}

	store.Set([]byte(party1Key), party1Bytes)
	store.Set([]byte(party2Key), party2Bytes)
	return nil
}

// GetPartyData retrieves TSS party data for a wallet
func (s *Storage) GetPartyData(ctx context.Context, walletAddress string) (*keygen.LocalPartySaveData, *keygen.LocalPartySaveData, error) {
	store := prefix.NewStore(s.store, []byte(types.PartyDataKey))
	party1Key := fmt.Sprintf("%s_party1", walletAddress)
	party2Key := fmt.Sprintf("%s_party2", walletAddress)

	party1Bytes := store.Get([]byte(party1Key))
	party2Bytes := store.Get([]byte(party2Key))

	if party1Bytes == nil || party2Bytes == nil {
		return nil, nil, fmt.Errorf("party data not found for wallet: %s", walletAddress)
	}

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
func (s *Storage) DeletePartyData(ctx context.Context, walletAddress string) error {
	store := prefix.NewStore(s.store, []byte(types.PartyDataKey))
	party1Key := fmt.Sprintf("%s_party1", walletAddress)
	party2Key := fmt.Sprintf("%s_party2", walletAddress)
	
	store.Delete([]byte(party1Key))
	store.Delete([]byte(party2Key))
	return nil
}
