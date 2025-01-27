package types

// Event types and attribute keys for the keyless module
const (
	EventTypeGrantPermission  = "grant_permission"
	EventTypeRevokePermission = "revoke_permission"
	EventTypeCreateWallet     = "create_wallet"
	EventTypeRecoverWallet    = "recover_wallet"
	EventTypeRotateKey        = "rotate_key"

	EventTypeWalletCreated         = "wallet_created"
	EventTypeWalletRecovered       = "wallet_recovered"
	EventTypeKeyRotationInitiated  = "key_rotation_initiated"
	EventTypeKeyRotationCompleted  = "key_rotation_completed"
	EventTypeTransactionSigned     = "transaction_signed"
	EventTypeTransactionBatchSigned = "transaction_batch_signed"
	EventTypeBatchSignRequested    = "batch_sign_requested"
	EventTypeRecoveryStarted       = "recovery_started"
	EventTypeRecoveryCompleted     = "recovery_completed"
	EventTypePermissionGranted     = "permission_granted"
	EventTypePermissionRevoked     = "permission_revoked"

	AttributeKeyWalletID      = "wallet_id"
	AttributeKeyGrantee       = "grantee"
	AttributeKeyPermissions   = "permissions"
	AttributeKeyCreator       = "creator"
	AttributeKeyWalletAddress = "wallet_address"
	AttributeKeyChainID       = "chain_id"
	AttributeKeyStatus        = "status"
	AttributeKeyNewPubKey     = "new_pub_key"
	AttributeKeyVersion       = "version"
	AttributeKeyTxHash       = "tx_hash"
	AttributeKeyBatchSize    = "batch_size"
	AttributeKeyBatchStatus  = "batch_status"
	AttributeKeyRecoveryAddress = "recovery_address"
	AttributeKeyNewOwner        = "new_owner"
	AttributeKeyTimestamp       = "timestamp"
)
