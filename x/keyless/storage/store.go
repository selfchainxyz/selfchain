package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

var (
	// Key prefixes for the store
	SharePrefix     = []byte("share/")
	VersionPrefix   = []byte("version/")
	MetadataPrefix  = []byte("metadata/")
)

// Store implements secure storage for key shares
type Store struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
}

// NewStore creates a new secure storage instance
func NewStore(storeKey storetypes.StoreKey, cdc codec.BinaryCodec) *Store {
	return &Store{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// StoreShare stores an encrypted share
func (s *Store) StoreShare(ctx sdk.Context, walletID string, share []byte, key []byte) error {
	if len(share) == 0 {
		return errors.New("empty share")
	}
	if len(key) == 0 {
		return errors.New("empty encryption key")
	}

	// Encrypt the share
	encryptedShare, err := encrypt(share, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt share: %w", err)
	}

	// Get current version
	version := s.getNextVersion(ctx, walletID)

	// Store encrypted share
	store := ctx.KVStore(s.storeKey)
	shareKey := append(SharePrefix, []byte(walletID)...)
	store.Set(shareKey, encryptedShare)

	// Update metadata
	now := time.Now()
	metadata := types.ShareMetadata{
		Version:     version,
		CreatedAt:   &now,
		UpdatedAt:   &now,
		BackupState: types.BackupStatus_BACKUP_STATUS_NONE,
	}
	if err := s.storeMetadata(ctx, walletID, &metadata); err != nil {
		return fmt.Errorf("failed to store metadata: %w", err)
	}

	return nil
}

// RetrieveShare retrieves and decrypts a share
func (s *Store) RetrieveShare(ctx sdk.Context, walletID string, key []byte) ([]byte, error) {
	store := ctx.KVStore(s.storeKey)
	shareKey := append(SharePrefix, []byte(walletID)...)
	
	encryptedShare := store.Get(shareKey)
	if encryptedShare == nil {
		return nil, fmt.Errorf("share not found for wallet: %s", walletID)
	}

	// Decrypt the share
	share, err := decrypt(encryptedShare, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt share: %w", err)
	}

	return share, nil
}

// UpdateShare updates an existing share
func (s *Store) UpdateShare(ctx sdk.Context, walletID string, newShare []byte, key []byte) error {
	// Check if share exists
	store := ctx.KVStore(s.storeKey)
	shareKey := append(SharePrefix, []byte(walletID)...)
	if !store.Has(shareKey) {
		return fmt.Errorf("share not found for wallet: %s", walletID)
	}

	// Encrypt and store new share
	encryptedShare, err := encrypt(newShare, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt share: %w", err)
	}

	// Update version
	version := s.getNextVersion(ctx, walletID)

	// Store new share
	store.Set(shareKey, encryptedShare)

	// Update metadata
	metadata, err := s.getMetadata(ctx, walletID)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}
	now := time.Now()
	metadata.Version = version
	metadata.UpdatedAt = &now

	if err := s.storeMetadata(ctx, walletID, metadata); err != nil {
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	return nil
}

// DeleteShare deletes a share and its metadata
func (s *Store) DeleteShare(ctx sdk.Context, walletID string) error {
	store := ctx.KVStore(s.storeKey)
	
	// Delete share
	shareKey := append(SharePrefix, []byte(walletID)...)
	store.Delete(shareKey)

	// Delete metadata
	metadataKey := append(MetadataPrefix, []byte(walletID)...)
	store.Delete(metadataKey)

	// Delete version
	versionKey := append(VersionPrefix, []byte(walletID)...)
	store.Delete(versionKey)

	return nil
}

// Helper functions

func (s *Store) getNextVersion(ctx sdk.Context, walletID string) uint64 {
	store := ctx.KVStore(s.storeKey)
	versionKey := append(VersionPrefix, []byte(walletID)...)
	
	var version uint64 = 1
	bz := store.Get(versionKey)
	if bz != nil {
		version = sdk.BigEndianToUint64(bz) + 1
	}
	
	store.Set(versionKey, sdk.Uint64ToBigEndian(version))
	return version
}

func (s *Store) storeMetadata(ctx sdk.Context, walletID string, metadata *types.ShareMetadata) error {
	store := ctx.KVStore(s.storeKey)
	metadataKey := append(MetadataPrefix, []byte(walletID)...)
	
	bz, err := s.cdc.Marshal(metadata)
	if err != nil {
		return err
	}
	
	store.Set(metadataKey, bz)
	return nil
}

func (s *Store) getMetadata(ctx sdk.Context, walletID string) (*types.ShareMetadata, error) {
	store := ctx.KVStore(s.storeKey)
	metadataKey := append(MetadataPrefix, []byte(walletID)...)
	
	bz := store.Get(metadataKey)
	if bz == nil {
		return nil, fmt.Errorf("no metadata found for wallet %s", walletID)
	}
	
	var metadata types.ShareMetadata
	if err := s.cdc.Unmarshal(bz, &metadata); err != nil {
		return nil, err
	}
	return &metadata, nil
}

// Encryption helpers

func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
