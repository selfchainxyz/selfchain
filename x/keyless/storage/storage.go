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

// Storage interface defines storage operations
type Storage interface {
	SavePartyShare(ctx context.Context, key string, share *types.EncryptedShare) error
	GetPartyShare(ctx context.Context, key string) (*types.EncryptedShare, error)
	DeletePartyShare(ctx context.Context, key string)
	SavePartyData(ctx context.Context, walletAddress string, party1Data, party2Data *keygen.LocalPartySaveData) error
	GetPartyData(ctx context.Context, walletAddress string) (*keygen.LocalPartySaveData, *keygen.LocalPartySaveData, error)
	DeletePartyData(ctx context.Context, walletAddress string) error
	GetWalletMetadata(ctx context.Context, walletAddress string) (*types.WalletMetadata, error)
	SaveWalletMetadata(ctx context.Context, walletAddress string, metadata *types.WalletMetadata) error
}

// StorageImpl manages data persistence
type StorageImpl struct {
	store sdk.KVStore
}

var _ Storage = (*StorageImpl)(nil)

// NewStorage creates a new storage instance
func NewStorage(store sdk.KVStore) Storage {
	return &StorageImpl{
		store: store,
	}
}

// SavePartyShare stores an encrypted key share
func (s *StorageImpl) SavePartyShare(ctx context.Context, key string, share *types.EncryptedShare) error {
	store := prefix.NewStore(s.store, []byte(types.KeyShareKey))
	bz, err := share.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal share: %w", err)
	}
	store.Set([]byte(key), bz)
	return nil
}

// GetPartyShare retrieves an encrypted key share
func (s *StorageImpl) GetPartyShare(ctx context.Context, key string) (*types.EncryptedShare, error) {
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
func (s *StorageImpl) DeletePartyShare(ctx context.Context, key string) {
	store := prefix.NewStore(s.store, []byte(types.KeyShareKey))
	store.Delete([]byte(key))
}

// SavePartyData stores TSS party data for a wallet
func (s *StorageImpl) SavePartyData(ctx context.Context, walletAddress string, party1Data, party2Data *keygen.LocalPartySaveData) error {
	store := prefix.NewStore(s.store, []byte(types.PartyDataKey))
	
	// Marshal party1 data
	party1Bytes, err := json.Marshal(party1Data)
	if err != nil {
		return fmt.Errorf("failed to marshal party1 data: %w", err)
	}
	
	// Marshal party2 data
	party2Bytes, err := json.Marshal(party2Data)
	if err != nil {
		return fmt.Errorf("failed to marshal party2 data: %w", err)
	}
	
	// Store both party data
	store.Set([]byte(walletAddress+"_party1"), party1Bytes)
	store.Set([]byte(walletAddress+"_party2"), party2Bytes)
	
	return nil
}

// GetPartyData retrieves TSS party data for a wallet
func (s *StorageImpl) GetPartyData(ctx context.Context, walletAddress string) (*keygen.LocalPartySaveData, *keygen.LocalPartySaveData, error) {
	store := prefix.NewStore(s.store, []byte(types.PartyDataKey))
	
	// Get party1 data
	party1Bytes := store.Get([]byte(walletAddress + "_party1"))
	if party1Bytes == nil {
		return nil, nil, fmt.Errorf("party1 data not found for wallet: %s", walletAddress)
	}
	
	// Get party2 data
	party2Bytes := store.Get([]byte(walletAddress + "_party2"))
	if party2Bytes == nil {
		return nil, nil, fmt.Errorf("party2 data not found for wallet: %s", walletAddress)
	}
	
	// Unmarshal party1 data
	var party1Data keygen.LocalPartySaveData
	if err := json.Unmarshal(party1Bytes, &party1Data); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal party1 data: %w", err)
	}
	
	// Unmarshal party2 data
	var party2Data keygen.LocalPartySaveData
	if err := json.Unmarshal(party2Bytes, &party2Data); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal party2 data: %w", err)
	}
	
	return &party1Data, &party2Data, nil
}

// DeletePartyData removes TSS party data for a wallet
func (s *StorageImpl) DeletePartyData(ctx context.Context, walletAddress string) error {
	store := prefix.NewStore(s.store, []byte(types.PartyDataKey))
	store.Delete([]byte(walletAddress + "_party1"))
	store.Delete([]byte(walletAddress + "_party2"))
	return nil
}

// GetWalletMetadata retrieves metadata for a wallet
func (s *StorageImpl) GetWalletMetadata(ctx context.Context, walletAddress string) (*types.WalletMetadata, error) {
	store := prefix.NewStore(s.store, []byte(types.WalletMetadataKey))
	bz := store.Get([]byte(walletAddress))
	if bz == nil {
		return nil, fmt.Errorf("metadata not found for wallet: %s", walletAddress)
	}

	var metadata types.WalletMetadata
	if err := metadata.Unmarshal(bz); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}
	return &metadata, nil
}

// SaveWalletMetadata stores metadata for a wallet
func (s *StorageImpl) SaveWalletMetadata(ctx context.Context, walletAddress string, metadata *types.WalletMetadata) error {
	store := prefix.NewStore(s.store, []byte(types.WalletMetadataKey))
	bz, err := metadata.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	store.Set([]byte(walletAddress), bz)
	return nil
}
