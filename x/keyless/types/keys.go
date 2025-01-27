package types

import (
	"encoding/binary"
)

const (
	// ModuleName defines the module name
	ModuleName = "keyless"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_keyless"

	// Version defines the current version the keyless module
	Version = 1
)

const (
	// KeyPrefix is the prefix for all keys in the store
	KeyPrefix = "keyless/"

	// WalletKey is the prefix for wallet keys
	WalletKey = KeyPrefix + "wallet/"

	// KeyRotationKey is the prefix for key rotation keys
	KeyRotationKey = KeyPrefix + "key_rotation/"

	// AuthorizationKeyPrefix is the prefix for authorization keys
	AuthorizationKeyPrefix = KeyPrefix + "authorization/"

	// KeyRotationStatusPrefix is the prefix for key rotation status
	KeyRotationStatusPrefix = KeyPrefix + "key_rotation_status/"

	// BatchSignStatusKey is the prefix for batch sign status
	BatchSignStatusKey = KeyPrefix + "batch_sign_status/"

	// ParamsKey is the prefix for params
	ParamsKey = KeyPrefix + "params/"

	// KeyShareKey is the prefix for key shares
	KeyShareKey = KeyPrefix + "key_share/"

	// PartyDataKey is the prefix for party data
	PartyDataKey = KeyPrefix + "party_data/"

	// PermissionKey is the prefix for permissions
	PermissionKey = KeyPrefix + "permission/"

	// SigningSessionKey is the prefix for signing sessions
	SigningSessionKey = KeyPrefix + "signing_session/"

	// RecoveryKeyPrefix is the prefix for recovery keys
	RecoveryKeyPrefix = KeyPrefix + "recovery/"

	// AuditEventKey is the prefix for audit events
	AuditEventKey = KeyPrefix + "audit_event/"
)

var (
	// KeyPrefixWallet is the prefix for wallet storage
	KeyPrefixWallet = []byte{0x01}

	// KeyPrefixPartyData is the prefix for party data storage
	KeyPrefixPartyData = []byte{0x02}

	// KeyPrefixSigningSession is the prefix for signing session storage
	KeyPrefixSigningSession = []byte{0x03}
)

// GetPrefixedKey returns a key with the module prefix
func GetPrefixedKey(key []byte) []byte {
	return append([]byte(KeyPrefix), key...)
}

// WalletStoreKey returns the store key to retrieve a Wallet from the index fields
func WalletStoreKey(walletId string) []byte {
	key := []byte(WalletKey)
	return append(key, []byte(walletId)...)
}

// AuthorizationKey returns the store key to retrieve an authorization from the index fields
func AuthorizationKey(creator, walletAddress string) []byte {
	key := []byte(AuthorizationKeyPrefix)
	key = append(key, []byte(creator)...)
	return append(key, []byte(walletAddress)...)
}

// KeyRotationStatusKey returns the store key to retrieve a key rotation status from the index fields
func KeyRotationStatusKey(walletId string) []byte {
	key := []byte(KeyRotationStatusPrefix)
	return append(key, []byte(walletId)...)
}

// GetKeyRotationKey returns the store key to retrieve a KeyRotation from the index fields
func GetKeyRotationKey(walletId string, version uint64) []byte {
	key := []byte(KeyRotationKey)
	key = append(key, []byte(walletId)...)
	versionBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(versionBytes, version)
	return append(key, versionBytes...)
}

// BatchSignStatusStoreKey returns the store key for batch sign status
func BatchSignStatusStoreKey(walletId string) []byte {
	key := []byte(BatchSignStatusKey)
	return append(key, []byte(walletId)...)
}

// ParamsStoreKey returns the store key for module parameters
func ParamsStoreKey() []byte {
	return []byte(ParamsKey)
}

// GetPermissionKey returns the store key for a permission
func GetPermissionKey(walletID, grantee string) []byte {
	return append(GetPermissionPrefix(walletID), []byte(grantee)...)
}

// GetPermissionPrefix returns the prefix for all permissions of a wallet
func GetPermissionPrefix(walletID string) []byte {
	return append([]byte(PermissionKey), []byte(walletID)...)
}

// GetRecoveryKey returns the key for storing recovery data
func GetRecoveryKey(walletID string) []byte {
	return append([]byte(RecoveryKeyPrefix), []byte(walletID)...)
}

// GetAuditEventKey returns the store key to retrieve an AuditEvent from the index fields
func GetAuditEventKey(walletId string, timestamp int64) []byte {
	key := []byte(AuditEventKey)
	key = append(key, []byte(walletId)...)
	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(timestamp))
	return append(key, timeBytes...)
}
